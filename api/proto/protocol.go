package proto

import (
	"github.com/nevercase/k8s-controller-custom-resource/api/group"
)

type ApiService string

const (
	SvcPing ApiService = "ping"

	SvcCreate ApiService = "create"
	SvcUpdate ApiService = "update"
	SvcDelete ApiService = "delete"
	SvcGet    ApiService = "get"
	SvcList   ApiService = "list"
	SvcWatch  ApiService = "watch"

	SvcResource ApiService = "resource"
)

const (
	CodeNone = iota
	CodeErr  = 10001
)

type Param struct {
	Service      string             `json:"service" protobuf:"bytes,1,opt,name=service"`
	ResourceType group.ResourceType `json:"resourceType" protobuf:"bytes,2,opt,name=resourceType"`
	NameSpace    string             `json:"nameSpace" protobuf:"bytes,3,opt,name=nameSpace"`
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
}

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

func GetResponse(param Param, data []byte) ([]byte, error) {
	r := Response{
		Code:   CodeNone,
		Param:  param,
		Result: data,
	}
	return r.Marshal()
}

func ErrorResponse(param Param) ([]byte, error) {
	r := Response{
		Code:  CodeErr,
		Param: param,
		//Result: "",
	}
	return r.Marshal()
}
