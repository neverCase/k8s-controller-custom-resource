package conf

import (
	"flag"
)

var (
	masterUrl      string
	kubeconfig     string
	apiservice     string
	dockerUrl      string
	dockerAdmin    string
	dockerPassword string
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterUrl, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&apiservice, "apiservice", "0.0.0.0:9090", "The address of the api server.")
	flag.StringVar(&dockerUrl, "dockerurl", "", "The address of the Harbor server.")
	flag.StringVar(&dockerAdmin, "dockeradmin", "", "The username of the Harbor's accoount")
	flag.StringVar(&dockerPassword, "dockerpwd", "", "The password of the Harbor's password")
}

type Config interface {
	MasterUrl() string
	KubeConfig() string
	ApiService() string
	DockerUrl() string
	DockerAdmin() string
	DockerPassword() string
}

type config struct {
	masterUrl      string
	kubeConfig     string
	apiService     string
	dockerUrl      string
	dockerAdmin    string
	dockerPassword string
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

func (c *config) DockerUrl() string {
	return c.dockerUrl
}

func (c *config) DockerAdmin() string {
	return c.dockerAdmin
}

func (c *config) DockerPassword() string {
	return c.dockerPassword
}

func Init() Config {
	return &config{
		masterUrl:      masterUrl,
		kubeConfig:     kubeconfig,
		apiService:     apiservice,
		dockerUrl:      dockerUrl,
		dockerAdmin:    dockerAdmin,
		dockerPassword: dockerPassword,
	}
}
