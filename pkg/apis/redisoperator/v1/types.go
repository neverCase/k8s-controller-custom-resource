package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Redis describes a Redis resource
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

// RedisSpec is the spec for a Network resource
type RedisOperatorSpec struct {
	DeploymentName string `json:"deploymentName"`
	MasterReplicas *int32 `json:"master_replicas"`
	SlaveReplicas  *int32 `json:"slave_replicas"`
}

// RedisStatus is the status for a Redis resource
type RedisOperatorStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RedisList is a list of Redis resources
type RedisOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []RedisOperator `json:"items"`
}
