package proto

import (
	_ "github.com/gogo/protobuf/gogoproto"
	_ "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/sortkeys"

	mysqlOperatorV1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/mysqloperator/v1"
)

const (
	SvcPing     = "ping"
	SvcList     = "list"
	SvcWatch    = "watch"
	SvcAdd      = "add"
	SvcUpdate   = "update"
	SvcDelete   = "delete"
	SvcResource = ""
)

const (
	CodeNone = iota
	CodeErr
)

type Request struct {
	Service string `json:"service" protobuf:"bytes,1,opt,name=service"`
	Data    string `json:"data" protobuf:"bytes,2,opt,name=data"`
}

type Response struct {
	Code   int32  `json:"code" protobuf:"varint,1,opt,name=code"`
	Result string `json:"result" protobuf:"bytes,2,opt,name=result"`
}

type List struct {
	Code   int32  `json:"code" protobuf:"varint,1,rep,name=code"`
	Result string `json:"result" protobuf:"bytes,2,rep,name=result"`
	//Mysql  mysqlOperatorV1.MysqlOperator `json:"mysql" protobuf:"bytes,3,opt,name=result,casttype"`
}

type Mysql struct {
	Mysql mysqlOperatorV1.MysqlOperator `json:"mysql" protobuf:"varint,1,opt,name=code,casttype=github.com/nevercase/k8s-controller-custom-resource/pkg/apis/mysqloperator/v1.MysqlOperator"`
}

func GetResponse(data string) Response {
	r := Response{
		Code:   CodeNone,
		Result: data,
	}
	return r
}
