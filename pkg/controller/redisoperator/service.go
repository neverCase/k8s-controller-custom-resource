package redisoperator

import (
	"fmt"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sCoreV1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
	redisOperatorV1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/redisoperator/v1"
)

func NewService(foo *redisOperatorV1.RedisOperator, rds *redisOperatorV1.RedisSpec) *coreV1.Service {
	var serviceName string
	var labels = map[string]string{
		"app":        operatorKindName,
		"controller": foo.Name,
		"role":       rds.Role,
	}
	serviceName = fmt.Sprintf(k8sCoreV1.ServiceNameTemplate, rds.Name)
	return &coreV1.Service{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      serviceName,
			Namespace: foo.Namespace,
			OwnerReferences: []metaV1.OwnerReference{
				*metaV1.NewControllerRef(foo, redisOperatorV1.SchemeGroupVersion.WithKind(operatorKindName)),
			},
			Labels: labels,
		},
		Spec: coreV1.ServiceSpec{
			Ports: []coreV1.ServicePort{
				{
					Port: RedisDefaultPort,
				},
			},
			Selector: labels,
		},
	}
}
