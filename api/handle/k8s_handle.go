package handle

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/klog"

	"github.com/nevercase/k8s-controller-custom-resource/api/group"
	"github.com/nevercase/k8s-controller-custom-resource/api/proto"
	mysqloperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/mysqloperator/v1"
	redisoperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/redisoperator/v1"
)

type KubernetesApiGetter interface {
	KubernetesApi() KubernetesApiInterface
}

type KubernetesApiInterface interface {
	Create(req proto.Param, obj []byte) (res []byte, err error)
	Delete(req proto.Param, obj []byte) (err error)
	Get(req proto.Param, obj []byte) (res []byte, err error)
	List(req proto.Param) ([]byte, error)
	Resources(req proto.Param) (res []byte, err error)
}

func NewKubernetesApiHandle(g group.Group) KubernetesApiInterface {
	return &k8sHandle{
		group: g,
	}
}

type k8sHandle struct {
	group group.Group
}

func (h *k8sHandle) Create(req proto.Param, obj []byte) (res []byte, err error) {
	var n interface{}
	switch req.ResourceType {
	case group.ConfigMap:
		var cm proto.ConfigMap
		if err = cm.Unmarshal(obj); err != nil {
			break
		}
		m := convertProtoToConfigMap(req, cm)
		if n, err = resourceCreateOrUpdate(h.group, req, m.Name, m); err != nil {
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
		if n, err = resourceCreateOrUpdate(h.group, req, m.Name, m); err != nil {
			break
		}
		v := n.(*corev1.Namespace)
		e := convertNameSpaceToProto(v)
		res, err = e.Marshal()
	case group.Service:
		var s proto.Service
		if err = s.Unmarshal(obj); err != nil {
			break
		}
		m := convertProtoToService(req, s)
		if n, err = resourceCreateOrUpdate(h.group, req, m.Name, m); err != nil {
			break
		}
		v := n.(*corev1.Service)
		e := convertServiceToProto(v)
		res, err = e.Marshal()
	case group.MysqlOperator:
		var mysqlCrd proto.MysqlCrd
		if err = mysqlCrd.Unmarshal(obj); err != nil {
			break
		}
		m := convertMysqlCrdToProto(req, mysqlCrd)
		if n, err = resourceCreateOrUpdate(h.group, req, m.Name, m); err != nil {
			break
		}
		v := n.(*mysqloperatorv1.MysqlOperator)
		e := convertProtoToMysqlCrd(v)
		res, err = e.Marshal()
	case group.RedisOperator:
		var redisCrd proto.RedisCrd
		if err = redisCrd.Unmarshal(obj); err != nil {
			break
		}
		m := convertProtoToRedisCrd(req, redisCrd)
		if n, err = resourceCreateOrUpdate(h.group, req, m.Name, m); err != nil {
			break
		}
		v := n.(*redisoperatorv1.RedisOperator)
		e := convertRedisCrdToProto(v)
		res, err = e.Marshal()
	case group.HelixOperator:
	}
	return proto.GetResponse(req, res)
}

func (h *k8sHandle) Delete(req proto.Param, obj []byte) (err error) {
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
	case group.Service:
		var s proto.Service
		if err = s.Unmarshal(obj); err != nil {
			break
		}
		name = s.Name
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

func (h *k8sHandle) Get(req proto.Param, obj []byte) (res []byte, err error) {
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
	case group.Service:
		var s proto.Service
		if err = s.Unmarshal(obj); err != nil {
			break
		}
		n, err = h.group.Resource().Get(req.ResourceType, req.NameSpace, s.Name)
		if err != nil {
			break
		}
		m := n.(*corev1.Service)
		e := convertServiceToProto(m)
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
		e := convertProtoToMysqlCrd(m)
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
		e := convertRedisCrdToProto(m)
		res, err = e.Marshal()
	case group.HelixOperator:
	}
	return proto.GetResponse(req, res)
}

func (h *k8sHandle) List(req proto.Param) (res []byte, err error) {
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
	case group.Service:
		m := proto.ServiceList{
			Items: make([]proto.Service, 0),
		}
		for _, v := range d.(*corev1.ServiceList).Items {
			m.Items = append(m.Items, convertServiceToProto(&v))
		}
		res, err = m.Marshal()
	case group.MysqlOperator:
		m := proto.MysqlCrdList{
			Items: make([]proto.MysqlCrd, 0),
		}
		for _, v := range d.(*mysqloperatorv1.MysqlOperatorList).Items {
			m.Items = append(m.Items, convertProtoToMysqlCrd(&v))
		}
		res, err = m.Marshal()
	case group.RedisOperator:
		m := proto.RedisCrdList{
			Items: make([]proto.RedisCrd, 0),
		}
		for _, v := range d.(*redisoperatorv1.RedisOperatorList).Items {
			m.Items = append(m.Items, convertRedisCrdToProto(&v))
		}
		res, err = m.Marshal()
	case group.HelixOperator:
	}
	if err != nil {
		return nil, err
	}
	return proto.GetResponse(req, res)
}

func (h *k8sHandle) Resources(req proto.Param) (res []byte, err error) {
	m := proto.ResourceList{
		Items: h.group.Resource().ResourceTypes(),
	}
	o, err := m.Marshal()
	if err != nil {
		return nil, err
	}
	return proto.GetResponse(req, o)
}

func resourceCreateOrUpdate(g group.Group, req proto.Param, specName string, m interface{}) (res interface{}, err error) {
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

func convertPodResourceLimitsToProto(res proto.PodResourceList) corev1.ResourceList {
	rl := make(map[corev1.ResourceName]resource.Quantity, 0)
	for k, v := range res {
		t, err := resource.ParseQuantity(v)
		if err != nil {
			klog.V(2).Info(err)
			continue
		}
		rl[corev1.ResourceName(k)] = t
	}
	return rl
}

func convertProtoToPodResourceLimits(rl corev1.ResourceList) proto.PodResourceList {
	res := make(map[string]string, 0)
	for k, v := range rl {
		res[string(k)] = v.String()
	}
	return res
}

func convertResourceRequirementsToProto(res proto.PodResourceRequirements) corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Limits:   convertPodResourceLimitsToProto(res.Limits),
		Requests: convertPodResourceLimitsToProto(res.Requests),
	}
}

func convertProtoToResourceRequirements(rl corev1.ResourceRequirements) proto.PodResourceRequirements {
	return proto.PodResourceRequirements{
		Limits:   convertProtoToPodResourceLimits(rl.Limits),
		Requests: convertProtoToPodResourceLimits(rl.Requests),
	}
}

func convertMysqlCrdToProto(req proto.Param, mysqlCrd proto.MysqlCrd) *mysqloperatorv1.MysqlOperator {
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
					Resources:  convertResourceRequirementsToProto(mysqlCrd.Master.PodResource),
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
					Resources:  convertResourceRequirementsToProto(mysqlCrd.Slave.PodResource),
				},
			},
		},
	}
}

func convertProtoToMysqlCrd(m *mysqloperatorv1.MysqlOperator) proto.MysqlCrd {
	return proto.MysqlCrd{
		Name: m.Name,
		Master: proto.NodeSpec{
			Name:             m.Spec.MasterSpec.Spec.Name,
			Replicas:         *m.Spec.MasterSpec.Spec.Replicas,
			Image:            m.Spec.MasterSpec.Spec.Image,
			ImagePullSecrets: m.Spec.MasterSpec.Spec.ImagePullSecrets[0].Name,
			VolumePath:       m.Spec.MasterSpec.Spec.VolumePath,
			PodResource:      convertProtoToResourceRequirements(m.Spec.MasterSpec.Spec.Resources),
		},
		Slave: proto.NodeSpec{
			Name:             m.Spec.SlaveSpec.Spec.Name,
			Replicas:         *m.Spec.SlaveSpec.Spec.Replicas,
			Image:            m.Spec.SlaveSpec.Spec.Image,
			ImagePullSecrets: m.Spec.SlaveSpec.Spec.ImagePullSecrets[0].Name,
			VolumePath:       m.Spec.SlaveSpec.Spec.VolumePath,
			PodResource:      convertProtoToResourceRequirements(m.Spec.SlaveSpec.Spec.Resources),
		},
	}
}

func convertProtoToRedisCrd(req proto.Param, redisCrd proto.RedisCrd) *redisoperatorv1.RedisOperator {
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
					Resources:  convertResourceRequirementsToProto(redisCrd.Master.PodResource),
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
					Resources:  convertResourceRequirementsToProto(redisCrd.Slave.PodResource),
				},
			},
		},
	}
}

func convertRedisCrdToProto(v *redisoperatorv1.RedisOperator) proto.RedisCrd {
	return proto.RedisCrd{
		Name: v.Name,
		Master: proto.NodeSpec{
			Name:             v.Spec.MasterSpec.Spec.Name,
			Replicas:         *v.Spec.MasterSpec.Spec.Replicas,
			Image:            v.Spec.MasterSpec.Spec.Image,
			ImagePullSecrets: v.Spec.MasterSpec.Spec.ImagePullSecrets[0].Name,
			VolumePath:       v.Spec.MasterSpec.Spec.VolumePath,
			PodResource:      convertProtoToResourceRequirements(v.Spec.MasterSpec.Spec.Resources),
		},
		Slave: proto.NodeSpec{
			Name:             v.Spec.SlaveSpec.Spec.Name,
			Replicas:         *v.Spec.SlaveSpec.Spec.Replicas,
			Image:            v.Spec.SlaveSpec.Spec.Image,
			ImagePullSecrets: v.Spec.SlaveSpec.Spec.ImagePullSecrets[0].Name,
			VolumePath:       v.Spec.SlaveSpec.Spec.VolumePath,
			PodResource:      convertProtoToResourceRequirements(v.Spec.SlaveSpec.Spec.Resources),
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

func convertProtoToService(req proto.Param, v proto.Service) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.Name,
			Namespace: req.NameSpace,
		},
		Spec: corev1.ServiceSpec{
			Ports:       convertProtoToServicePort(v.Ports),
			ClusterIP:   v.ClusterIP,
			Type:        corev1.ServiceType(v.Type),
			ExternalIPs: v.ExternalIPs,
		},
	}
}

func convertServiceToProto(s *corev1.Service) proto.Service {
	return proto.Service{
		Name:        s.Name,
		Ports:       convertServicePortToProto(s.Spec.Ports),
		ClusterIP:   s.Spec.ClusterIP,
		Type:        string(s.Spec.Type),
		ExternalIPs: s.Spec.ExternalIPs,
	}
}

func convertProtoToServicePort(p []proto.ServicePort) []corev1.ServicePort {
	res := make([]corev1.ServicePort, 0)
	for _, v := range p {
		res = append(res, corev1.ServicePort{
			Name:     v.Name,
			Protocol: corev1.Protocol(v.Protocol),
			Port:     v.Port,
			TargetPort: intstr.IntOrString{
				Type:   intstr.Type(v.TargetPort.Type),
				IntVal: v.TargetPort.IntVal,
				StrVal: v.TargetPort.StrVal,
			},
			NodePort: v.NodePort,
		})
	}
	return res
}

func convertServicePortToProto(p []corev1.ServicePort) []proto.ServicePort {
	res := make([]proto.ServicePort, 0)
	for _, v := range p {
		res = append(res, proto.ServicePort{
			Name:     v.Name,
			Protocol: string(v.Protocol),
			Port:     v.Port,
			TargetPort: proto.IntOrString{
				Type:   int32(v.TargetPort.Type),
				IntVal: v.TargetPort.IntVal,
				StrVal: v.TargetPort.StrVal,
			},
			NodePort: v.NodePort,
		})
	}
	return res
}
