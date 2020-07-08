package group

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog"
)

type Harbor interface {
	Http(method string, url string) (res *http.Response, err error)
	Login() error
	Projects() (res []Project, err error)
	Repositories(projectId int) (res []RepoRecord, err error)
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
	Login        HarborUrlSuffix = "login"
	SystemInfo   HarborUrlSuffix = "api/systeminfo"
	Projects     HarborUrlSuffix = "api/projects"                    // api/projects?page=1&page_size=15
	Repositories HarborUrlSuffix = "api/repositories?&project_id=%d" // api/repositories?page=1&page_size=15&project_id=2
)

func (h *harbor) SystemInfo() {

}

func (h *harbor) Http(method string, url string) (res *http.Response, err error) {
	var req *http.Request
	if req, err = http.NewRequest(method, url, nil); err != nil {
		klog.V(2).Info(err)
		return res, err
	}
	req.SetBasicAuth(h.admin, h.password)
	httpClient := http.Client{
		Timeout: time.Second * time.Duration(h.timeout),
	}
	if res, err = httpClient.Do(req); err != nil {
		klog.V(2).Info(err)
	}
	return res, err
}

func (h *harbor) Login() error {
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

func (h *harbor) Projects() (res []Project, err error) {
	var resp *http.Response
	if resp, err = h.Http("GET", fmt.Sprintf("%s/%v", h.url, Projects)); err != nil {
		return res, err
	}
	if resp.StatusCode == http.StatusOK {
		cont, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			klog.V(2).Info(err)
			return res, err
		}
		if err = resp.Body.Close(); err != nil {
			klog.V(2).Info(err)
			return res, err
		}
		if err = json.Unmarshal(cont, &res); err != nil {
			klog.V(2).Info(err)
			return res, err
		}
	}
	klog.Info(res)
	return res, nil
}

func (h *harbor) Repositories(projectId int) (res []RepoRecord, err error) {
	var (
		suffix string
		resp   *http.Response
	)
	suffix = fmt.Sprintf(string(Repositories), projectId)
	if resp, err = h.Http("GET", fmt.Sprintf("%s/%v", h.url, suffix)); err != nil {
		return res, err
	}
	if resp.StatusCode == http.StatusOK {
		cont, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			klog.V(2).Info(err)
			return res, err
		}
		if err = resp.Body.Close(); err != nil {
			klog.V(2).Info(err)
			return res, err
		}
		if err = json.Unmarshal(cont, &res); err != nil {
			klog.V(2).Info(err)
			return res, err
		}
	}
	klog.Info(res)
	return res, nil
}
