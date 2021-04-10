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
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	appslistersv1 "k8s.io/client-go/listers/apps/v1"
	"k8s.io/klog/v2"
)

type KubernetesDeployment interface {
	Get(nameSpace, specName string) (d *appsv1.Deployment, err error)
	Create(nameSpace string, d *appsv1.Deployment) (*appsv1.Deployment, error)
	Update(nameSpace string, d *appsv1.Deployment) (*appsv1.Deployment, error)
	Delete(nameSpace, specName string) error
	List(nameSpace, filterName string) (dl *appsv1.DeploymentList, err error)
	Patch(nameSpace string, name string, pt types.PatchType, data []byte, subResources ...string) (*appsv1.Deployment, error)
}

func NewKubernetesDeployment(kubeClientSet kubernetes.Interface, kubeInformerFactory kubeinformers.SharedInformerFactory) KubernetesDeployment {
	timeout, _ := env.GetExecutionTimeoutDuration()
	return &kubernetesDeployment{
		kubeClientSet:         kubeClientSet,
		deploymentsLister:     kubeInformerFactory.Apps().V1().Deployments().Lister(),
		executionTimeoutInSec: timeout,
	}
}

type kubernetesDeployment struct {
	kubeClientSet         kubernetes.Interface
	deploymentsLister     appslistersv1.DeploymentLister
	executionTimeoutInSec int64
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
	createOpt := metav1.CreateOptions{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(kd.executionTimeoutInSec))
	deployment, err := kd.kubeClientSet.AppsV1().Deployments(nameSpace).Create(ctx, d, createOpt)
	cancel()
	if err != nil {
		klog.V(2).Info(err)
	}
	return deployment, err
}

func (kd *kubernetesDeployment) Update(nameSpace string, d *appsv1.Deployment) (*appsv1.Deployment, error) {
	updateOpt := metav1.UpdateOptions{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(kd.executionTimeoutInSec))
	deployment, err := kd.kubeClientSet.AppsV1().Deployments(nameSpace).Update(ctx, d, updateOpt)
	cancel()
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
	opts := metav1.DeleteOptions{
		//GracePeriodSeconds: int64ToPointer(30),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(kd.executionTimeoutInSec))
	name := fmt.Sprintf(DeploymentNameTemplate, specName)
	err = kd.kubeClientSet.AppsV1().Deployments(nameSpace).Delete(ctx, name, opts)
	cancel()
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(kd.executionTimeoutInSec))
	dl, err = kd.kubeClientSet.AppsV1().Deployments(nameSpace).List(ctx, opts)
	cancel()
	if err != nil {
		klog.V(2).Info(err)
	}
	return dl, err
}

func (kd *kubernetesDeployment) Patch(nameSpace string, name string, pt types.PatchType, data []byte, subResources ...string) (*appsv1.Deployment, error) {
	opts := metav1.PatchOptions{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(kd.executionTimeoutInSec))
	dl, err := kd.kubeClientSet.AppsV1().Deployments(nameSpace).Patch(ctx, name, pt, data, opts, subResources...)
	cancel()
	if err != nil {
		klog.V(2).Info(err)
	}
	return dl, err
}
