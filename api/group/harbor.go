package group

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"k8s.io/klog"
)

type Harbor interface {
	Login() error
	Projects() error
}

func NewHarbor(url, admin, password string) Harbor {
	return &harbor{
		url:      url,
		admin:    admin,
		password: password,
		timeout:  10,
		cookie:   make([]*http.Cookie, 0),
	}
}

type harbor struct {
	url      string
	admin    string
	password string
	timeout  int

	cookie        []*http.Cookie
	cookieTimeout time.Time
}

type HarborUrlSuffix string

const (
	Login      HarborUrlSuffix = "login"
	SystemInfo HarborUrlSuffix = "api/systeminfo"
	Projects   HarborUrlSuffix = "api/projects"
)

func (h *harbor) SystemInfo() {

}

func (h *harbor) Login() error {
	fmt.Println("h:", h)
	var (
		req  *http.Request
		resp *http.Response
		err  error
	)
	req, err = http.NewRequest("GET", fmt.Sprintf("%s/%v", h.url, Login), nil)
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	req.SetBasicAuth(h.admin, h.password)
	httpClient := http.Client{
		Timeout: time.Second * time.Duration(h.timeout),
	}
	resp, err = httpClient.Do(req)
	if err != nil {
		klog.V(2).Info(err)
	}
	h.cookie = resp.Cookies()
	return err
}

func (h *harbor) Projects() error {
	var (
		req  *http.Request
		resp *http.Response
		err  error
	)
	req, err = http.NewRequest("GET", fmt.Sprintf("%s/%v", h.url, Projects), nil)
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	req.SetBasicAuth(h.admin, h.password)
	httpClient := http.Client{
		Timeout: time.Second * time.Duration(h.timeout),
	}
	resp, err = httpClient.Do(req)
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	if resp.StatusCode == http.StatusOK {
		cont, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			klog.V(2).Info(err)
			return err
		}
		if err = resp.Body.Close(); err != nil {
			klog.V(2).Info(err)
			return err
		}
		fmt.Println("cont:", string(cont))
	}
	return err
}
