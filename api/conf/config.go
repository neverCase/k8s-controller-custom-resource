package conf

import (
	"flag"
	harbor "github.com/nevercase/harbor-api"
	"k8s.io/klog/v2"
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
	masterUrl      string
	kubeconfig     string
	apiservice     string
	dockerUrl      arrayFlags
	dockerAdmin    arrayFlags
	dockerPassword arrayFlags
	rbacRulePath   string
	rbacMysqlPath  string
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterUrl, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&apiservice, "apiservice", "0.0.0.0:9090", "The address of the api server.")
	flag.Var(&dockerUrl, "dockerurl", "The address of the Harbor server.")
	flag.Var(&dockerAdmin, "dockeradmin", "The username of the Harbor's account")
	flag.Var(&dockerPassword, "dockerpwd", "The password of the Harbor's password")
	flag.StringVar(&rbacRulePath, "rbacRulePath", "", "The path of the rbac rule.")
	flag.StringVar(&rbacMysqlPath, "rbacMysqlPath", "", "The path of the rbac mysql.")
}

type Config interface {
	MasterUrl() string
	KubeConfig() string
	ApiService() string
	DockerHub() []harbor.Config
	RbacRulePath() string
	RbacMysqlPath() string
}

type config struct {
	masterUrl     string
	kubeConfig    string
	apiService    string
	dockerHub     []harbor.Config
	rbacRulePath  string
	rbacMysqlPath string
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

func (c *config) DockerHub() []harbor.Config {
	return c.dockerHub
}

func (c *config) RbacRulePath() string {
	return c.rbacRulePath
}
func (c *config) RbacMysqlPath() string {
	return c.rbacMysqlPath
}

func Init() Config {
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
	klog.Info("dockerhub:", dockerHub)
	return &config{
		masterUrl:     masterUrl,
		kubeConfig:    kubeconfig,
		apiService:    apiservice,
		dockerHub:     dockerHub,
		rbacRulePath:  rbacRulePath,
		rbacMysqlPath: rbacMysqlPath,
	}
}
