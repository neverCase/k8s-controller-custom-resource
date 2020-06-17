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
	List(rt ResourceType, nameSpace string, selector labels.Selector) (res interface{}, err error)
	Resource() []ResourceType
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

func (r *resource) Resource() []ResourceType {
	return r.options.GetOptionTypeList()
}
