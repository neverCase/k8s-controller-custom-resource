package handle

import (
	"fmt"
	harbor "github.com/nevercase/harbor-api"
	"github.com/nevercase/k8s-controller-custom-resource/api/group"
	"github.com/nevercase/k8s-controller-custom-resource/api/proto"
	"k8s.io/klog/v2"
	"sort"
)

type HarborApiGetter interface {
	HarborApi() HarborApiInterface
}

type HarborApiInterface interface {
	Core(req proto.Param, obj []byte) (res []byte, err error)
	Hubs() (res []byte, err error)
	Projects(url string) (res []byte, err error)
	Repositories(url string, projectId int) (res []byte, err error)
	Tags(url, imageName string) (res []byte, err error)
}

func NewHarborApi(g group.Group) HarborApiInterface {
	return &harborApi{
		group: g,
	}
}

type harborApi struct {
	group group.Group
}

func (ha *harborApi) Core(req proto.Param, obj []byte) (res []byte, err error) {
	if req.HarborRequest.Command != proto.Hubs && req.HarborRequest.HarborUrl == "" {
		err = fmt.Errorf("no hr.HarborUrl")
		klog.V(2).Info(err)
		return nil, err
	}
	switch req.HarborRequest.Command {
	case proto.Hubs:
		res, err = ha.Hubs()
	case proto.Projects:
		res, err = ha.Projects(req.HarborRequest.HarborUrl)
	case proto.Repositories:
		res, err = ha.Repositories(req.HarborRequest.HarborUrl, int(req.HarborRequest.ProjectID))
	case proto.Tags:
		res, err = ha.Tags(req.HarborRequest.HarborUrl, req.HarborRequest.ImageName)
	}
	if err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	return proto.GetResponse(req, res)
}

func (ha *harborApi) Hubs() (res []byte, err error) {
	t := make([]string, 0)
	t = ha.group.HarborHub().List()
	m := proto.HarborHubList{
		Items: make([]proto.HarborHub, 0),
	}
	for _, v := range t {
		m.Items = append(m.Items, proto.HarborHub{
			Name: v,
		})
	}
	return m.Marshal()
}

func (ha *harborApi) Projects(url string) (res []byte, err error) {
	h, err := ha.group.HarborHub().Get(url)
	if err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	t := make([]harbor.Project, 0)
	if t, err = h.Projects(); err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	m := proto.HarborProjectList{
		Items: make([]proto.HarborProject, 0),
	}
	for _, v := range t {
		m.Items = append(m.Items, proto.HarborProject{
			Name:      v.Name,
			ProjectID: int32(v.ProjectID),
		})
	}
	sort.Slice(m.Items, func(i, j int) bool {
		if m.Items[i].Name > m.Items[j].Name {
			return true
		}
		return false
	})
	return m.Marshal()
}

func (ha *harborApi) Repositories(url string, projectId int) (res []byte, err error) {
	h, err := ha.group.HarborHub().Get(url)
	if err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	t := make([]harbor.RepoRecord, 0)
	if t, err = h.Repositories(projectId); err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	m := proto.HarborRepositoryList{
		Items: make([]proto.HarborRepository, 0),
	}
	for _, v := range t {
		m.Items = append(m.Items, proto.HarborRepository{
			RepositoryID: int32(v.RepositoryID),
			Name:         v.Name,
			ProjectID:    int32(v.ProjectID),
		})
	}
	sort.Slice(m.Items, func(i, j int) bool {
		if m.Items[i].Name > m.Items[j].Name {
			return true
		}
		return false
	})
	return m.Marshal()
}

func (ha *harborApi) Tags(url, imageName string) (res []byte, err error) {
	h, err := ha.group.HarborHub().Get(url)
	if err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	t := make([]harbor.TagDetail, 0)
	if t, err = h.Tags(imageName); err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	m := proto.HarborTagList{
		Items: make([]proto.HarborTag, 0),
	}
	for _, v := range t {
		m.Items = append(m.Items, proto.HarborTag{
			Digest: v.Digest,
			Name:   v.Name,
		})
	}
	sort.Slice(m.Items, func(i, j int) bool {
		if m.Items[i].Name == "latest" {
			return true
		}
		if m.Items[j].Name == "latest" {
			return false
		}
		if m.Items[i].Name > m.Items[j].Name {
			return true
		}
		return false
	})
	return m.Marshal()
}
