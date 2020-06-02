package v1

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	corelistersv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
)

type KubernetesService interface {
	Get(nameSpace, specName string) (d *corev1.Service, err error)
	Create(nameSpace string, d *corev1.Service) (*corev1.Service, error)
	Update(nameSpace string, d *corev1.Service) (*corev1.Service, error)
	Delete(nameSpace, specName string) error
	List(nameSpace, filterName string) (sl *corev1.ServiceList, err error)
}

func NewKubernetesService(kubeClientSet kubernetes.Interface, kubeInformerFactory kubeinformers.SharedInformerFactory, recorder record.EventRecorder) KubernetesService {
	var kd KubernetesService = &kubernetesService{
		kubeClientSet:  kubeClientSet,
		servicesLister: kubeInformerFactory.Core().V1().Services().Lister(),
		recorder:       recorder,
	}
	return kd
}

type kubernetesService struct {
	kubeClientSet  kubernetes.Interface
	servicesLister corelistersv1.ServiceLister
	recorder       record.EventRecorder
}

func (kd *kubernetesService) Get(nameSpace, specName string) (d *corev1.Service, err error) {
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
	service, err := kd.servicesLister.Services(nameSpace).Get(name)
	return service, err
}

func (kd *kubernetesService) Create(nameSpace string, d *corev1.Service) (*corev1.Service, error) {
	service, err := kd.kubeClientSet.CoreV1().Services(nameSpace).Create(d)
	if err != nil {
		klog.V(2).Info(err)
	}
	return service, err
}

func (kd *kubernetesService) Update(nameSpace string, d *corev1.Service) (*corev1.Service, error) {
	service, err := kd.kubeClientSet.CoreV1().Services(nameSpace).Update(d)
	if err != nil {
		klog.V(2).Info(err)
	}
	return service, err
}

func (kd *kubernetesService) Delete(nameSpace, specName string) error {
	// Get the service with the name specified in spec
	_, err := kd.Get(nameSpace, specName)
	// If the resource doesn't exist, we'll return nil
	if errors.IsNotFound(err) {
		return nil
	}
	opts := &metav1.DeleteOptions{
		//GracePeriodSeconds: int64ToPointer(30),
	}
	name := fmt.Sprintf(ServiceNameTemplate, specName)
	err = kd.kubeClientSet.CoreV1().Services(nameSpace).Delete(name, opts)
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	return nil
}

func (kd *kubernetesService) List(nameSpace, filterName string) (sl *corev1.ServiceList, err error) {
	opts := metav1.ListOptions{
		LabelSelector: filterName,
	}
	sl, err = kd.kubeClientSet.CoreV1().Services(nameSpace).List(opts)
	if err != nil {
		klog.V(2).Info(err)
	}
	return sl, err
}
