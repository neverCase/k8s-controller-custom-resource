package mysqloperator

import (
	"fmt"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sCoreV1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
	mysqlOperatorV1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/mysqloperator/v1"
)

func NewService(foo *mysqlOperatorV1.MysqlOperator, rds *mysqlOperatorV1.MysqlSpec) *coreV1.Service {
	var serviceName string
	var labels = map[string]string{
		"app":        OperatorKindName,
		"controller": foo.Name,
		"role":       rds.Role,
	}
	serviceName = fmt.Sprintf(k8sCoreV1.ServiceNameTemplate, rds.Name)
	ports := []coreV1.ServicePort{
		{
			Port: MysqlDefaultPort,
		},
	}
	if len(rds.ServicePorts) > 0 {
		ports = rds.ServicePorts
	}
	return &coreV1.Service{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      serviceName,
			Namespace: foo.Namespace,
			OwnerReferences: []metaV1.OwnerReference{
				*metaV1.NewControllerRef(foo, mysqlOperatorV1.SchemeGroupVersion.WithKind(OperatorKindName)),
			},
			Labels: labels,
		},
		Spec: coreV1.ServiceSpec{
			Ports:    ports,
			Selector: labels,
		},
	}
}
