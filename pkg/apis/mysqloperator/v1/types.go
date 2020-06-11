package v1

import (
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//_ "github.com/gogo/protobuf/gogoproto"
	//_ "github.com/gogo/protobuf/proto"
	//_ "github.com/gogo/protobuf/sortkeys"
	_ "k8s.io/apimachinery/pkg/apis/testapigroup/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Mysql describes a MysqlOperator resource
type MysqlOperator struct {
	// TypeMeta is the metadata for the resource, like kind and apiversion
	metaV1.TypeMeta `json:",inline"`
	// ObjectMeta contains the metadata for the particular object, including
	// things like...
	//  - name
	//  - namespace
	//  - self link
	//  - labels
	//  - ... etc ...
	metaV1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec is the custom resource spec
	Spec MysqlOperatorSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

// MysqlSpec is the spec for a MysqlOperator resource
type MysqlOperatorSpec struct {
	MasterSpec MysqlCore `json:"masterSpec" protobuf:"bytes,1,rep,name=masterSpec"`
	SlaveSpec  MysqlCore `json:"slaveSpec" protobuf:"bytes,2,rep,name=slaveSpec"`
}

type MysqlCore struct {
	Spec   MysqlSpec   `json:"spec" protobuf:"bytes,1,rep,name=spec"`
	Status MysqlStatus `json:"status" protobuf:"bytes,2,rep,name=status"`
}

// MysqlSpec is the sub spec for a MysqlOperator resource
type MysqlSpec struct {
	// Name of the container specified as a DNS_LABEL.
	// Each container in a pod must have a unique name (DNS_LABEL).
	// Cannot be updated.
	Name string `json:"name" protobuf:"bytes,1,rep,name=name"`
	// Replicas is the number of desired replicas.
	// This is a pointer to distinguish between explicit zero and unspecified.
	// Defaults to 1.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller#what-is-a-replicationcontroller
	// +optional
	Replicas *int32 `json:"replicas,omitempty" protobuf:"bytes,2,rep,name=replicas"`
	// Docker image name.
	// More info: https://kubernetes.io/docs/concepts/containers/images
	// This field is optional to allow higher level config management to default or override
	// container images in workload controllers like Deployments and StatefulSets.
	// +optional
	Image string `json:"image,omitempty" protobuf:"bytes,3,opt,name=image"`
	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.
	// If specified, these secrets will be passed to individual puller implementations for them to use. For example,
	// in the case of docker, only DockerConfig type secrets are honored.
	// More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	ImagePullSecrets []coreV1.LocalObjectReference `json:"imagePullSecrets,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,4,rep,name=imagePullSecrets,casttype=k8s.io/api/core/v1.LocalObjectReference"`
	// List of environment variables to set in the container.
	// Cannot be updated.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	Env []coreV1.EnvVar `json:"env,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,5,rep,name=env,casttype=k8s.io/api/core/v1.EnvVar"`
	// Pod volumes to mount into the container's filesystem.
	// Cannot be updated.
	// +optional
	// +patchMergeKey=mountPath
	// +patchStrategy=merge
	VolumeMounts []coreV1.VolumeMount `json:"volumeMounts,omitempty" patchStrategy:"merge" patchMergeKey:"mountPath" protobuf:"bytes,6,rep,name=volumeMounts,casttype=k8s.io/api/core/v1.VolumeMount"`
	// List of ports to expose from the container. Exposing a port here gives
	// the system additional information about the network connections a
	// container uses, but is primarily informational. Not specifying a port here
	// DOES NOT prevent that port from being exposed. Any port which is
	// listening on the default "0.0.0.0" address inside a container will be
	// accessible from the network.
	// Cannot be updated.
	// +optional
	// +patchMergeKey=containerPort
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=containerPort
	// +listMapKey=protocol
	ContainerPorts []coreV1.ContainerPort `json:"containerPorts,omitempty" patchStrategy:"merge" patchMergeKey:"containerPort" protobuf:"bytes,7,rep,name=containerPorts,casttype=k8s.io/api/core/v1.ContainerPort"`
	// The list of ports that are exposed by this service.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
	// +patchMergeKey=port
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=port
	// +listMapKey=protocol
	ServicePorts []coreV1.ServicePort `json:"servicePorts,omitempty" patchStrategy:"merge" patchMergeKey:"port" protobuf:"bytes,8,rep,name=servicePorts,casttype=k8s.io/api/core/v1.ServicePort"`
	// The role of the server in the clusters.
	// such as: master, slave
	Role string `json:"role" protobuf:"bytes,9,rep,name=role"`
	// The config of the mysql instance.
	Config ServerConfig `json:"server_config" protobuf:"bytes,10,rep,name=serverConfig"`
}

// ServerConfig is the configuration for a MysqlDeploymentSpec of a MysqlOperator resource
type ServerConfig struct {
	ServerId    *int32 `json:"server_id" protobuf:"varint,1,opt,name=server_id,json=serverId"`
	Host        string `json:"host" protobuf:"bytes,2,opt,name=host"`
	User        string `json:"user" protobuf:"bytes,3,opt,name=user"`
	Password    string `json:"password" protobuf:"bytes,4,opt,name=password"`
	LogFile     string `json:"log_file" protobuf:"bytes,5,opt,name=log_file,json=logFile"`
	LogPosition string `json:"log_position" protobuf:"bytes,6,opt,name=log_position,json=logPosition"`
}

// MysqlSpecStatus is the status for a MysqlOperator resource
type MysqlStatus struct {
	// The generation observed by the deployment controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`

	// Total number of non-terminated pods targeted by this deployment (their labels match the selector).
	// +optional
	Replicas int32 `json:"replicas,omitempty" protobuf:"varint,2,opt,name=replicas"`

	// Total number of non-terminated pods targeted by this deployment that have the desired template spec.
	// +optional
	UpdatedReplicas int32 `json:"updatedReplicas,omitempty" protobuf:"varint,3,opt,name=updatedReplicas"`

	// Total number of ready pods targeted by this deployment.
	// +optional
	ReadyReplicas int32 `json:"readyReplicas,omitempty" protobuf:"varint,4,opt,name=readyReplicas"`

	// Total number of available pods (ready for at least minReadySeconds) targeted by this deployment.
	// +optional
	AvailableReplicas int32 `json:"availableReplicas,omitempty" protobuf:"varint,5,opt,name=availableReplicas"`

	// Total number of unavailable pods targeted by this deployment. This is the total number of
	// pods that are still required for the deployment to have 100% available capacity. They may
	// either be pods that are running but not yet available or pods that still have not been created.
	// +optional
	UnavailableReplicas int32 `json:"unavailableReplicas,omitempty" protobuf:"varint,6,opt,name=unavailableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MysqlList is a list of MysqlOperator resources
type MysqlOperatorList struct {
	metaV1.TypeMeta `json:",inline"`
	// +optional
	metaV1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []MysqlOperator `json:"items" protobuf:"bytes,2,rep,name=items"`
}
