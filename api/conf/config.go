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

type Config interface {
	MasterUrl() string
	KubeConfig() string
	ApiService() string
}

type config struct {
	masterUrl  string
	kubeConfig string
	apiService string
}

func (c *config) MasterUrl() string {
	return c.masterUrl
}

func (c *config) KubeConfig() string {
	return c.kubeConfig
}

func (c *config) ApiService() string {
	return c.apiService
}

func Init() Config {
	return &config{
		masterUrl:  masterUrl,
		kubeConfig: kubeconfig,
		apiService: apiservice,
	}
}
