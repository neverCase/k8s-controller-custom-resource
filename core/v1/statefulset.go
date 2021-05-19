package v1

import (
	"context"
	"fmt"
	"github.com/nevercase/k8s-controller-custom-resource/pkg/env"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	appslistersv1 "k8s.io/client-go/listers/apps/v1"
	"k8s.io/klog/v2"
)

type KubernetesStatefulSet interface {
	Get(nameSpace, specName string) (d *appsv1.StatefulSet, err error)
	Create(nameSpace string, d *appsv1.StatefulSet) (*appsv1.StatefulSet, error)
	Update(nameSpace string, d *appsv1.StatefulSet) (*appsv1.StatefulSet, error)
	Delete(nameSpace, specName string) error
	List(nameSpace, filterName string) (dl *appsv1.StatefulSetList, err error)
	Watch(nameSpace string, filter string) (w watch.Interface, err error)
	Patch(nameSpace string, name string, pt types.PatchType, data []byte, subResources ...string) (*appsv1.StatefulSet, error)
}

func NewKubernetesStatefulSet(kubeClientSet kubernetes.Interface, kubeInformerFactory kubeinformers.SharedInformerFactory) KubernetesStatefulSet {
	timeout, _ := env.GetExecutionTimeoutDuration()
	return &kubernetesStatefulSet{
		kubeClientSet:         kubeClientSet,
		statefulSetLister:     kubeInformerFactory.Apps().V1().StatefulSets().Lister(),
		executionTimeoutInSec: timeout,
	}
}

type kubernetesStatefulSet struct {
	kubeClientSet         kubernetes.Interface
	statefulSetLister     appslistersv1.StatefulSetLister
	executionTimeoutInSec int64
}

func (sts *kubernetesStatefulSet) Get(nameSpace, specName string) (ss *appsv1.StatefulSet, err error) {
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
	statefulSet, err := sts.statefulSetLister.StatefulSets(nameSpace).Get(name)
	return statefulSet, err
}

func (sts *kubernetesStatefulSet) Create(nameSpace string, ss *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	createOpt := metav1.CreateOptions{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(sts.executionTimeoutInSec))
	statefulSet, err := sts.kubeClientSet.AppsV1().StatefulSets(nameSpace).Create(ctx, ss, createOpt)
	cancel()
	if err != nil {
		klog.V(2).Info(err)
	}
	return statefulSet, err
}

func (sts *kubernetesStatefulSet) Update(nameSpace string, ss *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	opt := metav1.UpdateOptions{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(sts.executionTimeoutInSec))
	statefulSet, err := sts.kubeClientSet.AppsV1().StatefulSets(nameSpace).Update(ctx, ss, opt)
	cancel()
	if err != nil {
		klog.V(2).Info(err)
	}
	return statefulSet, err
}

func (sts *kubernetesStatefulSet) Delete(nameSpace, specName string) error {
	// Get the statefulSet with the name specified in spec
	_, err := sts.Get(nameSpace, specName)
	// If the resource doesn't exist, we'll return nil
	if errors.IsNotFound(err) {
		return nil
	}
	opts := metav1.DeleteOptions{
		//GracePeriodSeconds: int64ToPointer(30),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(sts.executionTimeoutInSec))
	name := fmt.Sprintf(StatefulSetNameTemplate, specName)
	err = sts.kubeClientSet.AppsV1().StatefulSets(nameSpace).Delete(ctx, name, opts)
	cancel()
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	return nil
}

func (sts *kubernetesStatefulSet) List(nameSpace, filterName string) (ssl *appsv1.StatefulSetList, err error) {
	timeout := int64(300)
	opts := metav1.ListOptions{
		LabelSelector:  filterName,
		TimeoutSeconds: &timeout,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(sts.executionTimeoutInSec))
	ssl, err = sts.kubeClientSet.AppsV1().StatefulSets(nameSpace).List(ctx, opts)
	cancel()
	if err != nil {
		klog.V(2).Info(err)
	}
	return ssl, err
}

func (sts *kubernetesStatefulSet) Watch(nameSpace string, filterName string) (w watch.Interface, err error) {
	timeout := int64(300)
	opts := metav1.ListOptions{
		LabelSelector:  filterName,
		TimeoutSeconds: &timeout,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(sts.executionTimeoutInSec))
	w, err = sts.kubeClientSet.AppsV1().StatefulSets(nameSpace).Watch(ctx, opts)
	cancel()
	if err != nil {
		klog.V(2).Info(err)
	}
	return w, err
}

func (sts *kubernetesStatefulSet) Patch(nameSpace string, name string, pt types.PatchType, data []byte, subResources ...string) (*appsv1.StatefulSet, error) {
	opt := metav1.PatchOptions{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(sts.executionTimeoutInSec))
	s, err := sts.kubeClientSet.AppsV1().StatefulSets(nameSpace).Patch(ctx, name, pt, data, opt, subResources...)
	cancel()
	if err != nil {
		klog.V(2).Info(err)
	}
	return s, err
}
