package conf

import (
	"flag"
)

var (
	masterUrl  string
	kubeconfig string
	apiservice string
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterUrl, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&apiservice, "apiservice", "0.0.0.0:9090", "The address of the api server.")
}

type Config struct {
	MasterUrl   string
	KubeConfig  string
	ApiService string
}

func Init() *Config {
	return &Config{
		MasterUrl:   masterUrl,
		KubeConfig:  kubeconfig,
		ApiService: apiservice,
	}
}
