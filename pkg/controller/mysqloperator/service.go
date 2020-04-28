package mysqloperator

import (
	coreV1 "k8s.io/api/core/v1"

	mysqlOperatorV1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/mysqloperator/v1"
)

func NewService(foo *mysqlOperatorV1.MysqlOperator, rds *mysqlOperatorV1.MysqlDeploymentSpec) *coreV1.Service {
	return &coreV1.Service{}
}
