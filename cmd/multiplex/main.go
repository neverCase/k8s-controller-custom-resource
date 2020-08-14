package main

import (
	"flag"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/Shanghai-Lunara/helixsaga-operator/pkg/controllers/helixsaga"
	k8sCoreV1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
	mysql "github.com/nevercase/k8s-controller-custom-resource/pkg/controller/mysqloperator"
	redis "github.com/nevercase/k8s-controller-custom-resource/pkg/controller/redisoperator"

	"github.com/nevercase/k8s-controller-custom-resource/pkg/signals"
)

var (
	masterURL  string
	kubeconfig string
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "masterurl", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	k8sClientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientSet: %s", err.Error())
	}

	controllerName := "multiplex-controller"
	opts := k8sCoreV1.NewOptions()
	mysqlOpt := mysql.NewOption(controllerName, cfg, stopCh)
	redisOpt := redis.NewOption(controllerName, cfg, stopCh)
	helixSagaOpt := helixsaga.NewOption(controllerName, cfg, stopCh)
	if err := opts.Add(mysqlOpt, redisOpt, helixSagaOpt); err != nil {
		klog.Fatal(err)
	}

	operator := k8sCoreV1.NewKubernetesOperator(k8sClientSet, stopCh, controllerName, opts)
	kc := k8sCoreV1.NewKubernetesController(operator)
	if err = kc.Run(4, stopCh); err != nil {
		klog.Fatalf("Error running multiplex-controller: %s", err.Error())
	}
}
