package handle

import (
	"reflect"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/klog/v2"

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
	Watch(broadcast chan []byte)
	Resources(req proto.Param) (res []byte, err error)
}

func NewKubernetesApiHandle(g group.Group, broadcast chan []byte) KubernetesApiInterface {
	kh := &k8sHandle{
		group: g,
	}
	go kh.Watch(broadcast)
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
		m := convertProtoToMysqlCrd(req, e)
		if n, err = resourceCreate(h.group, req, m.Name, m); err != nil {
			break
		}
		e = convertMysqlCrdToProto(n.(*mysqloperatorv1.MysqlOperator))
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
		m := convertProtoToMysqlCrd(req, e)
		if n, err = resourceUpdate(h.group, req, m.Name, m); err != nil {
			break
		}
		e = convertMysqlCrdToProto(n.(*mysqloperatorv1.MysqlOperator))
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
	case group.ServiceAccount:
		var e proto.ServiceAccount
		if err = e.Unmarshal(obj); err != nil {
			break
		}
		n, err = h.group.Resource().Get(req.ResourceType, req.NameSpace, e.Name)
		if err != nil {
			break
		}
		e = convertServiceAccountToProto(n.(*corev1.ServiceAccount))
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
		e = convertMysqlCrdToProto(n.(*mysqloperatorv1.MysqlOperator))
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
	case group.Pod:
		m := proto.PodList{
			Items: make([]proto.Pod, 0),
		}
		for _, v := range d.(*corev1.PodList).Items {
			m.Items = append(m.Items, convertPodToProto(&v))
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
	case group.ServiceAccount:
		m := proto.ServiceAccountList{
			Items: make([]proto.ServiceAccount, 0),
		}
		for _, v := range d.(*corev1.ServiceAccountList).Items {
			m.Items = append(m.Items, convertServiceAccountToProto(&v))
		}
		res, err = m.Marshal()
	case group.MysqlOperator:
		m := proto.MysqlCrdList{
			Items: make([]proto.MysqlCrd, 0),
		}
		for _, v := range d.(*mysqloperatorv1.MysqlOperatorList).Items {
			m.Items = append(m.Items, convertMysqlCrdToProto(&v))
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

func (h *k8sHandle) convertObjFromEvent(obj interface{}, et watch.EventType) (res []byte, err error) {
	req := proto.Param{
		Service:        string(proto.SvcWatch),
		WatchEventType: proto.EventType(et),
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
	case reflect.TypeOf(&corev1.Pod{}):
		var e proto.Pod
		n := obj.(*corev1.Pod)
		req.ResourceType = group.Pod
		e = convertPodToProto(n)
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
		e = convertMysqlCrdToProto(n)
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

func (h *k8sHandle) Watch(broadcast chan []byte) {
	for {
		select {
		case e, isClosed := <-h.group.WatchEvents():
			if !isClosed {
				return
			}
			res, err := h.convertObjFromEvent(e.Object, e.Type)
			if err != nil {
				klog.V(2).Info(err)
				continue
			}
			broadcast <- res
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

func convertProtoToServiceType(st corev1.ServiceType) proto.ServiceType {
	if st == "" {
		return proto.ServiceType(corev1.ServiceTypeClusterIP)
	}
	return proto.ServiceType(st)
}

func convertServiceTypeToProto(st proto.ServiceType) corev1.ServiceType {
	if st == "" {
		return corev1.ServiceTypeClusterIP
	}
	return corev1.ServiceType(st)
}

func convertPodToProto(p *corev1.Pod) proto.Pod {
	names := make([]string, 0)
	for _, v := range p.Spec.Containers {
		names = append(names, v.Name)
	}
	return proto.Pod{
		Name:            p.Name,
		Namespace:       p.Namespace,
		ResourceVersion: p.ResourceVersion,
		ContainerNames:  names,
		Status: proto.PodStatus{
			Phase:  proto.PodPhase(p.Status.Phase),
			HostIP: p.Status.HostIP,
			PodIP:  p.Status.PodIP,
		},
	}
}

func convertProtoToMysqlCrd(req proto.Param, mysqlCrd proto.MysqlCrd) *mysqloperatorv1.MysqlOperator {
	masterReplicas := mysqlCrd.Master.Replicas
	masterCollisionCount := *mysqlCrd.Master.Status.CollisionCount
	slaveReplicas := mysqlCrd.Slave.Replicas
	slaveCollisionCount := *mysqlCrd.Slave.Status.CollisionCount
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
					Replicas: &masterReplicas,
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
					ServiceType:    convertServiceTypeToProto(mysqlCrd.Master.ServiceType),
					Env:            convertProtoToEnvVar(mysqlCrd.Master.Env),
				},
				Status: mysqloperatorv1.MysqlStatus{
					ObservedGeneration: mysqlCrd.Master.Status.ObservedGeneration,
					Replicas:           mysqlCrd.Master.Status.Replicas,
					ReadyReplicas:      mysqlCrd.Master.Status.ReadyReplicas,
					CurrentReplicas:    mysqlCrd.Master.Status.CurrentReplicas,
					UpdatedReplicas:    mysqlCrd.Master.Status.UpdatedReplicas,
					CurrentRevision:    mysqlCrd.Master.Status.CurrentRevision,
					UpdateRevision:     mysqlCrd.Master.Status.UpdateRevision,
					CollisionCount:     &masterCollisionCount,
				},
			},
			SlaveSpec: mysqloperatorv1.MysqlCore{
				Spec: mysqloperatorv1.MysqlSpec{
					Name:     mysqlCrd.Slave.Name,
					Replicas: &slaveReplicas,
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
					ServiceType:    convertServiceTypeToProto(mysqlCrd.Slave.ServiceType),
					Env:            convertProtoToEnvVar(mysqlCrd.Slave.Env),
				},
				Status: mysqloperatorv1.MysqlStatus{
					ObservedGeneration: mysqlCrd.Slave.Status.ObservedGeneration,
					Replicas:           mysqlCrd.Slave.Status.Replicas,
					ReadyReplicas:      mysqlCrd.Slave.Status.ReadyReplicas,
					CurrentReplicas:    mysqlCrd.Slave.Status.CurrentReplicas,
					UpdatedReplicas:    mysqlCrd.Slave.Status.UpdatedReplicas,
					CurrentRevision:    mysqlCrd.Slave.Status.CurrentRevision,
					UpdateRevision:     mysqlCrd.Slave.Status.UpdateRevision,
					CollisionCount:     &slaveCollisionCount,
				},
			},
		},
	}
}

func convertMysqlCrdToProto(m *mysqloperatorv1.MysqlOperator) proto.MysqlCrd {
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
			ServiceType:      convertProtoToServiceType(m.Spec.MasterSpec.Spec.ServiceType),
			Env:              convertEnvVarToProto(m.Spec.MasterSpec.Spec.Env),
			Status: proto.Status{
				ObservedGeneration: m.Spec.MasterSpec.Status.ObservedGeneration,
				Replicas:           m.Spec.MasterSpec.Status.Replicas,
				ReadyReplicas:      m.Spec.MasterSpec.Status.ReadyReplicas,
				CurrentReplicas:    m.Spec.MasterSpec.Status.CurrentReplicas,
				UpdatedReplicas:    m.Spec.MasterSpec.Status.UpdatedReplicas,
				CurrentRevision:    m.Spec.MasterSpec.Status.CurrentRevision,
				UpdateRevision:     m.Spec.MasterSpec.Status.UpdateRevision,
				CollisionCount:     m.Spec.MasterSpec.Status.CollisionCount,
			},
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
			ServiceType:      convertProtoToServiceType(m.Spec.SlaveSpec.Spec.ServiceType),
			Env:              convertEnvVarToProto(m.Spec.SlaveSpec.Spec.Env),
			Status: proto.Status{
				ObservedGeneration: m.Spec.SlaveSpec.Status.ObservedGeneration,
				Replicas:           m.Spec.SlaveSpec.Status.Replicas,
				ReadyReplicas:      m.Spec.SlaveSpec.Status.ReadyReplicas,
				CurrentReplicas:    m.Spec.SlaveSpec.Status.CurrentReplicas,
				UpdatedReplicas:    m.Spec.SlaveSpec.Status.UpdatedReplicas,
				CurrentRevision:    m.Spec.SlaveSpec.Status.CurrentRevision,
				UpdateRevision:     m.Spec.SlaveSpec.Status.UpdateRevision,
				CollisionCount:     m.Spec.SlaveSpec.Status.CollisionCount,
			},
		},
	}
}

func convertProtoToRedisCrd(req proto.Param, redisCrd proto.RedisCrd) *redisoperatorv1.RedisOperator {
	masterReplicas := redisCrd.Master.Replicas
	masterCollisionCount := *redisCrd.Master.Status.CollisionCount
	slaveReplicas := redisCrd.Slave.Replicas
	slaveCollisionCount := *redisCrd.Slave.Status.CollisionCount
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
					Replicas: &masterReplicas,
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
					ServiceType:    convertServiceTypeToProto(redisCrd.Master.ServiceType),
					Env:            convertProtoToEnvVar(redisCrd.Master.Env),
				},
				Status: redisoperatorv1.RedisStatus{
					ObservedGeneration: redisCrd.Master.Status.ObservedGeneration,
					Replicas:           redisCrd.Master.Status.Replicas,
					ReadyReplicas:      redisCrd.Master.Status.ReadyReplicas,
					CurrentReplicas:    redisCrd.Master.Status.CurrentReplicas,
					UpdatedReplicas:    redisCrd.Master.Status.UpdatedReplicas,
					CurrentRevision:    redisCrd.Master.Status.CurrentRevision,
					UpdateRevision:     redisCrd.Master.Status.UpdateRevision,
					CollisionCount:     &masterCollisionCount,
				},
			},
			SlaveSpec: redisoperatorv1.RedisCore{
				Spec: redisoperatorv1.RedisSpec{
					Name:     redisCrd.Slave.Name,
					Replicas: &slaveReplicas,
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
					ServiceType:    convertServiceTypeToProto(redisCrd.Slave.ServiceType),
					Env:            convertProtoToEnvVar(redisCrd.Slave.Env),
				},
				Status: redisoperatorv1.RedisStatus{
					ObservedGeneration: redisCrd.Slave.Status.ObservedGeneration,
					Replicas:           redisCrd.Slave.Status.Replicas,
					ReadyReplicas:      redisCrd.Slave.Status.ReadyReplicas,
					CurrentReplicas:    redisCrd.Slave.Status.CurrentReplicas,
					UpdatedReplicas:    redisCrd.Slave.Status.UpdatedReplicas,
					CurrentRevision:    redisCrd.Slave.Status.CurrentRevision,
					UpdateRevision:     redisCrd.Slave.Status.UpdateRevision,
					CollisionCount:     &slaveCollisionCount,
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
			ServiceType:      convertProtoToServiceType(v.Spec.MasterSpec.Spec.ServiceType),
			Env:              convertEnvVarToProto(v.Spec.MasterSpec.Spec.Env),
			Status: proto.Status{
				ObservedGeneration: v.Spec.MasterSpec.Status.ObservedGeneration,
				Replicas:           v.Spec.MasterSpec.Status.Replicas,
				ReadyReplicas:      v.Spec.MasterSpec.Status.ReadyReplicas,
				CurrentReplicas:    v.Spec.MasterSpec.Status.CurrentReplicas,
				UpdatedReplicas:    v.Spec.MasterSpec.Status.UpdatedReplicas,
				CurrentRevision:    v.Spec.MasterSpec.Status.CurrentRevision,
				UpdateRevision:     v.Spec.MasterSpec.Status.UpdateRevision,
				CollisionCount:     v.Spec.MasterSpec.Status.CollisionCount,
			},
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
			ServiceType:      convertProtoToServiceType(v.Spec.SlaveSpec.Spec.ServiceType),
			Env:              convertEnvVarToProto(v.Spec.SlaveSpec.Spec.Env),
			Status: proto.Status{
				ObservedGeneration: v.Spec.SlaveSpec.Status.ObservedGeneration,
				Replicas:           v.Spec.SlaveSpec.Status.Replicas,
				ReadyReplicas:      v.Spec.SlaveSpec.Status.ReadyReplicas,
				CurrentReplicas:    v.Spec.SlaveSpec.Status.CurrentReplicas,
				UpdatedReplicas:    v.Spec.SlaveSpec.Status.UpdatedReplicas,
				CurrentRevision:    v.Spec.SlaveSpec.Status.CurrentRevision,
				UpdateRevision:     v.Spec.SlaveSpec.Status.UpdateRevision,
				CollisionCount:     v.Spec.SlaveSpec.Status.CollisionCount,
			},
		},
	}
}

func convertProtoToConfigMap(req proto.Param, v proto.ConfigMap) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:            v.Name,
			Namespace:       req.NameSpace,
			ResourceVersion: v.ResourceVersion,
		},
		Data: v.Data,
	}
}

func convertConfigMapToProto(c *corev1.ConfigMap) proto.ConfigMap {
	return proto.ConfigMap{
		Name:            c.Name,
		ResourceVersion: c.ResourceVersion,
		Data:            c.Data,
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
		Type:        proto.ServiceType(s.Spec.Type),
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

func convertServiceAccountToProto(s *corev1.ServiceAccount) proto.ServiceAccount {
	return proto.ServiceAccount{
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
		policy := helixsagaoperatorv1.WatchPolicyManual
		if v.Spec.WatchPolicy == helixsagaoperatorv1.WatchPolicyAuto {
			policy = helixsagaoperatorv1.WatchPolicyAuto
		}
		aft := &proto.Affinity{}

		if v.Spec.Affinity != nil &&
			v.Spec.Affinity.NodeAffinity != nil &&
			v.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution != nil &&
			len(v.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms) == 0 {
			aft = &proto.Affinity{
				NodeAffinity: &proto.NodeAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: &proto.NodeSelector{
						NodeSelectorTerms: convertNodeSelectorTermsToProto(v.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms),
					},
					PreferredDuringSchedulingIgnoredDuringExecution: make([]proto.PreferredSchedulingTerm, 0),
				},
				PodAffinity:     &proto.PodAffinity{},
				PodAntiAffinity: &proto.PodAntiAffinity{},
			}
		}
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
				ServiceType:      convertProtoToServiceType(v.Spec.ServiceType),
				Env:              convertEnvVarToProto(v.Spec.Env),
				Status: proto.Status{
					ObservedGeneration: v.Status.ObservedGeneration,
					Replicas:           v.Status.Replicas,
					ReadyReplicas:      v.Status.ReadyReplicas,
					CurrentReplicas:    v.Status.CurrentReplicas,
					UpdatedReplicas:    v.Status.UpdatedReplicas,
					CurrentRevision:    v.Status.CurrentRevision,
					UpdateRevision:     v.Status.UpdateRevision,
					CollisionCount:     v.Status.CollisionCount,
				},
			},
			Command:            v.Spec.Command,
			Args:               v.Spec.Args,
			WatchPolicy:        proto.WatchPolicy(policy),
			NodeSelector:       convertNodeSelectorElementToList(v.Spec.NodeSelector),
			ServiceAccountName: v.Spec.ServiceAccountName,
			Affinity:           aft,
			Tolerations:        convertTolerationsToProto(v.Spec.Tolerations),
		})
	}
	return res
}

func convertNodeSelectorElementToList(in map[string]string) []proto.NodeSelectorElement {
	res := make([]proto.NodeSelectorElement, 0)
	for k, v := range in {
		res = append(res, proto.NodeSelectorElement{
			Key:   k,
			Value: v,
		})
	}
	return res
}

func convertNodeSelectorTermsToProto(in []corev1.NodeSelectorTerm) []proto.NodeSelectorTerm {
	res := make([]proto.NodeSelectorTerm, len(in))
	for k, v := range in {
		res[k] = proto.NodeSelectorTerm{
			MatchExpressions: convertNodeSelectorRequirementsToProto(v.MatchExpressions),
			MatchFields:      convertNodeSelectorRequirementsToProto(v.MatchFields),
		}
	}
	return res
}

func convertProtoToNodeSelectorTerms(in []proto.NodeSelectorTerm) []corev1.NodeSelectorTerm {
	res := make([]corev1.NodeSelectorTerm, len(in))
	for k, v := range in {
		res[k] = corev1.NodeSelectorTerm{
			MatchExpressions: convertProtoToNodeSelectorRequirements(v.MatchExpressions),
			MatchFields:      convertProtoToNodeSelectorRequirements(v.MatchFields),
		}
	}
	return res
}

func convertNodeSelectorRequirementsToProto(in []corev1.NodeSelectorRequirement) []proto.NodeSelectorRequirement {
	res := make([]proto.NodeSelectorRequirement, len(in))
	for k, v := range in {
		res[k] = proto.NodeSelectorRequirement{
			Key:      v.Key,
			Operator: proto.NodeSelectorOperator(v.Operator),
			Values:   v.Values,
		}
	}
	return res
}

func convertProtoToNodeSelectorRequirements(in []proto.NodeSelectorRequirement) []corev1.NodeSelectorRequirement {
	res := make([]corev1.NodeSelectorRequirement, len(in))
	for k, v := range in {
		res[k] = corev1.NodeSelectorRequirement{
			Key:      v.Key,
			Operator: corev1.NodeSelectorOperator(v.Operator),
			Values:   v.Values,
		}
	}
	return res
}

func convertTolerationsToProto(in []corev1.Toleration) []proto.Toleration {
	res := make([]proto.Toleration, len(in))
	for k, v := range in {
		res[k] = proto.Toleration{
			Key:               v.Key,
			Operator:          proto.TolerationOperator(v.Operator),
			Value:             v.Value,
			Effect:            proto.TaintEffect(v.Effect),
			TolerationSeconds: v.TolerationSeconds,
		}
	}
	return res
}

func convertProtoToTolerations(in []proto.Toleration) []corev1.Toleration {
	res := make([]corev1.Toleration, len(in))
	for k, v := range in {
		res[k] = corev1.Toleration{
			Key:               v.Key,
			Operator:          corev1.TolerationOperator(v.Operator),
			Value:             v.Value,
			Effect:            corev1.TaintEffect(v.Effect),
			TolerationSeconds: v.TolerationSeconds,
		}
	}
	return res
}

func convertProtoToHelixSagaApp(a []proto.HelixSagaApp) []helixsagaoperatorv1.HelixSagaApp {
	res := make([]helixsagaoperatorv1.HelixSagaApp, 0)
	for _, v := range a {
		a := v.Spec.Replicas
		c := *v.Spec.Status.CollisionCount
		policy := proto.WatchPolicyManual
		if v.WatchPolicy == proto.WatchPolicyAuto {
			policy = proto.WatchPolicyAuto
		}
		aft := &corev1.Affinity{}
		if v.Affinity != nil &&
			v.Affinity.NodeAffinity != nil &&
			v.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution != nil &&
			len(v.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms) == 0 {
			aft = &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
						NodeSelectorTerms: convertProtoToNodeSelectorTerms(v.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms),
					},
					PreferredDuringSchedulingIgnoredDuringExecution: make([]corev1.PreferredSchedulingTerm, 0),
				},
				PodAffinity:     &corev1.PodAffinity{},
				PodAntiAffinity: &corev1.PodAntiAffinity{},
			}
		}
		res = append(res, helixsagaoperatorv1.HelixSagaApp{
			Spec: helixsagaoperatorv1.HelixSagaAppSpec{
				Name:     v.Spec.Name,
				Replicas: &a,
				Image:    v.Spec.Image,
				ImagePullSecrets: []corev1.LocalObjectReference{
					{
						Name: v.Spec.ImagePullSecrets,
					},
				},
				VolumePath:         v.Spec.VolumePath,
				Resources:          convertResourceRequirementsToProto(v.Spec.PodResource),
				ContainerPorts:     convertProtoToContainerPort(v.Spec.ContainerPorts),
				ServicePorts:       convertProtoToServicePort(v.Spec.ServicePorts),
				ServiceType:        convertServiceTypeToProto(v.Spec.ServiceType),
				Env:                convertProtoToEnvVar(v.Spec.Env),
				Command:            v.Command,
				Args:               v.Args,
				WatchPolicy:        helixsagaoperatorv1.WatchPolicy(policy),
				NodeSelector:       convertNodeSelectorElementToMap(v.NodeSelector),
				ServiceAccountName: v.ServiceAccountName,
				Affinity:           aft,
				Tolerations:        convertProtoToTolerations(v.Tolerations),
			},
			Status: helixsagaoperatorv1.HelixSagaAppStatus{
				ObservedGeneration: v.Spec.Status.ObservedGeneration,
				Replicas:           v.Spec.Status.Replicas,
				ReadyReplicas:      v.Spec.Status.ReadyReplicas,
				CurrentReplicas:    v.Spec.Status.CurrentReplicas,
				UpdatedReplicas:    v.Spec.Status.UpdatedReplicas,
				CurrentRevision:    v.Spec.Status.CurrentRevision,
				UpdateRevision:     v.Spec.Status.UpdateRevision,
				CollisionCount:     &c,
			},
		})
	}
	return res
}

func convertNodeSelectorElementToMap(in []proto.NodeSelectorElement) map[string]string {
	res := make(map[string]string, 0)
	for _, v := range in {
		res[v.Key] = v.Value
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
