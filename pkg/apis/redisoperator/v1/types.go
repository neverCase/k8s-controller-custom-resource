package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Redis describes a RedisOperator resource
type RedisOperator struct {
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
	Spec RedisOperatorSpec `json:"spec"`

	Status RedisOperatorStatus `json:"status"`
}

// RedisSpec is the spec for a RedisOperator resource
type RedisOperatorSpec struct {
	MasterSpec RedisDeploymentSpec `json:"master_spec"`
	SlaveSpec  RedisDeploymentSpec `json:"slave_spec"`
}

// RedisDeploymentSpec is the sub spec for a RedisOperator resource
type RedisDeploymentSpec struct {
	DeploymentName   string `json:"deploymentName"`
	Replicas         *int32 `json:"replicas"`
	Image            string `json:"image"`
	ImagePullSecrets string `json:"imagePullSecrets"`
}

// RedisOperatorStatus is the status for a RedisOperator resource
type RedisOperatorStatus struct {
	MasterStatus RedisDeploymentStatus `json:"master_status"`
	SlaveStatus  RedisDeploymentStatus `json:"slave_status"`
}

// RedisDeploymentStatus is the sub status for a RedisOperator resource
type RedisDeploymentStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RedisList is a list of RedisOperator resources
type RedisOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []RedisOperator `json:"items"`
}
