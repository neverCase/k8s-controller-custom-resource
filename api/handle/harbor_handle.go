package handle

import (
	harbor "github.com/nevercase/harbor-api"
	"github.com/nevercase/k8s-controller-custom-resource/api/group"
	"github.com/nevercase/k8s-controller-custom-resource/api/proto"
	"k8s.io/klog"
)

type HarborApiGetter interface {
	HarborApi() HarborApiInterface
}

type HarborApiInterface interface {
	Projects() (res []byte, err error)
	Repositories(projectId int) (res []byte, err error)
	Tags(imageName string) (res []byte, err error)
}

func NewHarborApi(g group.Group) HarborApiInterface {
	return &harborApi{
		group: g,
	}
}

type harborApi struct {
	group group.Group
}

func (ha *harborApi) Projects() (res []byte, err error) {
	t := make([]harbor.Project, 0)
	if t, err = ha.group.Harbor().Projects(); err != nil {
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
	return m.Marshal()
}

func (ha *harborApi) Repositories(projectId int) (res []byte, err error) {
	t := make([]harbor.RepoRecord, 0)
	if t, err = ha.group.Harbor().Repositories(projectId); err != nil {
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
	return m.Marshal()
}

func (ha *harborApi) Tags(imageName string) (res []byte, err error) {
	t := make([]harbor.TagDetail, 0)
	if t, err = ha.group.Harbor().Tags(imageName); err != nil {
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
	return m.Marshal()
}
