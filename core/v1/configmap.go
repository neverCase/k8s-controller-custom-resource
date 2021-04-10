package v1

import (
	"context"
	"fmt"
	"github.com/nevercase/k8s-controller-custom-resource/pkg/env"
	"time"

	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	kubeInformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	coreListersV1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/klog/v2"
)

type KubernetesConfigMap interface {
	Get(nameSpace, specDeploymentName string) (d *coreV1.ConfigMap, err error)
	Create(nameSpace, specDeploymentName string, d *coreV1.ConfigMap) (*coreV1.ConfigMap, error)
	Update(nameSpace string, d *coreV1.ConfigMap) (*coreV1.ConfigMap, error)
	Delete(nameSpace, specDeploymentName string) error
	List(nameSpace, filterName string) (sl *coreV1.ConfigMapList, err error)
}

func NewKubernetesConfigMap(kubeClientSet kubernetes.Interface, kubeInformerFactory kubeInformers.SharedInformerFactory) KubernetesConfigMap {
	timeout, _ := env.GetExecutionTimeoutDuration()
	return &kubernetesConfigMap{
		kubeClientSet:         kubeClientSet,
		configMapLister:       kubeInformerFactory.Core().V1().ConfigMaps().Lister(),
		executionTimeoutInSec: timeout,
	}
}

type kubernetesConfigMap struct {
	kubeClientSet         kubernetes.Interface
	configMapLister       coreListersV1.ConfigMapLister
	executionTimeoutInSec int64
}

func (kcm *kubernetesConfigMap) Get(nameSpace, specDeploymentName string) (d *coreV1.ConfigMap, err error) {
	var configMapName string
	if specDeploymentName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		utilruntime.HandleError(fmt.Errorf("%s: DeploymentName must be specified", specDeploymentName))
		return d, fmt.Errorf("%s: DeploymentName must be specified", specDeploymentName)
	}
	configMapName = fmt.Sprintf(ConfigMapTemplate, specDeploymentName)
	// Get the configMap with the name specified in spec
	configMap, err := kcm.configMapLister.ConfigMaps(nameSpace).Get(configMapName)
	return configMap, err
}

func (kcm *kubernetesConfigMap) Create(nameSpace, specDeploymentName string, d *coreV1.ConfigMap) (*coreV1.ConfigMap, error) {
	createOpt := metav1.CreateOptions{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(kcm.executionTimeoutInSec))
	configMap, err := kcm.kubeClientSet.CoreV1().ConfigMaps(nameSpace).Create(ctx, d, createOpt)
	cancel()
	if err != nil {
		klog.V(2).Info(err)
	}
	return configMap, err
}

func (kcm *kubernetesConfigMap) Update(nameSpace string, d *coreV1.ConfigMap) (*coreV1.ConfigMap, error) {
	updateOpt := metav1.UpdateOptions{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(kcm.executionTimeoutInSec))
	configMap, err := kcm.kubeClientSet.CoreV1().ConfigMaps(nameSpace).Update(ctx, d, updateOpt)
	cancel()
	if err != nil {
		klog.V(2).Info(err)
	}
	return configMap, err
}

func (kcm *kubernetesConfigMap) Delete(nameSpace, specDeploymentName string) error {
	// Get the configMap with the name specified in spec
	_, err := kcm.Get(nameSpace, specDeploymentName)
	// If the resource doesn't exist, we'll return nil
	if errors.IsNotFound(err) {
		return nil
	}
	opts := metav1.DeleteOptions{
		//GracePeriodSeconds: int64ToPointer(30),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(kcm.executionTimeoutInSec))
	configMapName := fmt.Sprintf(ConfigMapTemplate, specDeploymentName)
	err = kcm.kubeClientSet.CoreV1().ConfigMaps(nameSpace).Delete(ctx, configMapName, opts)
	cancel()
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	return nil
}

func (kcm *kubernetesConfigMap) List(nameSpace, filterName string) (sl *coreV1.ConfigMapList, err error) {
	opts := metav1.ListOptions{
		LabelSelector: filterName,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(kcm.executionTimeoutInSec))
	sl, err = kcm.kubeClientSet.CoreV1().ConfigMaps(nameSpace).List(ctx, opts)
	cancel()
	if err != nil {
		klog.V(2).Info(err)
	}
	return sl, err
}
