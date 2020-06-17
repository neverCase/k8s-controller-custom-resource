package service

import (
	"github.com/nevercase/k8s-controller-custom-resource/api/group"
	"github.com/nevercase/k8s-controller-custom-resource/api/proto"
	"k8s.io/apimachinery/pkg/labels"

	mysqloperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/mysqloperator/v1"
	redisoperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/redisoperator/v1"
)

type HandleInterface interface {
	List(param proto.Param) ([]byte, error)
}

func NewHandle(g group.Group) HandleInterface {
	return &handle{
		group: g,
	}
}

type handle struct {
	group group.Group
}

func (h *handle) List(req proto.Param) (res []byte, err error) {
	var o []byte
	var d interface{}
	var selector = labels.NewSelector()
	if d, err = h.group.Resource().List(req.ResourceType, req.NameSpace, selector); err != nil {
		return res, err
	}
	switch req.ResourceType {
	case group.MysqlOperator:
		m := proto.MysqlCrdList{
			Items: make([]proto.MysqlCrd, 0),
		}
		for _, v := range d.(*mysqloperatorv1.MysqlOperatorList).Items {
			m.Items = append(m.Items, proto.MysqlCrd{
				Name: v.Name,
				Master: proto.NodeSpec{
					Name:             v.Spec.MasterSpec.Spec.Name,
					Replicas:         *v.Spec.MasterSpec.Spec.Replicas,
					Image:            v.Spec.MasterSpec.Spec.Image,
					ImagePullSecrets: v.Spec.MasterSpec.Spec.ImagePullSecrets[0].Name,
				},
				Slave: proto.NodeSpec{
					Name:             v.Spec.SlaveSpec.Spec.Name,
					Replicas:         *v.Spec.SlaveSpec.Spec.Replicas,
					Image:            v.Spec.SlaveSpec.Spec.Image,
					ImagePullSecrets: v.Spec.SlaveSpec.Spec.ImagePullSecrets[0].Name,
				},
			})
		}
		o, err = m.Marshal()
	case group.RedisOperator:
		m := proto.RedisCrdList{
			Items: make([]proto.RedisCrd, 0),
		}
		for _, v := range d.(*redisoperatorv1.RedisOperatorList).Items {
			m.Items = append(m.Items, proto.RedisCrd{
				Name: v.Name,
				Master: proto.NodeSpec{
					Name:             v.Spec.MasterSpec.Spec.Name,
					Replicas:         *v.Spec.MasterSpec.Spec.Replicas,
					Image:            v.Spec.MasterSpec.Spec.Image,
					ImagePullSecrets: v.Spec.MasterSpec.Spec.ImagePullSecrets[0].Name,
				},
				Slave: proto.NodeSpec{
					Name:             v.Spec.SlaveSpec.Spec.Name,
					Replicas:         *v.Spec.SlaveSpec.Spec.Replicas,
					Image:            v.Spec.SlaveSpec.Spec.Image,
					ImagePullSecrets: v.Spec.SlaveSpec.Spec.ImagePullSecrets[0].Name,
				},
			})
		}
		o, err = m.Marshal()
	case group.HelixOperator:
	}
	if err != nil {
		return nil, err
	}
	return proto.GetResponse(req, string(o))
}
