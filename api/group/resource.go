package group

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	helixsagaoperatorv1 "github.com/Shanghai-Lunara/helixsaga-operator/pkg/apis/helixsaga/v1"
	helixsagaclientset "github.com/Shanghai-Lunara/helixsaga-operator/pkg/generated/helixsaga/clientset/versioned"
	mysqloperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/mysqloperator/v1"
	redisoperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/redisoperator/v1"
	mysqlclientset "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/mysqloperator/clientset/versioned"
	redisclientset "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/clientset/versioned"
)

type ResourceType string

const (
	// corev1
	ComponentStatus       ResourceType = "ComponentStatus"
	ConfigMap             ResourceType = "ConfigMap"
	Endpoint              ResourceType = "Endpoint"
	LimitRange            ResourceType = "LimitRange"
	Node                  ResourceType = "Node"
	NameSpace             ResourceType = "NameSpace"
	PersistentVolume      ResourceType = "PersistentVolume"
	PersistentVolumeClaim ResourceType = "PersistentVolumeClaim"
	Pod                   ResourceType = "Pod"
	PodTemplate           ResourceType = "PodTemplate"
	ReplicationController ResourceType = "ReplicationController"
	ResourceQuota         ResourceType = "ResourceQuota"
	Secret                ResourceType = "Secret"
	Service               ResourceType = "Service"
	ServiceAccount        ResourceType = "ServiceAccount"

	// appv1
	Deployment  ResourceType = "Deployment"
	StatefulSet ResourceType = "StatefulSet"

	// custom resource definition
	MysqlOperator     ResourceType = "MysqlOperator"
	RedisOperator     ResourceType = "RedisOperator"
	HelixSagaOperator ResourceType = "HelixSagaOperator"
)

// ResourceGetter has a method to return a ResourceInterface.
// A group's client should implement this interface.
type ResourceGetter interface {
	Resource() ResourceInterface
}

// ResourceInterface has methods to work with all Kubernetes resources include custom resource definitions.
type ResourceInterface interface {
	Create(rt ResourceType, nameSpace string, obj interface{}) (res interface{}, err error)
	Update(rt ResourceType, nameSpace string, obj interface{}) (res interface{}, err error)
	Delete(rt ResourceType, nameSpace, specName string) (err error)
	Get(rt ResourceType, nameSpace, specName string) (res interface{}, err error)
	List(rt ResourceType, nameSpace string, selector labels.Selector) (res interface{}, err error)
	Watch(rt ResourceType, nameSpace string, selector labels.Selector, eventsChan chan watch.Event) (err error)
	ResourceTypes() []ResourceType
}

// NewResource returns a ResourceInterface
func NewResource(ctx context.Context, masterUrl, kubeconfigPath string, eventsChan chan watch.Event) ResourceInterface {
	cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}
	mysql, err := mysqlclientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building mysqlclientset: %s", err.Error())
	}
	redis, err := redisclientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building redisclientset: %s", err.Error())
	}
	helixsaga, err := helixsagaclientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building redisclientset: %s", err.Error())
	}
	opts := NewOptions()
	var empty interface{}
	opts.Add(
		NewOption(NameSpace, empty),
		NewOption(ConfigMap, empty),
		NewOption(Pod, empty),
		NewOption(Secret, empty),
		NewOption(Service, empty),
		NewOption(ServiceAccount, empty),
		NewOption(MysqlOperator, mysql),
		NewOption(RedisOperator, redis),
		NewOption(HelixSagaOperator, helixsaga),
	)
	ctx2, cancel := context.WithCancel(ctx)
	var r = &resource{
		kubeClientSet: kubeClient,
		options:       opts,
		ctx:           ctx2,
		cancel:        cancel,
	}
	for _, v := range opts.GetOptionTypeList() {
		if v != Pod && v != ConfigMap && v != MysqlOperator && v != RedisOperator && v != HelixSagaOperator {
			continue
		}
		if err := r.Watch(v, "", labels.NewSelector(), eventsChan); err != nil {
			klog.V(2).Infof("Error watching ResourceType:%v err: %s", v, err)
		}
	}
	return r
}

// resource implements ResourceInterface
type resource struct {
	kubeClientSet kubernetes.Interface
	options       Options

	ctx    context.Context
	cancel context.CancelFunc
}

func (r *resource) Create(rt ResourceType, nameSpace string, obj interface{}) (res interface{}, err error) {
	var opt Option
	switch rt {
	case ConfigMap:
		res, err = r.kubeClientSet.CoreV1().ConfigMaps(nameSpace).Create(obj.(*corev1.ConfigMap))
	case NameSpace:
		res, err = r.kubeClientSet.CoreV1().Namespaces().Create(obj.(*corev1.Namespace))
	case Service:
		res, err = r.kubeClientSet.CoreV1().Services(nameSpace).Create(obj.(*corev1.Service))
	case MysqlOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*mysqlclientset.Clientset).NevercaseV1().MysqlOperators(nameSpace).Create(obj.(*mysqloperatorv1.MysqlOperator))
	case RedisOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*redisclientset.Clientset).NevercaseV1().RedisOperators(nameSpace).Create(obj.(*redisoperatorv1.RedisOperator))
	case HelixSagaOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*helixsagaclientset.Clientset).NevercaseV1().HelixSagas(nameSpace).Create(obj.(*helixsagaoperatorv1.HelixSaga))
	}
	if err != nil {
		klog.V(2).Info(err)
	}
	return res, err
}

func (r *resource) Update(rt ResourceType, nameSpace string, obj interface{}) (res interface{}, err error) {
	var opt Option
	switch rt {
	case ConfigMap:
		res, err = r.kubeClientSet.CoreV1().ConfigMaps(nameSpace).Update(obj.(*corev1.ConfigMap))
	case Service:
		res, err = r.kubeClientSet.CoreV1().Services(nameSpace).UpdateStatus(obj.(*corev1.Service))
	case MysqlOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*mysqlclientset.Clientset).NevercaseV1().MysqlOperators(nameSpace).Update(obj.(*mysqloperatorv1.MysqlOperator))
	case RedisOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*redisclientset.Clientset).NevercaseV1().RedisOperators(nameSpace).Update(obj.(*redisoperatorv1.RedisOperator))
	case HelixSagaOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*helixsagaclientset.Clientset).NevercaseV1().HelixSagas(nameSpace).Update(obj.(*helixsagaoperatorv1.HelixSaga))
	}
	if err != nil {
		klog.V(2).Info(err)
	}
	return res, err
}

func (r *resource) Delete(rt ResourceType, nameSpace, specName string) (err error) {
	var opt Option
	var delOpts = &metav1.DeleteOptions{}
	switch rt {
	case ConfigMap:
		err = r.kubeClientSet.CoreV1().ConfigMaps(nameSpace).Delete(specName, delOpts)
	case NameSpace:
		err = r.kubeClientSet.CoreV1().Namespaces().Delete(specName, delOpts)
	case Service:
		err = r.kubeClientSet.CoreV1().Services(nameSpace).Delete(specName, delOpts)
	case MysqlOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		err = opt.Get().(*mysqlclientset.Clientset).NevercaseV1().MysqlOperators(nameSpace).Delete(specName, delOpts)
	case RedisOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		err = opt.Get().(*redisclientset.Clientset).NevercaseV1().RedisOperators(nameSpace).Delete(specName, delOpts)
	case HelixSagaOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		err = opt.Get().(*helixsagaclientset.Clientset).NevercaseV1().HelixSagas(nameSpace).Delete(specName, delOpts)
	}
	if err != nil {
		klog.V(2).Info(err)
	}
	return err
}

func (r *resource) Get(rt ResourceType, nameSpace, specName string) (res interface{}, err error) {
	var opt Option
	var getOpts = metav1.GetOptions{}
	switch rt {
	case ConfigMap:
		res, err = r.kubeClientSet.CoreV1().ConfigMaps(nameSpace).Get(specName, getOpts)
	case NameSpace:
		res, err = r.kubeClientSet.CoreV1().Namespaces().Get(specName, getOpts)
	case Pod:
		res, err = r.kubeClientSet.CoreV1().Pods(nameSpace).Get(specName, getOpts)
	case Service:
		res, err = r.kubeClientSet.CoreV1().Services(nameSpace).Get(specName, getOpts)
	case Secret:
		res, err = r.kubeClientSet.CoreV1().Secrets(nameSpace).Get(specName, getOpts)
	case ServiceAccount:
		res, err = r.kubeClientSet.CoreV1().ServiceAccounts(nameSpace).Get(specName, getOpts)
	case MysqlOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*mysqlclientset.Clientset).NevercaseV1().MysqlOperators(nameSpace).Get(specName, getOpts)
	case RedisOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*redisclientset.Clientset).NevercaseV1().RedisOperators(nameSpace).Get(specName, getOpts)
	case HelixSagaOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*helixsagaclientset.Clientset).NevercaseV1().HelixSagas(nameSpace).Get(specName, getOpts)
	}
	if err != nil {
		klog.V(2).Info(err)
	}
	return res, err
}

func (r *resource) List(rt ResourceType, nameSpace string, selector labels.Selector) (res interface{}, err error) {
	var opt Option
	var opts = metav1.ListOptions{
		LabelSelector: selector.String(),
	}
	switch rt {
	case ConfigMap:
		res, err = r.kubeClientSet.CoreV1().ConfigMaps(nameSpace).List(opts)
	case Node:
		res, err = r.kubeClientSet.CoreV1().Nodes().List(opts)
	case NameSpace:
		res, err = r.kubeClientSet.CoreV1().Namespaces().List(opts)
	case Pod:
		res, err = r.kubeClientSet.CoreV1().Pods(nameSpace).List(opts)
	case Service:
		res, err = r.kubeClientSet.CoreV1().Services(nameSpace).List(opts)
	case Secret:
		res, err = r.kubeClientSet.CoreV1().Secrets(nameSpace).List(opts)
	case ServiceAccount:
		res, err = r.kubeClientSet.CoreV1().ServiceAccounts(nameSpace).List(opts)
	case MysqlOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*mysqlclientset.Clientset).NevercaseV1().MysqlOperators(nameSpace).List(opts)
	case RedisOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*redisclientset.Clientset).NevercaseV1().RedisOperators(nameSpace).List(opts)
	case HelixSagaOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*helixsagaclientset.Clientset).NevercaseV1().HelixSagas(nameSpace).List(opts)
	}
	if err != nil {
		klog.V(2).Info(err)
	}
	return res, err
}

func (r *resource) Watch(rt ResourceType, nameSpace string, selector labels.Selector, eventsChan chan watch.Event) (err error) {
	var opt Option
	timeout := int64(3600)
	var opts = metav1.ListOptions{
		LabelSelector:  selector.String(),
		TimeoutSeconds: &timeout,
	}
	var res watch.Interface
	switch rt {
	case ConfigMap:
		res, err = r.kubeClientSet.CoreV1().ConfigMaps(nameSpace).Watch(opts)
	case Node:
		res, err = r.kubeClientSet.CoreV1().Nodes().Watch(opts)
	case NameSpace:
		res, err = r.kubeClientSet.CoreV1().Namespaces().Watch(opts)
	case Pod:
		res, err = r.kubeClientSet.CoreV1().Pods(nameSpace).Watch(opts)
	case Service:
		res, err = r.kubeClientSet.CoreV1().Services(nameSpace).Watch(opts)
	case Secret:
		res, err = r.kubeClientSet.CoreV1().Secrets(nameSpace).Watch(opts)
	case MysqlOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*mysqlclientset.Clientset).NevercaseV1().MysqlOperators(nameSpace).Watch(opts)
	case RedisOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*redisclientset.Clientset).NevercaseV1().RedisOperators(nameSpace).Watch(opts)
	case HelixSagaOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*helixsagaclientset.Clientset).NevercaseV1().HelixSagas(nameSpace).Watch(opts)
	}
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	go func() {
		defer func() {
			if err := r.Watch(rt, nameSpace, selector, eventsChan); err != nil {
				klog.Fatal(err)
			}
		}()
		for {
			select {
			case e, isClosed := <-res.ResultChan():
				klog.Infof("resource watch resourceType:%v obj:%v", rt, e)
				if !isClosed {
					klog.Infof("resource watch resourceType:%v closed", rt)
					res.Stop()
					return
				}
				eventsChan <- e
			case <-r.ctx.Done():
				res.Stop()
				return
			}
		}
	}()
	return nil
}

func (r *resource) ResourceTypes() []ResourceType {
	return r.options.GetOptionTypeList()
}
