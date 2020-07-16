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
	Projects     HarborCommand = "projects"
	Repositories HarborCommand = "repositories"
	Tags         HarborCommand = "tags"
)

const (
	CodeNone = iota
	CodeErr  = 10001
)

type Param struct {
	Service      string             `json:"service" protobuf:"bytes,1,opt,name=service"`
	ResourceType group.ResourceType `json:"resourceType" protobuf:"bytes,2,opt,name=resourceType"`
	NameSpace    string             `json:"nameSpace" protobuf:"bytes,3,opt,name=nameSpace"`
	Command      HarborCommand      `json:"command" protobuf:"bytes,4,opt,name=command"`
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
	Name   string   `json:"name" protobuf:"bytes,1,rep,name=Name"`
	Master NodeSpec `json:"master" protobuf:"bytes,2,rep,name=master"`
	Slave  NodeSpec `json:"slave" protobuf:"bytes,3,rep,name=slave"`
}

type RedisCrdList struct {
	Items []RedisCrd `json:"items" protobuf:"bytes,1,rep,name=items"`
}

type RedisCrd struct {
	Name   string   `json:"name" protobuf:"bytes,1,rep,name=Name"`
	Master NodeSpec `json:"master" protobuf:"bytes,2,rep,name=master"`
	Slave  NodeSpec `json:"slave" protobuf:"bytes,3,rep,name=slave"`
}

type NodeSpec struct {
	Name             string `json:"name" protobuf:"bytes,1,rep,name=Name"`
	Replicas         int32  `json:"replicas" protobuf:"varint,2,opt,name=replicas"`
	Image            string `json:"image" protobuf:"bytes,3,rep,name=image"`
	ImagePullSecrets string `json:"imagePullSecrets" protobuf:"bytes,4,rep,name=imagePullSecrets"`
	// The path of the nas disk which was mounted on the machine
	VolumePath string `json:"volumePath" protobuf:"bytes,5,rep,name=volumePath"`
	// PodResource
	PodResource PodResourceRequirements `json:"podResource" protobuf:"bytes,6,rep,name=podResource"`
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
	// Data contains the configuration data.
	// Each key must consist of alphanumeric characters, '-', '_' or '.'.
	// Values with non-UTF-8 byte sequences must use the BinaryData field.
	// The keys stored in Data must not overlap with the keys in
	// the BinaryData field, this is enforced during validation process.
	// +optional
	Data map[string]string `json:"data" protobuf:"bytes,2,rep,name=data"`
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
	Type string `json:"type,omitempty" protobuf:"bytes,4,opt,name=type"`
	// externalIPs is a list of IP addresses for which nodes in the cluster
	// will also accept traffic for this service.  These IPs are not managed by
	// Kubernetes.  The user is responsible for ensuring that traffic arrives
	// at a node with this IP.  A common example is external load-balancers
	// that are not part of the Kubernetes system.
	// +optional
	ExternalIPs []string `json:"externalIPs,omitempty" protobuf:"bytes,5,rep,name=externalIPs"`
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

// The shortage ot the harbor projects
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
	Command   HarborCommand `json:"command" protobuf:"bytes,1,opt,name=command"`
	ProjectID int32         `json:"projectId" protobuf:"varint,2,opt,name=projectId"`
	ImageName string        `json:"imageName" protobuf:"bytes,3,opt,name=imageName"`
}
