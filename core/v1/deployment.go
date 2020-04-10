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
	Get(nameSpace, specDeploymentName string) (d *appsv1.Deployment, err error)
	Create(nameSpace, specDeploymentName string, d *appsv1.Deployment) (*appsv1.Deployment, error)
	Update(nameSpace string, d *appsv1.Deployment) (*appsv1.Deployment, error)
	Delete(nameSpace, specDeploymentName string) error
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

func (kd *kubernetesDeployment) Get(nameSpace, specDeploymentName string) (d *appsv1.Deployment, err error) {
	var deploymentName string
	if specDeploymentName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		utilruntime.HandleError(fmt.Errorf("%s: DeploymentName must be specified", specDeploymentName))
		return d, fmt.Errorf("%s: DeploymentName must be specified", specDeploymentName)
	}
	deploymentName = fmt.Sprintf(DeploymentNameTemplate, specDeploymentName)
	// Get the deployment with the name specified in RedisOperator.spec
	deployment, err := kd.deploymentsLister.Deployments(nameSpace).Get(deploymentName)
	return deployment, err
}

func (kd *kubernetesDeployment) Create(nameSpace, specDeploymentName string, d *appsv1.Deployment) (*appsv1.Deployment, error) {
	_, err := kd.kubeClientSet.AppsV1().Deployments(nameSpace).Create(d)
	if err != nil {
		klog.V(2).Info(err)
	}
	return d, err
}

func (kd *kubernetesDeployment) Update(nameSpace string, d *appsv1.Deployment) (*appsv1.Deployment, error) {
	deployment, err := kd.kubeClientSet.AppsV1().Deployments(nameSpace).Update(d)
	if err != nil {
		klog.V(2).Info(err)
	}
	return deployment, err
}

func (kd *kubernetesDeployment) Delete(nameSpace, specDeploymentName string) error {
	// Get the deployment with the name specified in RedisOperator.spec
	_, err := kd.Get(nameSpace, specDeploymentName)
	// If the resource doesn't exist, we'll return nil
	if errors.IsNotFound(err) {
		return nil
	}
	opts := &metav1.DeleteOptions{
		//GracePeriodSeconds: int64ToPointer(30),
	}
	deploymentName := fmt.Sprintf(DeploymentNameTemplate, specDeploymentName)
	err = kd.kubeClientSet.AppsV1().Deployments(nameSpace).Delete(deploymentName, opts)
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	return nil
}
