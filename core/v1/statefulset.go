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

type KubernetesStatefulSet interface {
	Get(nameSpace, specDeploymentName string) (d *appsv1.StatefulSet, err error)
	Create(nameSpace, specDeploymentName string, d *appsv1.StatefulSet) (*appsv1.StatefulSet, error)
	Update(nameSpace string, d *appsv1.StatefulSet) (*appsv1.StatefulSet, error)
	Delete(nameSpace, specDeploymentName string) error
	List(nameSpace, filterName string) (dl *appsv1.StatefulSetList, err error)
}

func NewKubernetesStatefulSet(kubeClientSet kubernetes.Interface, kubeInformerFactory kubeinformers.SharedInformerFactory, recorder record.EventRecorder) KubernetesStatefulSet {
	var kd KubernetesStatefulSet = &kubernetesStatefulSet{
		kubeClientSet:     kubeClientSet,
		statefulSetLister: kubeInformerFactory.Apps().V1().StatefulSets().Lister(),
		recorder:          recorder,
	}
	return kd
}

type kubernetesStatefulSet struct {
	kubeClientSet     kubernetes.Interface
	statefulSetLister appslistersv1.StatefulSetLister
	recorder          record.EventRecorder
}

func (kd *kubernetesStatefulSet) Get(nameSpace, specDeploymentName string) (d *appsv1.StatefulSet, err error) {
	var deploymentName string
	if specDeploymentName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		utilruntime.HandleError(fmt.Errorf("%s: DeploymentName must be specified", specDeploymentName))
		return d, fmt.Errorf("%s: DeploymentName must be specified", specDeploymentName)
	}
	deploymentName = fmt.Sprintf(StatefulSetNameTemplate, specDeploymentName)
	// Get the statefulSet with the name specified in RedisOperator.spec
	statefulSet, err := kd.statefulSetLister.StatefulSets(nameSpace).Get(deploymentName)
	return statefulSet, err
}

func (kd *kubernetesStatefulSet) Create(nameSpace, specDeploymentName string, d *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	statefulSet, err := kd.kubeClientSet.AppsV1().StatefulSets(nameSpace).Create(d)
	if err != nil {
		klog.V(2).Info(err)
	}
	return statefulSet, err
}

func (kd *kubernetesStatefulSet) Update(nameSpace string, d *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	statefulSet, err := kd.kubeClientSet.AppsV1().StatefulSets(nameSpace).Update(d)
	if err != nil {
		klog.V(2).Info(err)
	}
	return statefulSet, err
}

func (kd *kubernetesStatefulSet) Delete(nameSpace, specDeploymentName string) error {
	// Get the statefulSet with the name specified in RedisOperator.spec
	_, err := kd.Get(nameSpace, specDeploymentName)
	// If the resource doesn't exist, we'll return nil
	if errors.IsNotFound(err) {
		return nil
	}
	opts := &metav1.DeleteOptions{
		//GracePeriodSeconds: int64ToPointer(30),
	}
	deploymentName := fmt.Sprintf(StatefulSetNameTemplate, specDeploymentName)
	err = kd.kubeClientSet.AppsV1().StatefulSets(nameSpace).Delete(deploymentName, opts)
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	return nil
}

func (kd *kubernetesStatefulSet) List(nameSpace, filterName string) (dl *appsv1.StatefulSetList, err error) {
	opts := metav1.ListOptions{
		LabelSelector: filterName,
	}
	dl, err = kd.kubeClientSet.AppsV1().StatefulSets(nameSpace).List(opts)
	if err != nil {
		klog.V(2).Info(err)
	}
	return dl, err
}
