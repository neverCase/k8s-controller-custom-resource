package redisoperator

import (
	"fmt"
	"regexp"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"

	redisoperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/redisoperator/v1"
)

func NewService(foo *redisoperatorv1.RedisOperator, rds *redisoperatorv1.RedisDeploymentSpec) *corev1.Service {
	var serviceName, role string
	res, err := regexp.Match(`master`, []byte(rds.DeploymentName))
	if err != nil {
		klog.V(2).Info(err)
	}
	if res {
		role = MasterName
	} else {
		role = SlaveName
	}
	var labels = map[string]string{
		"app":        operatorKindName,
		"controller": foo.Name,
		"role":       role,
	}
	serviceName = fmt.Sprintf(ServiceNameTemplate, rds.DeploymentName)
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: foo.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(foo, redisoperatorv1.SchemeGroupVersion.WithKind(operatorKindName)),
			},
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port: RedisDefaultPort,
				},
			},
			Selector: labels,
		},
	}
}
