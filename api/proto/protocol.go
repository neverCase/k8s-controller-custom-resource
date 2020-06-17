package proto

import (
	//_ "github.com/gogo/protobuf/gogoproto"
	//_ "github.com/gogo/protobuf/proto"
	//_ "github.com/gogo/protobuf/sortkeys"
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
	Data  string `json:"data" protobuf:"bytes,2,opt,name=data"`
}

type Response struct {
	Code   int32  `json:"code" protobuf:"varint,1,opt,name=code"`
	Param  Param  `protobuf:"bytes,2,opt,name=param"`
	Result string `json:"result" protobuf:"bytes,3,opt,name=result"`
}

type MysqlCrd struct {
	Name   string   `json:"name" protobuf:"bytes,1,rep,name=Name"`
	Master NodeSpec `json:"master" protobuf:"bytes,2,rep,name=master"`
	Slave  NodeSpec `json:"slave" protobuf:"bytes,3,rep,name=slave"`
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
}

type Resources struct {
	Resources []group.ResourceType `json:"resources" protobuf:"bytes,1,rep,name=resources"`
}

func GetResponse(param Param, data string) ([]byte, error) {
	r := Response{
		Code:   CodeNone,
		Param:  param,
		Result: data,
	}
	return r.Marshal()
}

func ErrorResponse(param Param) ([]byte, error) {
	r := Response{
		Code:   CodeErr,
		Param:  param,
		Result: "",
	}
	return r.Marshal()
}
