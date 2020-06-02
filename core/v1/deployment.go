package v1

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	appslistersv1 "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
)

type KubernetesDeployment interface {
	Get(nameSpace, specName string) (d *appsv1.Deployment, err error)
	Create(nameSpace string, d *appsv1.Deployment) (*appsv1.Deployment, error)
	Update(nameSpace string, d *appsv1.Deployment) (*appsv1.Deployment, error)
	Delete(nameSpace, specName string) error
	List(nameSpace, filterName string) (dl *appsv1.DeploymentList, err error)
}

func NewKubernetesDeployment(kubeClientSet kubernetes.Interface, kubeInformerFactory kubeinformers.SharedInformerFactory, recorder record.EventRecorder) KubernetesDeployment {
	var kd KubernetesDeployment = &kubernetesDeployment{
		kubeClientSet:     kubeClientSet,
		deploymentsLister: kubeInformerFactory.Apps().V1().Deployments().Lister(),
		recorder:          recorder,
	}
	return kd
}

type kubernetesDeployment struct {
	kubeClientSet     kubernetes.Interface
	deploymentsLister appslistersv1.DeploymentLister
	recorder          record.EventRecorder
}

func (kd *kubernetesDeployment) Get(nameSpace, specName string) (d *appsv1.Deployment, err error) {
	var name string
	if specName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		utilruntime.HandleError(fmt.Errorf("%s: name must be specified", specName))
		return d, fmt.Errorf("%s: name must be specified", specName)
	}
	name = fmt.Sprintf(DeploymentNameTemplate, specName)
	// Get the deployment with the name specified in spec
	deployment, err := kd.deploymentsLister.Deployments(nameSpace).Get(name)
	return deployment, err
}

func (kd *kubernetesDeployment) Create(nameSpace string, d *appsv1.Deployment) (*appsv1.Deployment, error) {
	deployment, err := kd.kubeClientSet.AppsV1().Deployments(nameSpace).Create(d)
	if err != nil {
		klog.V(2).Info(err)
	}
	return deployment, err
}

func (kd *kubernetesDeployment) Update(nameSpace string, d *appsv1.Deployment) (*appsv1.Deployment, error) {
	deployment, err := kd.kubeClientSet.AppsV1().Deployments(nameSpace).Update(d)
	if err != nil {
		klog.V(2).Info(err)
	}
	return deployment, err
}

func (kd *kubernetesDeployment) Delete(nameSpace, specName string) error {
	// Get the deployment with the name specified in spec
	_, err := kd.Get(nameSpace, specName)
	// If the resource doesn't exist, we'll return nil
	if errors.IsNotFound(err) {
		return nil
	}
	opts := &metav1.DeleteOptions{
		//GracePeriodSeconds: int64ToPointer(30),
	}
	name := fmt.Sprintf(DeploymentNameTemplate, specName)
	err = kd.kubeClientSet.AppsV1().Deployments(nameSpace).Delete(name, opts)
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	return nil
}

func (kd *kubernetesDeployment) List(nameSpace, filterName string) (dl *appsv1.DeploymentList, err error) {
	opts := metav1.ListOptions{
		LabelSelector: filterName,
	}
	dl, err = kd.kubeClientSet.AppsV1().Deployments(nameSpace).List(opts)
	if err != nil {
		klog.V(2).Info(err)
	}
	return dl, err
}
