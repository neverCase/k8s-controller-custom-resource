package main

import (
	"flag"

	"github.com/nevercase/k8s-controller-custom-resource/api/conf"
	"github.com/nevercase/k8s-controller-custom-resource/api/service"
	"github.com/nevercase/k8s-controller-custom-resource/api/v1"
	"github.com/nevercase/k8s-controller-custom-resource/pkg/signals"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	config := conf.Init()

	cfg, err := clientcmd.BuildConfigFromFlags(config.MasterUrl(), config.KubeConfig())
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	g := v1.NewGroup(cfg)
	_ = g

	s := service.NewService(config)
	s.Listen()

	//m, err := g.Mysql().MysqloperatorV1().MysqlOperators(apiV1.NamespaceDefault).Create(newMysql())
	//if err != nil {
	//	klog.V(2).Info(err)
	//}
	//klog.V(4).Info(m)

	<-stopCh

	//err = g.Mysql().MysqloperatorV1().MysqlOperators(apiV1.NamespaceDefault).Delete("example-mysql", &metaV1.DeleteOptions{})
	//if err != nil {
	//	klog.V(2).Info(err)
	//}
}

//func newMysql() *mysqlOperatorV1.MysqlOperator {
//	var a int32 = 1
//	var b int32 = 4
//	return &mysqlOperatorV1.MysqlOperator{
//		ObjectMeta: metaV1.ObjectMeta{
//			Name:      "example-mysql",
//			Namespace: apiV1.NamespaceDefault,
//		},
//		Spec: mysqlOperatorV1.MysqlOperatorSpec{
//			MasterSpec: mysqlOperatorV1.MysqlDeploymentSpec{
//				DeploymentName:   "test-mysql",
//				Replicas:         &a,
//				Image:            "domain/mysql-slave:1.0",
//				ImagePullSecrets: "private-harbor",
//				Configuration:    mysqlOperatorV1.MysqlConfig{},
//			},
//			SlaveSpec: mysqlOperatorV1.MysqlDeploymentSpec{
//				DeploymentName:   "test-mysql",
//				Replicas:         &b,
//				Image:            "domain/mysql-slave:1.0",
//				ImagePullSecrets: "private-harbor",
//				Configuration:    mysqlOperatorV1.MysqlConfig{},
//			},
//		},
//	}
//}
