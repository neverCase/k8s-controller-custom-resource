package mysqloperator

import (
	"github.com/Shanghai-Lunara/helixsaga-operator/pkg/serviceloadbalancer"
	corev1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8scorev1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
	mysqlOperatorV1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/mysqloperator/v1"
)

func NewService(foo *mysqlOperatorV1.MysqlOperator, rds *mysqlOperatorV1.MysqlSpec) *corev1.Service {
	var serviceName string
	var labels = map[string]string{
		k8scorev1.LabelApp:        OperatorKindName,
		k8scorev1.LabelController: foo.Name,
		k8scorev1.LabelRole:       rds.Role,
	}
	serviceName = k8scorev1.GetServiceName(rds.Name)
	ports := []corev1.ServicePort{
		{
			Port: MysqlDefaultPort,
		},
	}
	if len(rds.ServicePorts) > 0 {
		ports = rds.ServicePorts
	}
	annotations := make(map[string]string, 0)
	switch rds.ServiceType {
	case corev1.ServiceTypeLoadBalancer:
		annotations = serviceloadbalancer.Get().Annotations
		if rds.ServiceWhiteList == true {
			for k, v := range serviceloadbalancer.Get().WhiteList {
				annotations[k] = v
			}
		}
	}
	return &corev1.Service{
		ObjectMeta: metaV1.ObjectMeta{
			Annotations: annotations,
			Name:        serviceName,
			Namespace:   foo.Namespace,
			OwnerReferences: []metaV1.OwnerReference{
				*metaV1.NewControllerRef(foo, mysqlOperatorV1.SchemeGroupVersion.WithKind(OperatorKindName)),
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
