package proto

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
