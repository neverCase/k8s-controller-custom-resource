package group

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"

	mysqloperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/mysqloperator/v1"
	redisoperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/redisoperator/v1"
	mysqlclientset "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/mysqloperator/clientset/versioned"
	redisclientset "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/clientset/versioned"
)

type ResourceType string

const (
	ConfigMap   ResourceType = "ConfigMap"
	Deployment  ResourceType = "Deployment"
	Pod         ResourceType = "Pod"
	Secret      ResourceType = "Secret"
	Service     ResourceType = "Service"
	StatefulSet ResourceType = "StatefulSet"

	MysqlOperator ResourceType = "MysqlOperator"
	RedisOperator ResourceType = "RedisOperator"
	HelixOperator ResourceType = "HelixOperator"
)

type ResourceGetter interface {
	Resource() ResourceInterface
}

type ResourceInterface interface {
	Create(rt ResourceType, nameSpace string, obj interface{}) (res interface{}, err error)
	Update(rt ResourceType, nameSpace string, obj interface{}) (res interface{}, err error)
	Delete(rt ResourceType, nameSpace, specName string) (err error)
	Get(rt ResourceType, nameSpace, specName string) (res interface{}, err error)
	List(rt ResourceType, nameSpace string, selector labels.Selector) (res interface{}, err error)
	ResourceTypes() []ResourceType
}

func NewResource(masterUrl, kubeconfigPath string) ResourceInterface {
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
	opts := NewOptions()
	mysqlOpt := NewOption(MysqlOperator, mysql)
	redisOpt := NewOption(RedisOperator, redis)
	opts.Add(mysqlOpt, redisOpt)
	var r = &resource{
		kubeClientSet: kubeClient,
		options:       opts,
	}
	return r
}

type resource struct {
	kubeClientSet kubernetes.Interface
	options       Options
}

func (r *resource) Create(rt ResourceType, nameSpace string, obj interface{}) (res interface{}, err error) {
	var opt Option
	switch rt {
	case ConfigMap:
		res, err = r.kubeClientSet.CoreV1().ConfigMaps(nameSpace).Create(obj.(*corev1.ConfigMap))
	case MysqlOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*mysqlclientset.Clientset).MysqloperatorV1().MysqlOperators(nameSpace).Create(obj.(*mysqloperatorv1.MysqlOperator))
	case RedisOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*redisclientset.Clientset).RedisoperatorV1().RedisOperators(nameSpace).Create(obj.(*redisoperatorv1.RedisOperator))
	case HelixOperator:
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
	case MysqlOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*mysqlclientset.Clientset).MysqloperatorV1().MysqlOperators(nameSpace).Update(obj.(*mysqloperatorv1.MysqlOperator))
	case RedisOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*redisclientset.Clientset).RedisoperatorV1().RedisOperators(nameSpace).Update(obj.(*redisoperatorv1.RedisOperator))
	case HelixOperator:
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
	case MysqlOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		err = opt.Get().(*mysqlclientset.Clientset).MysqloperatorV1().MysqlOperators(nameSpace).Delete(specName, delOpts)
	case RedisOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		err = opt.Get().(*redisclientset.Clientset).RedisoperatorV1().RedisOperators(nameSpace).Delete(specName, delOpts)
	case HelixOperator:
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
	case MysqlOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*mysqlclientset.Clientset).MysqloperatorV1().MysqlOperators(nameSpace).Get(specName, getOpts)
	case RedisOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*redisclientset.Clientset).RedisoperatorV1().RedisOperators(nameSpace).Get(specName, getOpts)
	case HelixOperator:
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
	case MysqlOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*mysqlclientset.Clientset).MysqloperatorV1().MysqlOperators(nameSpace).List(opts)
	case RedisOperator:
		if opt, err = r.options.Get(rt); err != nil {
			break
		}
		res, err = opt.Get().(*redisclientset.Clientset).RedisoperatorV1().RedisOperators(nameSpace).List(opts)
	case HelixOperator:
	}
	if err != nil {
		klog.V(2).Info(err)
	}
	return res, err
}

func (r *resource) ResourceTypes() []ResourceType {
	return r.options.GetOptionTypeList()
}