package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Mysql describes a MysqlOperator resource
type MysqlOperator struct {
	// TypeMeta is the metadata for the resource, like kind and apiversion
	metav1.TypeMeta `json:",inline"`
	// ObjectMeta contains the metadata for the particular object, including
	// things like...
	//  - name
	//  - namespace
	//  - self link
	//  - labels
	//  - ... etc ...
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the custom resource spec
	Spec MysqlOperatorSpec `json:"spec"`

	Status MysqlOperatorStatus `json:"status"`
}

// MysqlSpec is the spec for a MysqlOperator resource
type MysqlOperatorSpec struct {
	MasterSpec MysqlDeploymentSpec `json:"master_spec"`
	SlaveSpec  MysqlDeploymentSpec `json:"slave_spec"`
}

// MysqlDeploymentSpec is the sub spec for a MysqlOperator resource
type MysqlDeploymentSpec struct {
	DeploymentName   string      `json:"deploymentName"`
	Replicas         *int32      `json:"replicas"`
	Image            string      `json:"image"`
	ImagePullSecrets string      `json:"imagePullSecrets"`
	Configuration    MysqlConfig `json:"self_configuration"`
}

// MysqlConfig is the configuration for a MysqlDeploymentSpec of a MysqlOperator resource
type MysqlConfig struct {
	Host        string `json:"host"`
	User        string `json:"user"`
	Password    string `json:"password"`
	LogFile     string `json:"log_file"`
	LogPosition string `json:"log_position"`
}

// MysqlOperatorStatus is the status for a MysqlOperator resource
type MysqlOperatorStatus struct {
	MasterStatus MysqlDeploymentStatus `json:"master_status"`
	SlaveStatus  MysqlDeploymentStatus `json:"slave_status"`
}

// MysqlDeploymentStatus is the sub status for a MysqlOperator resource
type MysqlDeploymentStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MysqlList is a list of MysqlOperator resources
type MysqlOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []MysqlOperator `json:"items"`
}
