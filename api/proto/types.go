package proto

import (
	"github.com/nevercase/k8s-controller-custom-resource/api/group"
)

type ApiService string

const (
	SvcPing     ApiService = "ping"
	SvcCreate   ApiService = "create"
	SvcUpdate   ApiService = "update"
	SvcDelete   ApiService = "delete"
	SvcGet      ApiService = "get"
	SvcList     ApiService = "list"
	SvcWatch    ApiService = "watch"
	SvcResource ApiService = "resource"
	SvcHarbor   ApiService = "harbor"
)

type HarborCommand string

const (
	Hubs         HarborCommand = "hubs"
	Projects     HarborCommand = "projects"
	Repositories HarborCommand = "repositories"
	Tags         HarborCommand = "tags"
)

const (
	CodeNone = iota
	CodeErr  = 10001
)

// EventType defines the possible types of events.
type EventType string

const (
	EventAdded    EventType = "ADDED"
	EventModified EventType = "MODIFIED"
	EventDeleted  EventType = "DELETED"
	EventBookmark EventType = "BOOKMARK"
	EventError    EventType = "ERROR"
)

type Param struct {
	Service      string             `json:"service" protobuf:"bytes,1,opt,name=service"`
	ResourceType group.ResourceType `json:"resourceType" protobuf:"bytes,2,opt,name=resourceType"`
	NameSpace    string             `json:"nameSpace" protobuf:"bytes,3,opt,name=nameSpace"`

	WatchEventType EventType `json:"watchEventType" protobuf:"bytes,4,opt,name=watchEventType"`

	// The request parameters of the HarborApi
	// only used for HarborApi
	HarborRequest HarborRequest `json:"harborRequest" protobuf:"bytes,5,opt,name=harborRequest"`
}

type Request struct {
	Param Param  `protobuf:"bytes,1,opt,name=param"`
	Data  []byte `json:"data" protobuf:"bytes,2,opt,name=data"`
}

type Response struct {
	Code   int32  `json:"code" protobuf:"varint,1,opt,name=code"`
	Param  Param  `protobuf:"bytes,2,opt,name=param"`
	Result []byte `json:"result" protobuf:"bytes,3,opt,name=result"`
}

type MysqlCrdList struct {
	Items []MysqlCrd `json:"items" protobuf:"bytes,1,rep,name=items"`
}

type MysqlCrd struct {
	Name string `json:"name" protobuf:"bytes,1,rep,name=name"`
	// An opaque value that represents the internal version of this object that can
	// be used by clients to determine when objects have changed. May be used for optimistic
	// concurrency, change detection, and the watch operation on a resource or set of resources.
	// Clients must treat these values as opaque and passed unmodified back to the server.
	// They may only be valid for a particular resource or set of resources.
	//
	// Populated by the system.
	// Read-only.
	// Value must be treated as opaque by clients and .
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency
	// +optional
	ResourceVersion string   `json:"resourceVersion,omitempty" protobuf:"bytes,2,opt,name=resourceVersion"`
	Master          NodeSpec `json:"master" protobuf:"bytes,3,rep,name=master"`
	Slave           NodeSpec `json:"slave" protobuf:"bytes,4,rep,name=slave"`
}

type RedisCrdList struct {
	Items []RedisCrd `json:"items" protobuf:"bytes,1,rep,name=items"`
}

type RedisCrd struct {
	Name string `json:"name" protobuf:"bytes,1,rep,name=name"`
	// An opaque value that represents the internal version of this object that can
	// be used by clients to determine when objects have changed. May be used for optimistic
	// concurrency, change detection, and the watch operation on a resource or set of resources.
	// Clients must treat these values as opaque and passed unmodified back to the server.
	// They may only be valid for a particular resource or set of resources.
	//
	// Populated by the system.
	// Read-only.
	// Value must be treated as opaque by clients and .
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency
	// +optional
	ResourceVersion string   `json:"resourceVersion,omitempty" protobuf:"bytes,2,opt,name=resourceVersion"`
	Master          NodeSpec `json:"master" protobuf:"bytes,3,rep,name=master"`
	Slave           NodeSpec `json:"slave" protobuf:"bytes,4,rep,name=slave"`
}

type NodeSpec struct {
	Name             string `json:"name" protobuf:"bytes,1,rep,name=name"`
	Replicas         int32  `json:"replicas" protobuf:"varint,2,opt,name=replicas"`
	Image            string `json:"image" protobuf:"bytes,3,rep,name=image"`
	ImagePullSecrets string `json:"imagePullSecrets" protobuf:"bytes,4,rep,name=imagePullSecrets"`
	// The path of the nas disk which was mounted on the machine
	VolumePath string `json:"volumePath" protobuf:"bytes,5,rep,name=volumePath"`
	// PodResource
	PodResource PodResourceRequirements `json:"podResource" protobuf:"bytes,6,rep,name=podResource"`
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
	ContainerPorts []ContainerPort `json:"containerPorts,omitempty" patchStrategy:"merge" patchMergeKey:"containerPort" protobuf:"bytes,7,rep,name=containerPorts"`
	// The list of ports that are exposed by this service.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
	// +patchMergeKey=port
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=port
	// +listMapKey=protocol
	ServicePorts []ServicePort `json:"servicePorts,omitempty" patchStrategy:"merge" patchMergeKey:"port" protobuf:"bytes,8,rep,name=servicePorts"`
	// The type of the service
	ServiceType ServiceType `json:"serviceType" protobuf:"bytes,9,rep,name=serviceType"`
	// List of environment variables to set in the container.
	// Cannot be updated.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	Env []EnvVar `json:"env,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,10,rep,name=env"`
	// The status of the specific StatefulSet which was owned by the crd
	// +optional
	Status Status `json:"status" protobuf:"bytes,11,rep,name=status"`
}

type EnvVar struct {
	// Name of the environment variable. Must be a C_IDENTIFIER.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`

	// Optional: no more than one of the following may be specified.

	// Variable references $(VAR_NAME) are expanded
	// using the previous defined environment variables in the container and
	// any service environment variables. If a variable cannot be resolved,
	// the reference in the input string will be unchanged. The $(VAR_NAME)
	// syntax can be escaped with a double $$, ie: $$(VAR_NAME). Escaped
	// references will never be expanded, regardless of whether the variable
	// exists or not.
	// Defaults to "".
	// +optional
	Value string `json:"value,omitempty" protobuf:"bytes,2,opt,name=value"`
	// Source for the environment variable's value. Cannot be used if value is not empty.
	// +optional
	//ValueFrom *corev1.EnvVarSource `json:"valueFrom,omitempty" protobuf:"bytes,3,opt,name=valueFrom"`
}

// PodResourceRequirements describes the compute resource requirements.
type PodResourceRequirements struct {
	// Limits describes the maximum amount of compute resources allowed.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/
	// +optional
	Limits PodResourceList `json:"limits,omitempty" protobuf:"bytes,1,rep,name=limits"`
	// Requests describes the minimum amount of compute resources required.
	// If Requests is omitted for a container, it defaults to Limits if that is explicitly specified,
	// otherwise to an implementation-defined value.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/
	// +optional
	Requests PodResourceList `json:"requests,omitempty" protobuf:"bytes,2,rep,name=requests"`
}

// PodResourceList is a set of (resource name, quantity) pairs.
type PodResourceList map[string]string

type ResourceList struct {
	Items []group.ResourceType `json:"items" protobuf:"bytes,1,rep,name=items"`
}

type ConfigMapList struct {
	Items []ConfigMap `json:"items" protobuf:"bytes,1,rep,name=items"`
}

type ConfigMap struct {
	Name string `json:"name" protobuf:"bytes,1,rep,name=Name"`
	// An opaque value that represents the internal version of this object that can
	// be used by clients to determine when objects have changed. May be used for optimistic
	// concurrency, change detection, and the watch operation on a resource or set of resources.
	// Clients must treat these values as opaque and passed unmodified back to the server.
	// They may only be valid for a particular resource or set of resources.
	//
	// Populated by the system.
	// Read-only.
	// Value must be treated as opaque by clients and .
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency
	// +optional
	ResourceVersion string `json:"resourceVersion,omitempty" protobuf:"bytes,2,opt,name=resourceVersion"`
	// Data contains the configuration data.
	// Each key must consist of alphanumeric characters, '-', '_' or '.'.
	// Values with non-UTF-8 byte sequences must use the BinaryData field.
	// The keys stored in Data must not overlap with the keys in
	// the BinaryData field, this is enforced during validation process.
	// +optional
	Data map[string]string `json:"data" protobuf:"bytes,3,rep,name=data"`
}

type NameSpaceList struct {
	Items []NameSpace `json:"items" protobuf:"bytes,1,rep,name=items"`
}

type NameSpace struct {
	Name string `json:"name" protobuf:"bytes,1,rep,name=Name"`
}

type ServiceList struct {
	Items []Service `json:"items" protobuf:"bytes,1,rep,name=items"`
}

type Service struct {
	Name string `json:"name" protobuf:"bytes,1,rep,name=Name"`
	// The list of ports that are exposed by this service.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
	// +patchMergeKey=port
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=port
	// +listMapKey=protocol
	Ports []ServicePort `json:"ports,omitempty" patchStrategy:"merge" patchMergeKey:"port" protobuf:"bytes,2,rep,name=ports"`
	// clusterIP is the IP address of the service and is usually assigned
	// randomly by the master. If an address is specified manually and is not in
	// use by others, it will be allocated to the service; otherwise, creation
	// of the service will fail. This field can not be changed through updates.
	// Valid values are "None", empty string (""), or a valid IP address. "None"
	// can be specified for headless services when proxying is not required.
	// Only applies to types ClusterIP, NodePort, and LoadBalancer. Ignored if
	// type is ExternalName.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
	// +optional
	ClusterIP string `json:"clusterIP,omitempty" protobuf:"bytes,3,opt,name=clusterIP"`
	// type determines how the Service is exposed. Defaults to ClusterIP. Valid
	// options are ExternalName, ClusterIP, NodePort, and LoadBalancer.
	// "ExternalName" maps to the specified externalName.
	// "ClusterIP" allocates a cluster-internal IP address for load-balancing to
	// endpoints. Endpoints are determined by the selector or if that is not
	// specified, by manual construction of an Endpoints object. If clusterIP is
	// "None", no virtual IP is allocated and the endpoints are published as a
	// set of endpoints rather than a stable IP.
	// "NodePort" builds on ClusterIP and allocates a port on every node which
	// routes to the clusterIP.
	// "LoadBalancer" builds on NodePort and creates an
	// external load-balancer (if supported in the current cloud) which routes
	// to the clusterIP.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types
	// +optional
	Type ServiceType `json:"type,omitempty" protobuf:"bytes,4,opt,name=type"`
	// externalIPs is a list of IP addresses for which nodes in the cluster
	// will also accept traffic for this service.  These IPs are not managed by
	// Kubernetes.  The user is responsible for ensuring that traffic arrives
	// at a node with this IP.  A common example is external load-balancers
	// that are not part of the Kubernetes system.
	// +optional
	ExternalIPs []string `json:"externalIPs,omitempty" protobuf:"bytes,5,rep,name=externalIPs"`
}

// ContainerPort represents a network port in a single container.
type ContainerPort struct {
	// If specified, this must be an IANA_SVC_NAME and unique within the pod. Each
	// named port in a pod must have a unique name. Name for the port that can be
	// referred to by services.
	// +optional
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	// Number of port to expose on the host.
	// If specified, this must be a valid port number, 0 < x < 65536.
	// If HostNetwork is specified, this must match ContainerPort.
	// Most containers do not need this.
	// +optional
	HostPort int32 `json:"hostPort,omitempty" protobuf:"varint,2,opt,name=hostPort"`
	// Number of port to expose on the pod's IP address.
	// This must be a valid port number, 0 < x < 65536.
	ContainerPort int32 `json:"containerPort" protobuf:"varint,3,opt,name=containerPort"`
	// Protocol for port. Must be UDP, TCP, or SCTP.
	// Defaults to "TCP".
	// +optional
	Protocol string `json:"protocol,omitempty" protobuf:"bytes,4,opt,name=protocol"`
	// What host IP to bind the external port to.
	// +optional
	HostIP string `json:"hostIP,omitempty" protobuf:"bytes,5,opt,name=hostIP"`
}

type ServicePort struct {
	// The name of this port within the service. This must be a DNS_LABEL.
	// All ports within a ServiceSpec must have unique names. When considering
	// the endpoints for a Service, this must match the 'name' field in the
	// EndpointPort.
	// Optional if only one ServicePort is defined on this service.
	// +optional
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	// The IP protocol for this port. Supports "TCP", "UDP", and "SCTP".
	// Default is TCP.
	// +optional
	Protocol string `json:"protocol,omitempty" protobuf:"bytes,2,opt,name=protocol"`
	// The port that will be exposed by this service.
	Port int32 `json:"port" protobuf:"varint,3,opt,name=port"`
	// Number or name of the port to access on the pods targeted by the service.
	// Number must be in the range 1 to 65535. Name must be an IANA_SVC_NAME.
	// If this is a string, it will be looked up as a named port in the
	// target Pod's container ports. If this is not specified, the value
	// of the 'port' field is used (an identity map).
	// This field is ignored for services with clusterIP=None, and should be
	// omitted or set equal to the 'port' field.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#defining-a-service
	// +optional
	TargetPort IntOrString `json:"targetPort,omitempty" protobuf:"bytes,4,opt,name=targetPort"`
	// The port on each node on which this service is exposed when type=NodePort or LoadBalancer.
	// Usually assigned by the system. If specified, it will be allocated to the service
	// if unused or else creation of the service will fail.
	// Default is to auto-allocate a port if the ServiceType of this Service requires one.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#type-nodeport
	// +optional
	NodePort int32 `json:"nodePort,omitempty" protobuf:"varint,5,opt,name=nodePort"`
}

// Service Type string describes ingress methods for a service
type ServiceType string

const (
	// ServiceTypeClusterIP means a service will only be accessible inside the
	// cluster, via the cluster IP.
	ServiceTypeClusterIP ServiceType = "ClusterIP"

	// ServiceTypeNodePort means a service will be exposed on one port of
	// every node, in addition to 'ClusterIP' type.
	ServiceTypeNodePort ServiceType = "NodePort"

	// ServiceTypeLoadBalancer means a service will be exposed via an
	// external load balancer (if the cloud provider supports it), in addition
	// to 'NodePort' type.
	ServiceTypeLoadBalancer ServiceType = "LoadBalancer"

	// ServiceTypeExternalName means a service consists of only a reference to
	// an external name that kubedns or equivalent will return as a CNAME
	// record, with no exposing or proxying of any pods involved.
	ServiceTypeExternalName ServiceType = "ExternalName"
)

type IntOrString struct {
	Type   int32  `protobuf:"varint,1,opt,name=type"`
	IntVal int32  `protobuf:"varint,2,opt,name=intVal"`
	StrVal string `protobuf:"bytes,3,opt,name=strVal"`
}

type Secret struct {
	Name      string `json:"name" protobuf:"bytes,1,opt,name=name"`
	NameSpace string `json:"nameSpace" protobuf:"bytes,2,opt,name=nameSpace"`
}

type SecretList struct {
	Items []Secret `json:"items" protobuf:"bytes,1,rep,name=items"`
}

type ServiceAccount struct {
	Name      string `json:"name" protobuf:"bytes,1,opt,name=name"`
	NameSpace string `json:"nameSpace" protobuf:"bytes,2,opt,name=nameSpace"`
}

type ServiceAccountList struct {
	Items []ServiceAccount `json:"items" protobuf:"bytes,1,rep,name=items"`
}

// The shortage ot the harbor projects
type HarborHub struct {
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
}

type HarborHubList struct {
	Items []HarborHub `json:"items" protobuf:"bytes,1,rep,name=items"`
}

type HarborProject struct {
	ProjectID int32  `json:"projectId" protobuf:"varint,1,opt,name=projectId"`
	Name      string `json:"name" protobuf:"bytes,2,opt,name=name"`
}

type HarborProjectList struct {
	Items []HarborProject `json:"items" protobuf:"bytes,1,rep,name=items"`
}

type HarborRepository struct {
	RepositoryID int32  `json:"repositoryId" protobuf:"bytes,1,opt,name=repositoryId"`
	Name         string `json:"name" protobuf:"bytes,2,opt,name=name"`
	ProjectID    int32  `json:"projectId" protobuf:"varint,3,opt,name=projectId"`
}

type HarborRepositoryList struct {
	Items []HarborRepository `json:"items" protobuf:"bytes,1,rep,name=items"`
}

type HarborTag struct {
	Digest string `json:"digest" protobuf:"bytes,1,rep,name=digest"`
	Name   string `json:"name" protobuf:"bytes,2,opt,name=name"`
}

type HarborTagList struct {
	Items []HarborTag `json:"items" protobuf:"bytes,1,rep,name=items"`
}

type HarborRequest struct {
	HarborUrl string        `json:"harborUrl" protobuf:"bytes,1,opt,name=harborUrl"`
	Command   HarborCommand `json:"command" protobuf:"bytes,2,opt,name=command"`
	ProjectID int32         `json:"projectId" protobuf:"varint,3,opt,name=projectId"`
	ImageName string        `json:"imageName" protobuf:"bytes,4,opt,name=imageName"`
}

// HelixSaga Operator
type HelixSagaCrd struct {
	Name string `json:"name" protobuf:"bytes,1,rep,name=name"`
	// An opaque value that represents the internal version of this object that can
	// be used by clients to determine when objects have changed. May be used for optimistic
	// concurrency, change detection, and the watch operation on a resource or set of resources.
	// Clients must treat these values as opaque and passed unmodified back to the server.
	// They may only be valid for a particular resource or set of resources.
	//
	// Populated by the system.
	// Read-only.
	// Value must be treated as opaque by clients and .
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency
	// +optional
	ResourceVersion string                   `json:"resourceVersion,omitempty" protobuf:"bytes,2,opt,name=resourceVersion"`
	ConfigMap       HelixSagaConfigMapVolume `json:"configMap" protobuf:"bytes,3,rep,name=configMap"`
	Applications    []HelixSagaApp           `json:"applications" protobuf:"bytes,4,rep,name=applications"`
}

type WatchPolicy string

const (
	WatchPolicyAuto   WatchPolicy = "auto"
	WatchPolicyManual WatchPolicy = "manual"
)

type NodeSelectorElement struct {
	Key   string `json:"key" protobuf:"bytes,1,name=key"`
	Value string `json:"value" protobuf:"bytes,2,name=value"`
}

type HelixSagaApp struct {
	Spec NodeSpec `json:"spec" protobuf:"bytes,1,rep,name=spec"`
	// Entrypoint array. Not executed within a shell.
	// The docker image's ENTRYPOINT is used if this is not provided.
	// Variable references $(VAR_NAME) are expanded using the container's environment. If a variable
	// cannot be resolved, the reference in the input string will be unchanged. The $(VAR_NAME) syntax
	// can be escaped with a double $$, ie: $$(VAR_NAME). Escaped references will never be expanded,
	// regardless of whether the variable exists or not.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
	// +optional
	Command []string `json:"command,omitempty" protobuf:"bytes,2,rep,name=command"`
	// Arguments to the entrypoint.
	// The docker image's CMD is used if this is not provided.
	// Variable references $(VAR_NAME) are expanded using the container's environment. If a variable
	// cannot be resolved, the reference in the input string will be unchanged. The $(VAR_NAME) syntax
	// can be escaped with a double $$, ie: $$(VAR_NAME). Escaped references will never be expanded,
	// regardless of whether the variable exists or not.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
	// +optional
	Args []string `json:"args,omitempty" protobuf:"bytes,3,rep,name=args"`
	// Watch policy for the present app.
	// One of Auto, Manual.
	// Default to Manual.
	WatchPolicy WatchPolicy `json:"watchPolicy" protobuf:"bytes,4,rep,name=watchPolicy"`
	// NodeSelector is a selector which must be true for the pod to fit on a node.
	// Selector which must match a node's labels for the pod to be scheduled on that node.
	// More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
	// +optional
	NodeSelector []NodeSelectorElement `json:"nodeSelector,omitempty" protobuf:"bytes,5,rep,name=nodeSelector"`
	// ServiceAccountName is the name of the ServiceAccount to use to run this pod.
	// More info: https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty" protobuf:"bytes,6,opt,name=serviceAccountName"`
	// If specified, the pod's scheduling constraints
	// +optional
	Affinity *Affinity `json:"affinity,omitempty" protobuf:"bytes,7,opt,name=affinity"`
	// If specified, the pod's tolerations.
	// +optional
	Tolerations []Toleration `json:"tolerations,omitempty" protobuf:"bytes,8,opt,name=tolerations"`
}

type HelixSagaCrdList struct {
	Items []HelixSagaCrd `json:"items" protobuf:"bytes,1,rep,name=items"`
}

type HelixSagaConfigMapVolume struct {
	Volume      Volume      `json:"volume" protobuf:"bytes,1,rep,name=volume"`
	VolumeMount VolumeMount `json:"volumeMount" protobuf:"bytes,2,rep,name=volumeMount"`
}

// Volume represents a named volume in a pod that may be accessed by any container in the pod.
type Volume struct {
	// Volume's name.
	// Must be a DNS_LABEL and unique within the pod.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// VolumeSource represents the location and type of the mounted volume.
	// If not specified, the Volume is implied to be an EmptyDir.
	// This implied behavior is deprecated and will be removed in a future version.
	VolumeSource `json:",inline" protobuf:"bytes,2,opt,name=volumeSource"`
}

// Represents the source of a volume to mount.
// Only one of its members may be specified.
type VolumeSource struct {
	// ConfigMap's name.
	// Must be a DNS_LABEL and unique within the pod.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// ConfigMap represents a configMap that should populate this volume
	// +optional
	ConfigMap *ConfigMapVolumeSource `json:"configMap,omitempty" protobuf:"bytes,2,opt,name=configMap"`
}

// Adapts a ConfigMap into a volume.
//
// The contents of the target ConfigMap's Data field will be presented in a
// volume as files using the keys in the Data field as the file names, unless
// the items element is populated with specific mappings of keys to paths.
// ConfigMap volumes support ownership management and SELinux relabeling.
type ConfigMapVolumeSource struct {
	//LocalObjectReference `json:",inline" protobuf:"bytes,1,opt,name=localObjectReference"`
	// If unspecified, each key-value pair in the Data field of the referenced
	// ConfigMap will be projected into the volume as a file whose name is the
	// key and content is the value. If specified, the listed keys will be
	// projected into the specified paths, and unlisted keys will not be
	// present. If a key is specified which is not present in the ConfigMap,
	// the volume setup will error unless it is marked optional. Paths must be
	// relative and may not contain the '..' path or start with '..'.
	// +optional
	Items []KeyToPath `json:"items,omitempty" protobuf:"bytes,1,rep,name=items"`
	// Optional: mode bits to use on created files by default. Must be a
	// value between 0 and 0777. Defaults to 0644.
	// Directories within the path are not affected by this setting.
	// This might be in conflict with other options that affect the file
	// mode, like fsGroup, and the result can be other mode bits set.
	// +optional
	DefaultMode *int32 `json:"defaultMode,omitempty" protobuf:"varint,2,opt,name=defaultMode"`
	// Specify whether the ConfigMap or its keys must be defined
	// +optional
	Optional *bool `json:"optional,omitempty" protobuf:"varint,3,opt,name=optional"`
}

// Maps a string key to a path within a volume.
type KeyToPath struct {
	// The key to project.
	Key string `json:"key" protobuf:"bytes,1,opt,name=key"`

	// The relative path of the file to map the key to.
	// May not be an absolute path.
	// May not contain the path element '..'.
	// May not start with the string '..'.
	Path string `json:"path" protobuf:"bytes,2,opt,name=path"`
	// Optional: mode bits to use on this file, must be a value between 0
	// and 0777. If not specified, the volume defaultMode will be used.
	// This might be in conflict with other options that affect the file
	// mode, like fsGroup, and the result can be other mode bits set.
	// +optional
	Mode *int32 `json:"mode,omitempty" protobuf:"varint,3,opt,name=mode"`
}

// VolumeMount describes a mounting of a Volume within a container.
type VolumeMount struct {
	// This must match the Name of a Volume.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Mounted read-only if true, read-write otherwise (false or unspecified).
	// Defaults to false.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,2,opt,name=readOnly"`
	// Path within the container at which the volume should be mounted.  Must
	// not contain ':'.
	MountPath string `json:"mountPath" protobuf:"bytes,3,opt,name=mountPath"`
	// Path within the volume from which the container's volume should be mounted.
	// Defaults to "" (volume's root).
	// +optional
	SubPath string `json:"subPath,omitempty" protobuf:"bytes,4,opt,name=subPath"`
	// mountPropagation determines how mounts are propagated from the host
	// to container and the other way around.
	// When not set, MountPropagationNone is used.
	// This field is beta in 1.10.
	// +optional
	//MountPropagation *MountPropagationMode `json:"mountPropagation,omitempty" protobuf:"bytes,5,opt,name=mountPropagation,casttype=MountPropagationMode"`
	// Expanded path within the volume from which the container's volume should be mounted.
	// Behaves similarly to SubPath but environment variable references $(VAR_NAME) are expanded using the container's environment.
	// Defaults to "" (volume's root).
	// SubPathExpr and SubPath are mutually exclusive.
	// +optional
	SubPathExpr string `json:"subPathExpr,omitempty" protobuf:"bytes,6,opt,name=subPathExpr"`
}

// Status
type Status struct {
	// observedGeneration is the most recent generation observed for this StatefulSet. It corresponds to the
	// StatefulSet's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`

	// replicas is the number of Pods created by the StatefulSet controller.
	Replicas int32 `json:"replicas" protobuf:"varint,2,opt,name=replicas"`

	// readyReplicas is the number of Pods created by the StatefulSet controller that have a Ready Condition.
	ReadyReplicas int32 `json:"readyReplicas,omitempty" protobuf:"varint,3,opt,name=readyReplicas"`

	// currentReplicas is the number of Pods created by the StatefulSet controller from the StatefulSet version
	// indicated by currentRevision.
	CurrentReplicas int32 `json:"currentReplicas,omitempty" protobuf:"varint,4,opt,name=currentReplicas"`

	// updatedReplicas is the number of Pods created by the StatefulSet controller from the StatefulSet version
	// indicated by updateRevision.
	UpdatedReplicas int32 `json:"updatedReplicas,omitempty" protobuf:"varint,5,opt,name=updatedReplicas"`

	// currentRevision, if not empty, indicates the version of the StatefulSet used to generate Pods in the
	// sequence [0,currentReplicas).
	CurrentRevision string `json:"currentRevision,omitempty" protobuf:"bytes,6,opt,name=currentRevision"`

	// updateRevision, if not empty, indicates the version of the StatefulSet used to generate Pods in the sequence
	// [replicas-updatedReplicas,replicas)
	UpdateRevision string `json:"updateRevision,omitempty" protobuf:"bytes,7,opt,name=updateRevision"`

	// collisionCount is the count of hash collisions for the StatefulSet. The StatefulSet controller
	// uses this field as a collision avoidance mechanism when it needs to create the name for the
	// newest ControllerRevision.
	// +optional
	CollisionCount *int32 `json:"collisionCount,omitempty" protobuf:"varint,9,opt,name=collisionCount"`
}

// Pod is a collection of containers that can run on a host. This resource is created
// by clients and scheduled onto hosts.
type Pod struct {
	// Name must be unique within a namespace. Is required when creating resources, although
	// some resources may allow a client to request the generation of an appropriate name
	// automatically. Name is primarily intended for creation idempotence and configuration
	// definition.
	// Cannot be updated.
	// More info: http://kubernetes.io/docs/user-guide/identifiers#names
	// +optional
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	// Namespace defines the space within each name must be unique. An empty namespace is
	// equivalent to the "default" namespace, but "default" is the canonical representation.
	// Not all objects are required to be scoped to a namespace - the value of this field for
	// those objects will be empty.
	//
	// Must be a DNS_LABEL.
	// Cannot be updated.
	// More info: http://kubernetes.io/docs/user-guide/namespaces
	// +optional
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,2,opt,name=namespace"`
	// An opaque value that represents the internal version of this object that can
	// be used by clients to determine when objects have changed. May be used for optimistic
	// concurrency, change detection, and the watch operation on a resource or set of resources.
	// Clients must treat these values as opaque and passed unmodified back to the server.
	// They may only be valid for a particular resource or set of resources.
	//
	// Populated by the system.
	// Read-only.
	// Value must be treated as opaque by clients and .
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency
	// +optional
	ResourceVersion string `json:"resourceVersion,omitempty" protobuf:"bytes,3,opt,name=resourceVersion"`
	// ContainerNames are the list of all the containers' names inside the pod
	ContainerNames []string `json:"containerNames" protobuf:"bytes,4,opt,name=containerNames"`
	// Most recently observed status of the pod.
	// This data may not be up to date.
	// Populated by the system.
	// Read-only.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Status PodStatus `json:"status" protobuf:"bytes,5,opt,name=status"`
}

type PodStatus struct {
	// +optional
	Phase PodPhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=PodPhase"`
	// IP address of the host to which the pod is assigned. Empty if not yet scheduled.
	// +optional
	HostIP string `json:"hostIP,omitempty" protobuf:"bytes,2,opt,name=hostIP"`
	// IP address allocated to the pod. Routable at least within the cluster.
	// Empty if not yet allocated.
	// +optional
	PodIP string `json:"podIP,omitempty" protobuf:"bytes,3,opt,name=podIP"`
}

// PodPhase is a label for the condition of a pod at the current time.
type PodPhase string

// These are the valid statuses of pods.
const (
	// PodPending means the pod has been accepted by the system, but one or more of the containers
	// has not been started. This includes time before being bound to a node, as well as time spent
	// pulling images onto the host.
	PodPending PodPhase = "Pending"
	// PodRunning means the pod has been bound to a node and all of the containers have been started.
	// At least one container is still running or is in the process of being restarted.
	PodRunning PodPhase = "Running"
	// PodSucceeded means that all containers in the pod have voluntarily terminated
	// with a container exit code of 0, and the system is not going to restart any of these containers.
	PodSucceeded PodPhase = "Succeeded"
	// PodFailed means that all containers in the pod have terminated, and at least one container has
	// terminated in a failure (exited with a non-zero exit code or was stopped by the system).
	PodFailed PodPhase = "Failed"
	// PodUnknown means that for some reason the state of the pod could not be obtained, typically due
	// to an error in communicating with the host of the pod.
	PodUnknown PodPhase = "Unknown"
)

type PodList struct {
	Items []Pod `json:"items" protobuf:"bytes,1,rep,name=items"`
}
