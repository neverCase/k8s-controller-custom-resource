package service

import (
	"github.com/nevercase/k8s-controller-custom-resource/api/group"
	"github.com/nevercase/k8s-controller-custom-resource/api/proto"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	mysqloperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/mysqloperator/v1"
	redisoperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/redisoperator/v1"
)

type HandleInterface interface {
	List(req proto.Param) ([]byte, error)
	Resources(req proto.Param) (res []byte, err error)
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

func (h *handle) Resources(req proto.Param) (res []byte, err error) {
	m := proto.Resources{
		Resources: h.group.Resource().ResourceTypes(),
	}
	o, err := m.Marshal()
	if err != nil {
		return nil, err
	}
	return proto.GetResponse(req, string(o))
}

func (h *handle) Create(req proto.Param, obj interface{}) (res []byte, err error) {
	switch req.ResourceType {
	case group.MysqlOperator:
		mysqlCrd := obj.(proto.MysqlCrd)
		m := &mysqloperatorv1.MysqlOperator{
			ObjectMeta: metav1.ObjectMeta{
				Name:      mysqlCrd.Name,
				Namespace: req.NameSpace,
			},
			Spec: mysqloperatorv1.MysqlOperatorSpec{
				MasterSpec: mysqloperatorv1.MysqlCore{
					Spec: mysqloperatorv1.MysqlSpec{
						Name:     mysqlCrd.Master.Name,
						Replicas: &mysqlCrd.Master.Replicas,
						Image:    mysqlCrd.Master.Image,
						ImagePullSecrets: []corev1.LocalObjectReference{
							{
								Name: mysqlCrd.Master.ImagePullSecrets,
							},
						},
					},
				},
				SlaveSpec: mysqloperatorv1.MysqlCore{
					Spec: mysqloperatorv1.MysqlSpec{
						Name:     mysqlCrd.Slave.Name,
						Replicas: &mysqlCrd.Slave.Replicas,
						Image:    mysqlCrd.Slave.Image,
						ImagePullSecrets: []corev1.LocalObjectReference{
							{
								Name: mysqlCrd.Slave.ImagePullSecrets,
							},
						},
					},
				},
			},
		}
		n, err := h.group.Resource().Create(req.ResourceType, req.NameSpace, m)
		if err != nil {
			break
		}
		v := n.(mysqloperatorv1.MysqlOperator)
		e := proto.MysqlCrd{
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
		}
		res, err = e.Marshal()
	case group.RedisOperator:
		redisCrd := obj.(proto.RedisCrd)
		m := &redisoperatorv1.RedisOperator{
			ObjectMeta: metav1.ObjectMeta{
				Name:      redisCrd.Name,
				Namespace: req.NameSpace,
			},
			Spec: redisoperatorv1.RedisOperatorSpec{
				MasterSpec: redisoperatorv1.RedisCore{
					Spec: redisoperatorv1.RedisSpec{
						Name:     redisCrd.Master.Name,
						Replicas: &redisCrd.Master.Replicas,
						Image:    redisCrd.Master.Image,
						ImagePullSecrets: []corev1.LocalObjectReference{
							{
								Name: redisCrd.Master.ImagePullSecrets,
							},
						},
					},
				},
				SlaveSpec: redisoperatorv1.RedisCore{
					Spec: redisoperatorv1.RedisSpec{
						Name:     redisCrd.Slave.Name,
						Replicas: &redisCrd.Slave.Replicas,
						Image:    redisCrd.Slave.Image,
						ImagePullSecrets: []corev1.LocalObjectReference{
							{
								Name: redisCrd.Slave.ImagePullSecrets,
							},
						},
					},
				},
			},
		}
		n, err := h.group.Resource().Create(req.ResourceType, req.NameSpace, m)
		if err != nil {
			break
		}
		v := n.(redisoperatorv1.RedisOperator)
		e := proto.RedisCrd{
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
		}
		res, err = e.Marshal()
	case group.HelixOperator:
	}
	return res, err
}
