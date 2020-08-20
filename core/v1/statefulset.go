package v1

import (
	"fmt"

	appsV1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	appslistersv1 "k8s.io/client-go/listers/apps/v1"
	"k8s.io/klog"
)

type KubernetesStatefulSet interface {
	Get(nameSpace, specName string) (d *appsV1.StatefulSet, err error)
	Create(nameSpace string, d *appsV1.StatefulSet) (*appsV1.StatefulSet, error)
	Update(nameSpace string, d *appsV1.StatefulSet) (*appsV1.StatefulSet, error)
	Delete(nameSpace, specName string) error
	List(nameSpace, filterName string) (dl *appsV1.StatefulSetList, err error)
	Watch(nameSpace string, filter string) (w watch.Interface, err error)
}

func NewKubernetesStatefulSet(kubeClientSet kubernetes.Interface, kubeInformerFactory kubeinformers.SharedInformerFactory) KubernetesStatefulSet {
	return &kubernetesStatefulSet{
		kubeClientSet:     kubeClientSet,
		statefulSetLister: kubeInformerFactory.Apps().V1().StatefulSets().Lister(),
	}
}

type kubernetesStatefulSet struct {
	kubeClientSet     kubernetes.Interface
	statefulSetLister appslistersv1.StatefulSetLister
}

func (kss *kubernetesStatefulSet) Get(nameSpace, specName string) (ss *appsV1.StatefulSet, err error) {
	var name string
	if specName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		utilruntime.HandleError(fmt.Errorf("%s: name must be specified", specName))
		return ss, fmt.Errorf("%s: name must be specified", specName)
	}
	name = fmt.Sprintf(StatefulSetNameTemplate, specName)
	// Get the statefulSet with the name specified in spec
	statefulSet, err := kss.statefulSetLister.StatefulSets(nameSpace).Get(name)
	return statefulSet, err
}

func (kss *kubernetesStatefulSet) Create(nameSpace string, ss *appsV1.StatefulSet) (*appsV1.StatefulSet, error) {
	statefulSet, err := kss.kubeClientSet.AppsV1().StatefulSets(nameSpace).Create(ss)
	if err != nil {
		klog.V(2).Info(err)
	}
	return statefulSet, err
}

func (kss *kubernetesStatefulSet) Update(nameSpace string, ss *appsV1.StatefulSet) (*appsV1.StatefulSet, error) {
	statefulSet, err := kss.kubeClientSet.AppsV1().StatefulSets(nameSpace).Update(ss)
	if err != nil {
		klog.V(2).Info(err)
	}
	return statefulSet, err
}

func (kss *kubernetesStatefulSet) Delete(nameSpace, specName string) error {
	// Get the statefulSet with the name specified in spec
	_, err := kss.Get(nameSpace, specName)
	// If the resource doesn't exist, we'll return nil
	if errors.IsNotFound(err) {
		return nil
	}
	opts := &metaV1.DeleteOptions{
		//GracePeriodSeconds: int64ToPointer(30),
	}
	name := fmt.Sprintf(StatefulSetNameTemplate, specName)
	err = kss.kubeClientSet.AppsV1().StatefulSets(nameSpace).Delete(name, opts)
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	return nil
}

func (kss *kubernetesStatefulSet) List(nameSpace, filterName string) (ssl *appsV1.StatefulSetList, err error) {
	timeout := int64(300)
	opts := metaV1.ListOptions{
		LabelSelector:  filterName,
		TimeoutSeconds: &timeout,
	}
	ssl, err = kss.kubeClientSet.AppsV1().StatefulSets(nameSpace).List(opts)
	if err != nil {
		klog.V(2).Info(err)
	}
	return ssl, err
}

func (kss *kubernetesStatefulSet) Watch(nameSpace string, filterName string) (w watch.Interface, err error) {
	timeout := int64(300)
	opts := metaV1.ListOptions{
		LabelSelector:  filterName,
		TimeoutSeconds: &timeout,
	}
	w, err = kss.kubeClientSet.AppsV1().StatefulSets(nameSpace).Watch(opts)
	if err != nil {
		klog.V(2).Info(err)
	}
	return w, err
}
