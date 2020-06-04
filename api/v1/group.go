package v1

import (
	"k8s.io/client-go/rest"
	"k8s.io/klog"

	mysqlClientSet "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/mysqloperator/clientset/versioned"
	redisClientSet "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/clientset/versioned"
)

type Group interface {
	Mysql() *mysqlClientSet.Clientset
	Redis() *redisClientSet.Clientset
}

type group struct {
	mysql *mysqlClientSet.Clientset
	redis *redisClientSet.Clientset
}

func NewGroup(cfg *rest.Config) Group {
	mysql, err := mysqlClientSet.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building example mysqlClientSet: %s", err.Error())
	}
	redis, err := redisClientSet.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building example mysqlClientSet: %s", err.Error())
	}
	var g = &group{
		mysql: mysql,
		redis: redis,
	}
	return g
}

func (g *group) Mysql() *mysqlClientSet.Clientset {
	return g.mysql
}

func (g *group) Redis() *redisClientSet.Clientset {
	return g.redis
}
