package handle

import (
	"reflect"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/klog"

	helixsagaoperatorv1 "github.com/Shanghai-Lunara/helixsaga-operator/pkg/apis/helixsaga/v1"
	mysqloperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/mysqloperator/v1"
	redisoperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/redisoperator/v1"

	"github.com/nevercase/k8s-controller-custom-resource/api/group"
	"github.com/nevercase/k8s-controller-custom-resource/api/proto"
)

type KubernetesApiGetter interface {
	KubernetesApi() KubernetesApiInterface
}

type KubernetesApiInterface interface {
	Create(req proto.Param, obj []byte) (res []byte, err error)
	Update(req proto.Param, obj []byte) (res []byte, err error)
	Delete(req proto.Param, obj []byte) (err error)
	Get(req proto.Param, obj []byte) (res []byte, err error)
	List(req proto.Param) ([]byte, error)
	Watch()
	Resources(req proto.Param) (res []byte, err error)
}

func NewKubernetesApiHandle(g group.Group) KubernetesApiInterface {
	kh := &k8sHandle{
		group: g,
	}
	go kh.Watch()
	return kh
}

type k8sHandle struct {
	group group.Group
}

func (h *k8sHandle) Create(req proto.Param, obj []byte) (res []byte, err error) {
	var n interface{}
	switch req.ResourceType {
	case group.ConfigMap:
		var e proto.ConfigMap
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		m := convertProtoToConfigMap(req, e)
		if n, err = resourceCreate(h.group, req, m.Name, m); err != nil {
			break
		}
		e = convertConfigMapToProto(n.(*corev1.ConfigMap))
		res, err = e.Marshal()
	case group.NameSpace:
		var e proto.NameSpace
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		m := convertProtoToNameSpace(e)
		if n, err = resourceCreate(h.group, req, m.Name, m); err != nil {
			break
		}
		e = convertNameSpaceToProto(n.(*corev1.Namespace))
		res, err = e.Marshal()
	case group.Service:
		var e proto.Service
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		m := convertProtoToService(req, e)
		if n, err = resourceCreate(h.group, req, m.Name, m); err != nil {
			break
		}
		e = convertServiceToProto(n.(*corev1.Service))
		res, err = e.Marshal()
	case group.MysqlOperator:
		var e proto.MysqlCrd
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		m := convertMysqlCrdToProto(req, e)
		if n, err = resourceCreate(h.group, req, m.Name, m); err != nil {
			break
		}
		e = convertProtoToMysqlCrd(n.(*mysqloperatorv1.MysqlOperator))
		res, err = e.Marshal()
	case group.RedisOperator:
		var e proto.RedisCrd
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		m := convertProtoToRedisCrd(req, e)
		if n, err = resourceCreate(h.group, req, m.Name, m); err != nil {
			break
		}
		e = convertRedisCrdToProto(n.(*redisoperatorv1.RedisOperator))
		res, err = e.Marshal()
	case group.HelixSagaOperator:
		var e proto.HelixSagaCrd
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		m := convertProtoToHelixSagaCrd(req, e)
		if n, err = resourceCreate(h.group, req, m.Name, m); err != nil {
			break
		}
		e = covertHelixSagaCrdToProto(n.(*helixsagaoperatorv1.HelixSaga))
		res, err = e.Marshal()
	}
	if err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	return proto.GetResponse(req, res)
}

func (h *k8sHandle) Update(req proto.Param, obj []byte) (res []byte, err error) {
	var n interface{}
	switch req.ResourceType {
	case group.ConfigMap:
		var e proto.ConfigMap
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		m := convertProtoToConfigMap(req, e)
		if n, err = resourceUpdate(h.group, req, m.Name, m); err != nil {
			break
		}
		e = convertConfigMapToProto(n.(*corev1.ConfigMap))
		res, err = e.Marshal()
	case group.NameSpace:
		var e proto.NameSpace
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		m := convertProtoToNameSpace(e)
		if n, err = resourceUpdate(h.group, req, m.Name, m); err != nil {
			break
		}
		e = convertNameSpaceToProto(n.(*corev1.Namespace))
		res, err = e.Marshal()
	case group.Service:
		var e proto.Service
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		m := convertProtoToService(req, e)
		if n, err = resourceUpdate(h.group, req, m.Name, m); err != nil {
			break
		}
		e = convertServiceToProto(n.(*corev1.Service))
		res, err = e.Marshal()
	case group.MysqlOperator:
		var e proto.MysqlCrd
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		m := convertMysqlCrdToProto(req, e)
		if n, err = resourceUpdate(h.group, req, m.Name, m); err != nil {
			break
		}
		e = convertProtoToMysqlCrd(n.(*mysqloperatorv1.MysqlOperator))
		res, err = e.Marshal()
	case group.RedisOperator:
		var e proto.RedisCrd
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		m := convertProtoToRedisCrd(req, e)
		if n, err = resourceUpdate(h.group, req, m.Name, m); err != nil {
			break
		}
		e = convertRedisCrdToProto(n.(*redisoperatorv1.RedisOperator))
		res, err = e.Marshal()
	case group.HelixSagaOperator:
		var e proto.HelixSagaCrd
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		m := convertProtoToHelixSagaCrd(req, e)
		if n, err = resourceUpdate(h.group, req, m.Name, m); err != nil {
			break
		}
		e = covertHelixSagaCrdToProto(n.(*helixsagaoperatorv1.HelixSaga))
		res, err = e.Marshal()
	}
	if err != nil {
		klog.V(2).Info(err)
		return nil, err
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
		var e proto.MysqlCrd
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		name = e.Name
	case group.RedisOperator:
		var e proto.RedisCrd
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		name = e.Name
	case group.HelixSagaOperator:
		var e proto.HelixSagaCrd
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		name = e.Name
	}
	if err != nil {
		klog.V(2).Info(err)
		return err
	}
	err = h.group.Resource().Delete(req.ResourceType, req.NameSpace, name)
	return err
}

func (h *k8sHandle) Get(req proto.Param, obj []byte) (res []byte, err error) {
	var n interface{}
	switch req.ResourceType {
	case group.ConfigMap:
		var e proto.ConfigMap
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		n, err = h.group.Resource().Get(req.ResourceType, req.NameSpace, e.Name)
		if err != nil {
			break
		}
		e = convertConfigMapToProto(n.(*corev1.ConfigMap))
		res, err = e.Marshal()
	case group.NameSpace:
		var e proto.NameSpace
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		n, err = h.group.Resource().Get(req.ResourceType, req.NameSpace, e.Name)
		if err != nil {
			break
		}
		e = convertNameSpaceToProto(n.(*corev1.Namespace))
		res, err = e.Marshal()
	case group.Service:
		var e proto.Service
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		n, err = h.group.Resource().Get(req.ResourceType, req.NameSpace, e.Name)
		if err != nil {
			break
		}
		e = convertServiceToProto(n.(*corev1.Service))
		res, err = e.Marshal()
	case group.Secret:
		var e proto.Secret
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		n, err = h.group.Resource().Get(req.ResourceType, req.NameSpace, e.Name)
		if err != nil {
			break
		}
		e = convertSecretToProto(n.(*corev1.Secret))
		res, err = e.Marshal()
	case group.MysqlOperator:
		var e proto.MysqlCrd
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		n, err = h.group.Resource().Get(req.ResourceType, req.NameSpace, e.Name)
		if err != nil {
			break
		}
		e = convertProtoToMysqlCrd(n.(*mysqloperatorv1.MysqlOperator))
		res, err = e.Marshal()
	case group.RedisOperator:
		var e proto.RedisCrd
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		n, err = h.group.Resource().Get(req.ResourceType, req.NameSpace, e.Name)
		if err != nil {
			break
		}
		e = convertRedisCrdToProto(n.(*redisoperatorv1.RedisOperator))
		res, err = e.Marshal()
	case group.HelixSagaOperator:
		var e proto.HelixSagaCrd
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		n, err = h.group.Resource().Get(req.ResourceType, req.NameSpace, e.Name)
		if err != nil {
			break
		}
		e = covertHelixSagaCrdToProto(n.(*helixsagaoperatorv1.HelixSaga))
		res, err = e.Marshal()
	}
	if err != nil {
		klog.V(2).Info(err)
		return nil, err
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
	case group.Secret:
		m := proto.SecretList{
			Items: make([]proto.Secret, 0),
		}
		for _, v := range d.(*corev1.SecretList).Items {
			m.Items = append(m.Items, convertSecretToProto(&v))
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
	case group.HelixSagaOperator:
		m := proto.HelixSagaCrdList{
			Items: make([]proto.HelixSagaCrd, 0),
		}
		for _, v := range d.(*helixsagaoperatorv1.HelixSagaList).Items {
			m.Items = append(m.Items, covertHelixSagaCrdToProto(&v))
		}
		res, err = m.Marshal()
	}
	if err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	return proto.GetResponse(req, res)
}

func (h *k8sHandle) convertObjFromEvent(obj interface{}) (res []byte, err error) {
	req := proto.Param{
		Service: string(proto.SvcWatch),
	}
	switch reflect.TypeOf(obj) {
	case reflect.TypeOf(&corev1.ConfigMap{}):
		var e proto.ConfigMap
		n := obj.(*corev1.ConfigMap)
		req.ResourceType = group.ConfigMap
		req.NameSpace = n.Namespace
		e = convertConfigMapToProto(n)
		res, err = e.Marshal()
	case reflect.TypeOf(&corev1.Namespace{}):
		var e proto.NameSpace
		n := obj.(*corev1.Namespace)
		req.ResourceType = group.NameSpace
		e = convertNameSpaceToProto(n)
		res, err = e.Marshal()
	case reflect.TypeOf(&corev1.Service{}):
		var e proto.Service
		n := obj.(*corev1.Service)
		req.ResourceType = group.Service
		req.NameSpace = n.Namespace
		e = convertServiceToProto(n)
		res, err = e.Marshal()
	case reflect.TypeOf(&corev1.Secret{}):
		var e proto.Secret
		n := obj.(*corev1.Secret)
		req.ResourceType = group.Secret
		req.NameSpace = n.Namespace
		e = convertSecretToProto(n)
		res, err = e.Marshal()
	case reflect.TypeOf(&mysqloperatorv1.MysqlOperator{}):
		var e proto.MysqlCrd
		n := obj.(*mysqloperatorv1.MysqlOperator)
		req.ResourceType = group.MysqlOperator
		req.NameSpace = n.Namespace
		e = convertProtoToMysqlCrd(n)
		res, err = e.Marshal()
	case reflect.TypeOf(&redisoperatorv1.RedisOperator{}):
		var e proto.RedisCrd
		n := obj.(*redisoperatorv1.RedisOperator)
		req.ResourceType = group.RedisOperator
		req.NameSpace = n.Namespace
		e = convertRedisCrdToProto(n)
		res, err = e.Marshal()
	case reflect.TypeOf(&helixsagaoperatorv1.HelixSaga{}):
		var e proto.HelixSagaCrd
		n := obj.(*helixsagaoperatorv1.HelixSaga)
		req.ResourceType = group.HelixSagaOperator
		req.NameSpace = n.Namespace
		e = covertHelixSagaCrdToProto(n)
		res, err = e.Marshal()
	}
	if err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	return proto.GetResponse(req, res)
}

func (h *k8sHandle) Watch() {
	for {
		select {
		case e, isClosed := <-h.group.WatchEvents():
			klog.Info("watch h.group.WatchEvents:", e)
			if !isClosed {
				return
			}
			res, err := h.convertObjFromEvent(e.Object)
			if err != nil {
				klog.V(2).Info(err)
				continue
			}
			klog.Info("watch obj:", res)
		}
	}
}

func (h *k8sHandle) Resources(req proto.Param) (res []byte, err error) {
	m := proto.ResourceList{
		Items: h.group.Resource().ResourceTypes(),
	}
	o, err := m.Marshal()
	if err != nil {
		klog.V(2).Info(err)
		return nil, err
	}
	return proto.GetResponse(req, o)
}

func resourceCreate(g group.Group, req proto.Param, specName string, m interface{}) (res interface{}, err error) {
	_, err = g.Resource().Get(req.ResourceType, req.NameSpace, specName)
	if err != nil {
		if !errors.IsNotFound(err) {
			return res, err
		}
	}
	return g.Resource().Create(req.ResourceType, req.NameSpace, m)
}

func resourceUpdate(g group.Group, req proto.Param, specName string, m interface{}) (res interface{}, err error) {
	_, err = g.Resource().Get(req.ResourceType, req.NameSpace, specName)
	if err != nil {
		return nil, err
	}
	return g.Resource().Update(req.ResourceType, req.NameSpace, m)
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
			Name:            mysqlCrd.Name,
			Namespace:       req.NameSpace,
			ResourceVersion: mysqlCrd.ResourceVersion,
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
					VolumePath:     mysqlCrd.Master.VolumePath,
					Resources:      convertResourceRequirementsToProto(mysqlCrd.Master.PodResource),
					ContainerPorts: convertProtoToContainerPort(mysqlCrd.Master.ContainerPorts),
					ServicePorts:   convertProtoToServicePort(mysqlCrd.Master.ServicePorts),
					Env:            convertProtoToEnvVar(mysqlCrd.Master.Env),
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
					VolumePath:     mysqlCrd.Slave.VolumePath,
					Resources:      convertResourceRequirementsToProto(mysqlCrd.Slave.PodResource),
					ContainerPorts: convertProtoToContainerPort(mysqlCrd.Slave.ContainerPorts),
					ServicePorts:   convertProtoToServicePort(mysqlCrd.Slave.ServicePorts),
					Env:            convertProtoToEnvVar(mysqlCrd.Slave.Env),
				},
			},
		},
	}
}

func convertProtoToMysqlCrd(m *mysqloperatorv1.MysqlOperator) proto.MysqlCrd {
	return proto.MysqlCrd{
		Name:            m.Name,
		ResourceVersion: m.ResourceVersion,
		Master: proto.NodeSpec{
			Name:             m.Spec.MasterSpec.Spec.Name,
			Replicas:         *m.Spec.MasterSpec.Spec.Replicas,
			Image:            m.Spec.MasterSpec.Spec.Image,
			ImagePullSecrets: m.Spec.MasterSpec.Spec.ImagePullSecrets[0].Name,
			VolumePath:       m.Spec.MasterSpec.Spec.VolumePath,
			PodResource:      convertProtoToResourceRequirements(m.Spec.MasterSpec.Spec.Resources),
			ContainerPorts:   convertContainerPortToProto(m.Spec.MasterSpec.Spec.ContainerPorts),
			ServicePorts:     convertServicePortToProto(m.Spec.MasterSpec.Spec.ServicePorts),
			Env:              convertEnvVarToProto(m.Spec.MasterSpec.Spec.Env),
		},
		Slave: proto.NodeSpec{
			Name:             m.Spec.SlaveSpec.Spec.Name,
			Replicas:         *m.Spec.SlaveSpec.Spec.Replicas,
			Image:            m.Spec.SlaveSpec.Spec.Image,
			ImagePullSecrets: m.Spec.SlaveSpec.Spec.ImagePullSecrets[0].Name,
			VolumePath:       m.Spec.SlaveSpec.Spec.VolumePath,
			PodResource:      convertProtoToResourceRequirements(m.Spec.SlaveSpec.Spec.Resources),
			ContainerPorts:   convertContainerPortToProto(m.Spec.SlaveSpec.Spec.ContainerPorts),
			ServicePorts:     convertServicePortToProto(m.Spec.SlaveSpec.Spec.ServicePorts),
			Env:              convertEnvVarToProto(m.Spec.SlaveSpec.Spec.Env),
		},
	}
}

func convertProtoToRedisCrd(req proto.Param, redisCrd proto.RedisCrd) *redisoperatorv1.RedisOperator {
	return &redisoperatorv1.RedisOperator{
		ObjectMeta: metav1.ObjectMeta{
			Name:            redisCrd.Name,
			Namespace:       req.NameSpace,
			ResourceVersion: redisCrd.ResourceVersion,
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
					VolumePath:     redisCrd.Master.VolumePath,
					Resources:      convertResourceRequirementsToProto(redisCrd.Master.PodResource),
					ContainerPorts: convertProtoToContainerPort(redisCrd.Master.ContainerPorts),
					ServicePorts:   convertProtoToServicePort(redisCrd.Master.ServicePorts),
					Env:            convertProtoToEnvVar(redisCrd.Master.Env),
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
					VolumePath:     redisCrd.Slave.VolumePath,
					Resources:      convertResourceRequirementsToProto(redisCrd.Slave.PodResource),
					ContainerPorts: convertProtoToContainerPort(redisCrd.Slave.ContainerPorts),
					ServicePorts:   convertProtoToServicePort(redisCrd.Slave.ServicePorts),
					Env:            convertProtoToEnvVar(redisCrd.Slave.Env),
				},
			},
		},
	}
}

func convertRedisCrdToProto(v *redisoperatorv1.RedisOperator) proto.RedisCrd {
	return proto.RedisCrd{
		Name:            v.Name,
		ResourceVersion: v.ResourceVersion,
		Master: proto.NodeSpec{
			Name:             v.Spec.MasterSpec.Spec.Name,
			Replicas:         *v.Spec.MasterSpec.Spec.Replicas,
			Image:            v.Spec.MasterSpec.Spec.Image,
			ImagePullSecrets: v.Spec.MasterSpec.Spec.ImagePullSecrets[0].Name,
			VolumePath:       v.Spec.MasterSpec.Spec.VolumePath,
			PodResource:      convertProtoToResourceRequirements(v.Spec.MasterSpec.Spec.Resources),
			ContainerPorts:   convertContainerPortToProto(v.Spec.MasterSpec.Spec.ContainerPorts),
			ServicePorts:     convertServicePortToProto(v.Spec.MasterSpec.Spec.ServicePorts),
			Env:              convertEnvVarToProto(v.Spec.MasterSpec.Spec.Env),
		},
		Slave: proto.NodeSpec{
			Name:             v.Spec.SlaveSpec.Spec.Name,
			Replicas:         *v.Spec.SlaveSpec.Spec.Replicas,
			Image:            v.Spec.SlaveSpec.Spec.Image,
			ImagePullSecrets: v.Spec.SlaveSpec.Spec.ImagePullSecrets[0].Name,
			VolumePath:       v.Spec.SlaveSpec.Spec.VolumePath,
			PodResource:      convertProtoToResourceRequirements(v.Spec.SlaveSpec.Spec.Resources),
			ContainerPorts:   convertContainerPortToProto(v.Spec.SlaveSpec.Spec.ContainerPorts),
			ServicePorts:     convertServicePortToProto(v.Spec.SlaveSpec.Spec.ServicePorts),
			Env:              convertEnvVarToProto(v.Spec.SlaveSpec.Spec.Env),
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

func convertProtoToContainerPort(c []proto.ContainerPort) []corev1.ContainerPort {
	res := make([]corev1.ContainerPort, 0)
	for _, v := range c {
		res = append(res, corev1.ContainerPort{
			Name:          v.Name,
			HostPort:      v.HostPort,
			ContainerPort: v.ContainerPort,
			Protocol:      corev1.Protocol(v.Protocol),
			HostIP:        v.HostIP,
		})
	}
	return res
}

func convertContainerPortToProto(c []corev1.ContainerPort) []proto.ContainerPort {
	res := make([]proto.ContainerPort, 0)
	for _, v := range c {
		res = append(res, proto.ContainerPort{
			Name:          v.Name,
			HostPort:      v.HostPort,
			ContainerPort: v.ContainerPort,
			Protocol:      string(v.Protocol),
			HostIP:        v.HostIP,
		})
	}
	return res
}

//func convertProtoToSecret(req proto.Param, v proto.Service) corev1.Secret {
//	return corev1.Secret{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      v.Name,
//			Namespace: req.NameSpace,
//		},
//	}
//}

func convertSecretToProto(s *corev1.Secret) proto.Secret {
	return proto.Secret{
		Name:      s.Name,
		NameSpace: s.Namespace,
	}
}

func convertProtoToEnvVar(e []proto.EnvVar) []corev1.EnvVar {
	res := make([]corev1.EnvVar, 0)
	for _, v := range e {
		res = append(res, corev1.EnvVar{
			Name:  v.Name,
			Value: v.Value,
		})
	}
	return res
}

func convertEnvVarToProto(e []corev1.EnvVar) []proto.EnvVar {
	res := make([]proto.EnvVar, 0)
	for _, v := range e {
		res = append(res, proto.EnvVar{
			Name:  v.Name,
			Value: v.Value,
		})
	}
	return res
}

func covertHelixSagaCrdToProto(hs *helixsagaoperatorv1.HelixSaga) proto.HelixSagaCrd {
	return proto.HelixSagaCrd{
		Name:            hs.Name,
		ResourceVersion: hs.ResourceVersion,
		ConfigMap:       covertHelixSagaConfigMapVolumeToProto(hs.Spec.ConfigMap),
		Applications:    convertHelixSagaAppToProto(hs.Spec.Applications),
	}
}

func convertProtoToHelixSagaCrd(req proto.Param, hs proto.HelixSagaCrd) *helixsagaoperatorv1.HelixSaga {
	return &helixsagaoperatorv1.HelixSaga{
		ObjectMeta: metav1.ObjectMeta{
			Name:            hs.Name,
			Namespace:       req.NameSpace,
			ResourceVersion: hs.ResourceVersion,
		},
		Spec: helixsagaoperatorv1.HelixSagaSpec{
			ConfigMap:    covertProtoToHelixSagaConfigMapVolume(hs.ConfigMap),
			Applications: convertProtoToHelixSagaApp(hs.Applications),
		},
	}
}

func convertHelixSagaAppToProto(a []helixsagaoperatorv1.HelixSagaApp) []proto.HelixSagaApp {
	res := make([]proto.HelixSagaApp, 0)
	for _, v := range a {
		res = append(res, proto.HelixSagaApp{
			Spec: proto.NodeSpec{
				Name:             v.Spec.Name,
				Replicas:         *v.Spec.Replicas,
				Image:            v.Spec.Image,
				ImagePullSecrets: v.Spec.ImagePullSecrets[0].Name,
				VolumePath:       v.Spec.VolumePath,
				PodResource:      convertProtoToResourceRequirements(v.Spec.Resources),
				ContainerPorts:   convertContainerPortToProto(v.Spec.ContainerPorts),
				ServicePorts:     convertServicePortToProto(v.Spec.ServicePorts),
				Env:              convertEnvVarToProto(v.Spec.Env),
			},
			Command: v.Spec.Command,
			Args:    v.Spec.Args,
		})
	}
	return res
}

func convertProtoToHelixSagaApp(a []proto.HelixSagaApp) []helixsagaoperatorv1.HelixSagaApp {
	res := make([]helixsagaoperatorv1.HelixSagaApp, 0)
	for _, v := range a {
		res = append(res, helixsagaoperatorv1.HelixSagaApp{
			Spec: helixsagaoperatorv1.HelixSagaAppSpec{
				Name:     v.Spec.Name,
				Replicas: &v.Spec.Replicas,
				Image:    v.Spec.Image,
				ImagePullSecrets: []corev1.LocalObjectReference{
					{
						Name: v.Spec.ImagePullSecrets,
					},
				},
				VolumePath:     v.Spec.VolumePath,
				Resources:      convertResourceRequirementsToProto(v.Spec.PodResource),
				ContainerPorts: convertProtoToContainerPort(v.Spec.ContainerPorts),
				ServicePorts:   convertProtoToServicePort(v.Spec.ServicePorts),
				Env:            convertProtoToEnvVar(v.Spec.Env),
				Command:        v.Command,
				Args:           v.Args,
			},
		})
	}
	return res
}

func covertHelixSagaConfigMapVolumeToProto(c helixsagaoperatorv1.HelixSagaConfigMap) proto.HelixSagaConfigMapVolume {
	return proto.HelixSagaConfigMapVolume{
		Volume: proto.Volume{
			Name: c.Volume.Name,
			VolumeSource: proto.VolumeSource{
				Name: c.Volume.ConfigMap.Name,
				ConfigMap: &proto.ConfigMapVolumeSource{
					Items: covertKeyToPathToProto(c.Volume.ConfigMap.Items),
				},
			},
		},
		VolumeMount: proto.VolumeMount{
			Name:      c.VolumeMount.Name,
			MountPath: c.VolumeMount.MountPath,
		},
	}
}

func covertProtoToHelixSagaConfigMapVolume(c proto.HelixSagaConfigMapVolume) helixsagaoperatorv1.HelixSagaConfigMap {
	return helixsagaoperatorv1.HelixSagaConfigMap{
		Volume: corev1.Volume{
			Name: c.Volume.Name,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: c.Volume.VolumeSource.Name,
					},
					Items: covertProtoToKeyToPath(c.Volume.ConfigMap.Items),
				},
			},
		},
		VolumeMount: corev1.VolumeMount{
			Name:      c.VolumeMount.Name,
			MountPath: c.VolumeMount.MountPath,
		},
	}
}

func covertKeyToPathToProto(k []corev1.KeyToPath) []proto.KeyToPath {
	res := make([]proto.KeyToPath, 0)
	for _, v := range k {
		res = append(res, proto.KeyToPath{
			Key:  v.Key,
			Path: v.Path,
			Mode: v.Mode,
		})
	}
	return res
}

func covertProtoToKeyToPath(k []proto.KeyToPath) []corev1.KeyToPath {
	res := make([]corev1.KeyToPath, 0)
	for _, v := range k {
		res = append(res, corev1.KeyToPath{
			Key:  v.Key,
			Path: v.Path,
			Mode: v.Mode,
		})
	}
	return res
}
