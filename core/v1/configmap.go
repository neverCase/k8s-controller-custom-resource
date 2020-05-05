package v1

import (
	"fmt"

	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	kubeInformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	coreListersV1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
)

type KubernetesConfigMap interface {
	Get(nameSpace, specDeploymentName string) (d *coreV1.ConfigMap, err error)
	Create(nameSpace, specDeploymentName string, d *coreV1.ConfigMap) (*coreV1.ConfigMap, error)
	Update(nameSpace string, d *coreV1.ConfigMap) (*coreV1.ConfigMap, error)
	Delete(nameSpace, specDeploymentName string) error
	List(nameSpace, filterName string) (sl *coreV1.ConfigMapList, err error)
}

func NewKubernetesConfigMap(kubeClientSet kubernetes.Interface, kubeInformerFactory kubeInformers.SharedInformerFactory, recorder record.EventRecorder) KubernetesConfigMap {
	var kcf KubernetesConfigMap = &kubernetesConfigMap{
		kubeClientSet:   kubeClientSet,
		configMapLister: kubeInformerFactory.Core().V1().ConfigMaps().Lister(),
		recorder:        recorder,
	}
	return kcf
}

type kubernetesConfigMap struct {
	kubeClientSet   kubernetes.Interface
	configMapLister coreListersV1.ConfigMapLister
	recorder        record.EventRecorder
}

func (kcf *kubernetesConfigMap) Get(nameSpace, specDeploymentName string) (d *coreV1.ConfigMap, err error) {
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
	configMap, err := kcf.configMapLister.ConfigMaps(nameSpace).Get(configMapName)
	return configMap, err
}

func (kcf *kubernetesConfigMap) Create(nameSpace, specDeploymentName string, d *coreV1.ConfigMap) (*coreV1.ConfigMap, error) {
	configMap, err := kcf.kubeClientSet.CoreV1().ConfigMaps(nameSpace).Create(d)
	if err != nil {
		klog.V(2).Info(err)
	}
	return configMap, err
}

func (kcf *kubernetesConfigMap) Update(nameSpace string, d *coreV1.ConfigMap) (*coreV1.ConfigMap, error) {
	configMap, err := kcf.kubeClientSet.CoreV1().ConfigMaps(nameSpace).Update(d)
	if err != nil {
		klog.V(2).Info(err)
	}
	return configMap, err
}

func (kcf *kubernetesConfigMap) Delete(nameSpace, specDeploymentName string) error {
	// Get the configMap with the name specified in spec
	_, err := kcf.Get(nameSpace, specDeploymentName)
	// If the resource doesn't exist, we'll return nil
	if errors.IsNotFound(err) {
		return nil
	}
	opts := &metaV1.DeleteOptions{
		//GracePeriodSeconds: int64ToPointer(30),
	}
	configMapName := fmt.Sprintf(ConfigMapTemplate, specDeploymentName)
	err = kcf.kubeClientSet.CoreV1().ConfigMaps(nameSpace).Delete(configMapName, opts)
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	return nil
}

func (kcf *kubernetesConfigMap) List(nameSpace, filterName string) (sl *coreV1.ConfigMapList, err error) {
	opts := metaV1.ListOptions{
		LabelSelector: filterName,
	}
	sl, err = kcf.kubeClientSet.CoreV1().ConfigMaps(nameSpace).List(opts)
	if err != nil {
		klog.V(2).Info(err)
	}
	return sl, err
}
