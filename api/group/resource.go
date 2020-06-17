package group

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"

	mysqlClientSet "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/mysqloperator/clientset/versioned"
	redisClientSet "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/clientset/versioned"
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
	mysql, err := mysqlClientSet.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building mysqlClientSet: %s", err.Error())
	}
	redis, err := redisClientSet.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building redisClientSet: %s", err.Error())
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
	opts := metav1.ListOptions{
		LabelSelector: selector.String(),
	}
	switch rt {
	case ConfigMap:
		res, err = r.kubeClientSet.CoreV1().ConfigMaps(nameSpace).List(opts)
		if err != nil {
			klog.V(2).Info(err)
		}
	case MysqlOperator:
		obj, err := r.options.Get(rt)
		if err != nil {
			klog.V(2).Info(err)
		}
		mysql := obj.Get().(*mysqlClientSet.Clientset)
		res, err = mysql.MysqloperatorV1().MysqlOperators(nameSpace).List(opts)
		if err != nil {
			klog.V(2).Info(err)
		}
	case RedisOperator:
		obj, err := r.options.Get(rt)
		if err != nil {
			klog.V(2).Info(err)
		}
		redis := obj.Get().(*redisClientSet.Clientset)
		res, err = redis.RedisoperatorV1().RedisOperators(nameSpace).List(opts)
		if err != nil {
			klog.V(2).Info(err)
		}
	case HelixOperator:
	}
	return res, err
}
