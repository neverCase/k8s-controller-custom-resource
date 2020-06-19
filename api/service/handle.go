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
	Get(req proto.Param, obj []byte) (res []byte, err error)
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
	var n interface{}
	switch req.ResourceType {
	case group.ConfigMap:
		var cm proto.ConfigMap
		if err = cm.Unmarshal(obj); err != nil {
			break
		}
		m := convertProtoToConfigMap(req, cm)
		if n, err = resourceCreateWithUpdate(h.group, req, m.Name, m); err != nil {
			break
		}
		v := n.(*corev1.ConfigMap)
		e := convertConfigMapToProto(v)
		res, err = e.Marshal()
	case group.NameSpace:
		var ns proto.NameSpace
		if err = ns.Unmarshal(obj); err != nil {
			break
		}
		m := convertProtoToNameSpace(ns)
		if n, err = resourceCreateWithUpdate(h.group, req, m.Name, m); err != nil {
			break
		}
		v := n.(*corev1.Namespace)
		e := convertNameSpaceToProto(v)
		res, err = e.Marshal()
	case group.MysqlOperator:
		var mysqlCrd proto.MysqlCrd
		if err = mysqlCrd.Unmarshal(obj); err != nil {
			break
		}
		m := convertMysqlCrdToOperator(req, mysqlCrd)
		if n, err = resourceCreateWithUpdate(h.group, req, m.Name, m); err != nil {
			break
		}
		v := n.(*mysqloperatorv1.MysqlOperator)
		e := convertOperatorToMysqlCrd(v)
		res, err = e.Marshal()
	case group.RedisOperator:
		var redisCrd proto.RedisCrd
		if err = redisCrd.Unmarshal(obj); err != nil {
			break
		}
		m := convertRedisCrdToOperator(req, redisCrd)
		if n, err = resourceCreateWithUpdate(h.group, req, m.Name, m); err != nil {
			break
		}
		v := n.(*redisoperatorv1.RedisOperator)
		e := convertOperatorToRedisCrd(v)
		res, err = e.Marshal()
	case group.HelixOperator:
	}
	return proto.GetResponse(req, res)
}

func (h *handle) Delete(req proto.Param, obj []byte) (err error) {
	var name string
	switch req.ResourceType {
	case group.ConfigMap:
		var cm proto.ConfigMap
		if err = cm.Unmarshal(obj); err != nil {
			break
		}
		name = cm.Name
	case group.NameSpace:
		var ns proto.NameSpace
		if err = ns.Unmarshal(obj); err != nil {
			break
		}
		name = ns.Name
	case group.MysqlOperator:
		var mysqlCrd proto.MysqlCrd
		if err = mysqlCrd.Unmarshal(obj); err != nil {
			break
		}
		name = mysqlCrd.Name
	case group.RedisOperator:
		var redisCrd proto.RedisCrd
		if err = redisCrd.Unmarshal(obj); err != nil {
			break
		}
		name = redisCrd.Name
	case group.HelixOperator:
	}
	err = h.group.Resource().Delete(req.ResourceType, req.NameSpace, name)
	return err
}

func (h *handle) Get(req proto.Param, obj []byte) (res []byte, err error) {
	var n interface{}
	switch req.ResourceType {
	case group.ConfigMap:
		var cm proto.ConfigMap
		if err = cm.Unmarshal(obj); err != nil {
			break
		}
		n, err = h.group.Resource().Get(req.ResourceType, req.NameSpace, cm.Name)
		if err != nil {
			break
		}
		m := n.(*corev1.ConfigMap)
		e := convertConfigMapToProto(m)
		res, err = e.Marshal()
	case group.NameSpace:
		var ns proto.NameSpace
		if err = ns.Unmarshal(obj); err != nil {
			break
		}
		n, err = h.group.Resource().Get(req.ResourceType, req.NameSpace, ns.Name)
		if err != nil {
			break
		}
		m := n.(*corev1.Namespace)
		e := convertNameSpaceToProto(m)
		res, err = e.Marshal()
	case group.MysqlOperator:
		var mysqlCrd proto.MysqlCrd
		if err = mysqlCrd.Unmarshal(obj); err != nil {
			break
		}
		n, err = h.group.Resource().Get(req.ResourceType, req.NameSpace, mysqlCrd.Name)
		if err != nil {
			break
		}
		m := n.(*mysqloperatorv1.MysqlOperator)
		e := convertOperatorToMysqlCrd(m)
		res, err = e.Marshal()
	case group.RedisOperator:
		var redisCrd proto.RedisCrd
		if err = redisCrd.Unmarshal(obj); err != nil {
			break
		}
		n, err = h.group.Resource().Get(req.ResourceType, req.NameSpace, redisCrd.Name)
		if err != nil {
			break
		}
		m := n.(*redisoperatorv1.RedisOperator)
		e := convertOperatorToRedisCrd(m)
		res, err = e.Marshal()
	case group.HelixOperator:
	}
	return proto.GetResponse(req, res)
}

func (h *handle) List(req proto.Param) (res []byte, err error) {
	var d interface{}
	var selector = labels.NewSelector()
	if d, err = h.group.Resource().List(req.ResourceType, req.NameSpace, selector); err != nil {
		return res, err
	}
	switch req.ResourceType {
	case group.ConfigMap:
		m := proto.ConfigMapList{
			Items: make([]proto.ConfigMap, 0),
		}
		for _, v := range d.(*corev1.ConfigMapList).Items {
			m.Items = append(m.Items, convertConfigMapToProto(&v))
		}
		res, err = m.Marshal()
	case group.NameSpace:
		m := proto.NameSpaceList{
			Items: make([]proto.NameSpace, 0),
		}
		for _, v := range d.(*corev1.NamespaceList).Items {
			m.Items = append(m.Items, convertNameSpaceToProto(&v))
		}
		res, err = m.Marshal()
	case group.MysqlOperator:
		m := proto.MysqlCrdList{
			Items: make([]proto.MysqlCrd, 0),
		}
		for _, v := range d.(*mysqloperatorv1.MysqlOperatorList).Items {
			m.Items = append(m.Items, convertOperatorToMysqlCrd(&v))
		}
		res, err = m.Marshal()
	case group.RedisOperator:
		m := proto.RedisCrdList{
			Items: make([]proto.RedisCrd, 0),
		}
		for _, v := range d.(*redisoperatorv1.RedisOperatorList).Items {
			m.Items = append(m.Items, convertOperatorToRedisCrd(&v))
		}
		res, err = m.Marshal()
	case group.HelixOperator:
	}
	if err != nil {
		return nil, err
	}
	return proto.GetResponse(req, res)
}

func (h *handle) Resources(req proto.Param) (res []byte, err error) {
	m := proto.ResourceList{
		Items: h.group.Resource().ResourceTypes(),
	}
	o, err := m.Marshal()
	if err != nil {
		return nil, err
	}
	return proto.GetResponse(req, o)
}

func resourceCreateWithUpdate(g group.Group, req proto.Param, specName string, m interface{}) (res interface{}, err error) {
	_, err = g.Resource().Get(req.ResourceType, req.NameSpace, specName)
	if err != nil {
		if !errors.IsNotFound(err) {
			return
		}
		err = nil
		if res, err = g.Resource().Create(req.ResourceType, req.NameSpace, m); err != nil {
			return
		}
	} else {
		if res, err = g.Resource().Update(req.ResourceType, req.NameSpace, m); err != nil {
			return
		}
	}
	return
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

func convertProtoToConfigMap(req proto.Param, v proto.ConfigMap) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.Name,
			Namespace: req.NameSpace,
		},
		Data: v.Data,
	}
}

func convertConfigMapToProto(c *corev1.ConfigMap) proto.ConfigMap {
	return proto.ConfigMap{
		Name: c.Name,
		Data: c.Data,
	}
}

func convertProtoToNameSpace(v proto.NameSpace) *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: v.Name,
		},
	}
}

func convertNameSpaceToProto(c *corev1.Namespace) proto.NameSpace {
	return proto.NameSpace{
		Name: c.Name,
	}
}
