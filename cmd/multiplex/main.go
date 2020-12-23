package main

import (
	"flag"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	"github.com/Shanghai-Lunara/helixsaga-operator/pkg/controllers/helixsaga"
	harbor "github.com/nevercase/harbor-api"
	k8sCoreV1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
	mysql "github.com/nevercase/k8s-controller-custom-resource/pkg/controller/mysqloperator"
	redis "github.com/nevercase/k8s-controller-custom-resource/pkg/controller/redisoperator"
	"github.com/nevercase/k8s-controller-custom-resource/pkg/signals"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	masterURL      string
	kubeconfig     string
	dockerUrl      arrayFlags
	dockerAdmin    arrayFlags
	dockerPassword arrayFlags
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "masterurl", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.Var(&dockerUrl, "dockerurl", "The address of the Harbor server.")
	flag.Var(&dockerAdmin, "dockeradmin", "The username of the Harbor's account")
	flag.Var(&dockerPassword, "dockerpwd", "The password of the Harbor's password")
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

	dockerHub := make([]harbor.Config, 0)
	for k, url := range dockerUrl {
		var admin, password string
		if len(dockerAdmin) >= k+1 {
			admin = dockerAdmin[k]
		}
		if len(dockerPassword) >= k+1 {
			password = dockerPassword[k]
		}
		dockerHub = append(dockerHub, harbor.Config{
			Url:      url,
			Admin:    admin,
			Password: password,
		})
	}

	controllerName := "multiplex-controller"
	opts := k8sCoreV1.NewOptions()
	mysqlOpt := mysql.NewOption(controllerName, cfg, stopCh)
	redisOpt := redis.NewOption(controllerName, cfg, stopCh)
	helixSagaOpt := helixsaga.NewOption(controllerName, cfg, stopCh, dockerHub)
	if err := opts.Add(mysqlOpt, redisOpt, helixSagaOpt); err != nil {
		klog.Fatal(err)
	}

	operator := k8sCoreV1.NewKubernetesOperator(k8sClientSet, stopCh, controllerName, opts)
	kc := k8sCoreV1.NewKubernetesController(operator)
	if err = kc.Run(10, stopCh); err != nil {
		klog.Fatalf("Error running multiplex-controller: %s", err.Error())
	}
}
