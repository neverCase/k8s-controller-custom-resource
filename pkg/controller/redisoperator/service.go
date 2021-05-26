package redisoperator

import (
	"github.com/Shanghai-Lunara/helixsaga-operator/pkg/serviceloadbalancer"
	corev1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8scorev1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
	redisOperatorV1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/redisoperator/v1"
)

func NewService(foo *redisOperatorV1.RedisOperator, rds *redisOperatorV1.RedisSpec) *corev1.Service {
	var serviceName string
	var labels = map[string]string{
		k8scorev1.LabelApp:        OperatorKindName,
		k8scorev1.LabelController: foo.Name,
		k8scorev1.LabelRole:       rds.Role,
	}
	serviceName = k8scorev1.GetServiceName(rds.Name)
	ports := []corev1.ServicePort{
		{
			Port: RedisDefaultPort,
		},
	}
	if len(rds.ServicePorts) > 0 {
		ports = rds.ServicePorts
	}
	return &corev1.Service{
		ObjectMeta: metaV1.ObjectMeta{
			Annotations: serviceloadbalancer.Annotation(rds.ServiceType, rds.ServiceWhiteList),
			Name:        serviceName,
			Namespace:   foo.Namespace,
			OwnerReferences: []metaV1.OwnerReference{
				*metaV1.NewControllerRef(foo, redisOperatorV1.SchemeGroupVersion.WithKind(OperatorKindName)),
			},
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Type:     k8scorev1.GetServiceType(rds.ServiceType),
			Ports:    ports,
			Selector: labels,
		},
	}
}
