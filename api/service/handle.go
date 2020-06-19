package service

import (
	"github.com/nevercase/k8s-controller-custom-resource/api/group"
	"github.com/nevercase/k8s-controller-custom-resource/api/proto"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	mysqloperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/mysqloperator/v1"
	redisoperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/redisoperator/v1"
)

type HandleInterface interface {
	Create(req proto.Param, obj []byte) (res []byte, err error)
	Delete(req proto.Param, obj []byte) (err error)
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

func (h *handle) Create(req proto.Param, obj []byte) (res []byte, err error) {
	switch req.ResourceType {
	case group.MysqlOperator:
		var n interface{}
		var mysqlCrd proto.MysqlCrd
		if err = mysqlCrd.Unmarshal(obj); err != nil {
			break
		}
		m := convertMysqlCrdToOperator(req, mysqlCrd)
		_, err = h.group.Resource().Get(req.ResourceType, req.NameSpace, m.Name)
		if err != nil {
			if !errors.IsNotFound(err) {
				break
			}
			err = nil
			if n, err = h.group.Resource().Create(req.ResourceType, req.NameSpace, m); err != nil {
				break
			}
		} else {
			if n, err = h.group.Resource().Update(req.ResourceType, req.NameSpace, m); err != nil {
				break
			}
		}
		v := n.(*mysqloperatorv1.MysqlOperator)
		e := convertOperatorToMysqlCrd(v)
		res, err = e.Marshal()
	case group.RedisOperator:
		var n interface{}
		var redisCrd proto.RedisCrd
		if err = redisCrd.Unmarshal(obj); err != nil {
			break
		}
		m := convertRedisCrdToOperator(req, redisCrd)
		_, err = h.group.Resource().Get(req.ResourceType, req.NameSpace, m.Name)
		if err != nil {
			if !errors.IsNotFound(err) {
				break
			}
			err = nil
			if n, err = h.group.Resource().Create(req.ResourceType, req.NameSpace, m); err != nil {
				break
			}
		} else {
			if n, err = h.group.Resource().Update(req.ResourceType, req.NameSpace, m); err != nil {
				break
			}
		}
		v := n.(*redisoperatorv1.RedisOperator)
		e := convertOperatorToRedisCrd(v)
		res, err = e.Marshal()
	case group.HelixOperator:
	}
	return proto.GetResponse(req, res)
}

func (h *handle) Delete(req proto.Param, obj []byte) (err error) {
	switch req.ResourceType {
	case group.MysqlOperator:
		var mysqlCrd proto.MysqlCrd
		if err = mysqlCrd.Unmarshal(obj); err != nil {
			break
		}
		err = h.group.Resource().Delete(req.ResourceType, req.NameSpace, mysqlCrd.Name)
	case group.RedisOperator:
		var redisCrd proto.RedisCrd
		if err = redisCrd.Unmarshal(obj); err != nil {
			break
		}
		err = h.group.Resource().Delete(req.ResourceType, req.NameSpace, redisCrd.Name)
	case group.HelixOperator:
	}
	return err
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
			m.Items = append(m.Items, convertOperatorToMysqlCrd(&v))
		}
		o, err = m.Marshal()
	case group.RedisOperator:
		m := proto.RedisCrdList{
			Items: make([]proto.RedisCrd, 0),
		}
		for _, v := range d.(*redisoperatorv1.RedisOperatorList).Items {
			m.Items = append(m.Items, convertOperatorToRedisCrd(&v))
		}
		o, err = m.Marshal()
	case group.HelixOperator:
	}
	if err != nil {
		return nil, err
	}
	return proto.GetResponse(req, o)
}

func (h *handle) Resources(req proto.Param) (res []byte, err error) {
	m := proto.Resources{
		Resources: h.group.Resource().ResourceTypes(),
	}
	o, err := m.Marshal()
	if err != nil {
		return nil, err
	}
	return proto.GetResponse(req, o)
}

func convertMysqlCrdToOperator(req proto.Param, mysqlCrd proto.MysqlCrd) *mysqloperatorv1.MysqlOperator {
	return &mysqloperatorv1.MysqlOperator{
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
					VolumePath: mysqlCrd.Master.VolumePath,
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
					VolumePath: mysqlCrd.Slave.VolumePath,
				},
			},
		},
	}
}

func convertOperatorToMysqlCrd(m *mysqloperatorv1.MysqlOperator) proto.MysqlCrd {
	return proto.MysqlCrd{
		Name: m.Name,
		Master: proto.NodeSpec{
			Name:             m.Spec.MasterSpec.Spec.Name,
			Replicas:         *m.Spec.MasterSpec.Spec.Replicas,
			Image:            m.Spec.MasterSpec.Spec.Image,
			ImagePullSecrets: m.Spec.MasterSpec.Spec.ImagePullSecrets[0].Name,
			VolumePath:       m.Spec.MasterSpec.Spec.VolumePath,
		},
		Slave: proto.NodeSpec{
			Name:             m.Spec.SlaveSpec.Spec.Name,
			Replicas:         *m.Spec.SlaveSpec.Spec.Replicas,
			Image:            m.Spec.SlaveSpec.Spec.Image,
			ImagePullSecrets: m.Spec.SlaveSpec.Spec.ImagePullSecrets[0].Name,
			VolumePath:       m.Spec.SlaveSpec.Spec.VolumePath,
		},
	}
}

func convertRedisCrdToOperator(req proto.Param, redisCrd proto.RedisCrd) *redisoperatorv1.RedisOperator {
	return &redisoperatorv1.RedisOperator{
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
					VolumePath: redisCrd.Master.VolumePath,
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
					VolumePath: redisCrd.Slave.VolumePath,
				},
			},
		},
	}
}

func convertOperatorToRedisCrd(v *redisoperatorv1.RedisOperator) proto.RedisCrd {
	return proto.RedisCrd{
		Name: v.Name,
		Master: proto.NodeSpec{
			Name:             v.Spec.MasterSpec.Spec.Name,
			Replicas:         *v.Spec.MasterSpec.Spec.Replicas,
			Image:            v.Spec.MasterSpec.Spec.Image,
			ImagePullSecrets: v.Spec.MasterSpec.Spec.ImagePullSecrets[0].Name,
			VolumePath:       v.Spec.MasterSpec.Spec.VolumePath,
		},
		Slave: proto.NodeSpec{
			Name:             v.Spec.SlaveSpec.Spec.Name,
			Replicas:         *v.Spec.SlaveSpec.Spec.Replicas,
			Image:            v.Spec.SlaveSpec.Spec.Image,
			ImagePullSecrets: v.Spec.SlaveSpec.Spec.ImagePullSecrets[0].Name,
			VolumePath:       v.Spec.SlaveSpec.Spec.VolumePath,
		},
	}
}
