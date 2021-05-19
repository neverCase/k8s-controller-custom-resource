package v1

import (
	"context"
	"fmt"
	"github.com/nevercase/k8s-controller-custom-resource/pkg/env"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	corelistersv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/klog/v2"
)

type KubernetesService interface {
	Get(nameSpace, specName string) (d *corev1.Service, err error)
	Create(nameSpace string, d *corev1.Service) (*corev1.Service, error)
	Update(nameSpace string, d *corev1.Service) (*corev1.Service, error)
	Delete(nameSpace, specName string) error
	List(nameSpace, filterName string) (sl *corev1.ServiceList, err error)
}

func NewKubernetesService(kubeClientSet kubernetes.Interface, kubeInformerFactory kubeinformers.SharedInformerFactory) KubernetesService {
	timeout, _ := env.GetExecutionTimeoutDuration()
	return &kubernetesService{
		kubeClientSet:         kubeClientSet,
		servicesLister:        kubeInformerFactory.Core().V1().Services().Lister(),
		executionTimeoutInSec: timeout,
	}
}

type kubernetesService struct {
	kubeClientSet         kubernetes.Interface
	servicesLister        corelistersv1.ServiceLister
	executionTimeoutInSec int64
}

func (ks *kubernetesService) Get(nameSpace, specName string) (d *corev1.Service, err error) {
	var name string
	if specName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		utilruntime.HandleError(fmt.Errorf("%s: DeploymentName must be specified", specName))
		return d, fmt.Errorf("%s: DeploymentName must be specified", specName)
	}
	name = fmt.Sprintf(ServiceNameTemplate, specName)
	// Get the service with the name specified in spec
	service, err := ks.servicesLister.Services(nameSpace).Get(name)
	return service, err
}

func (ks *kubernetesService) Create(nameSpace string, d *corev1.Service) (*corev1.Service, error) {
	createOpt := metav1.CreateOptions{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(ks.executionTimeoutInSec))
	service, err := ks.kubeClientSet.CoreV1().Services(nameSpace).Create(ctx, d, createOpt)
	cancel()
	if err != nil {
		klog.V(2).Info(err)
	}
	return service, err
}

func (ks *kubernetesService) Update(nameSpace string, d *corev1.Service) (*corev1.Service, error) {
	updateOpt := metav1.UpdateOptions{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(ks.executionTimeoutInSec))
	service, err := ks.kubeClientSet.CoreV1().Services(nameSpace).Update(ctx, d, updateOpt)
	cancel()
	if err != nil {
		klog.V(2).Info(err)
	}
	return service, err
}

func (ks *kubernetesService) Delete(nameSpace, specName string) error {
	// Get the service with the name specified in spec
	_, err := ks.Get(nameSpace, specName)
	// If the resource doesn't exist, we'll return nil
	if errors.IsNotFound(err) {
		return nil
	}
	opts := metav1.DeleteOptions{
		//GracePeriodSeconds: int64ToPointer(30),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(ks.executionTimeoutInSec))
	name := fmt.Sprintf(ServiceNameTemplate, specName)
	err = ks.kubeClientSet.CoreV1().Services(nameSpace).Delete(ctx, name, opts)
	cancel()
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	return nil
}

func (ks *kubernetesService) List(nameSpace, filterName string) (sl *corev1.ServiceList, err error) {
	opts := metav1.ListOptions{
		LabelSelector: filterName,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(ks.executionTimeoutInSec))
	sl, err = ks.kubeClientSet.CoreV1().Services(nameSpace).List(ctx, opts)
	cancel()
	if err != nil {
		klog.V(2).Info(err)
	}
	return sl, err
}

func GetServiceType(st corev1.ServiceType) corev1.ServiceType {
	if st == "" {
		return corev1.ServiceTypeClusterIP
	}
	return st
}
