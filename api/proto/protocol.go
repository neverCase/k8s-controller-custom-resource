package proto

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
	Code   int32  `json:"code" protobuf:"varint,1,opt,name=code"`
	Result string `json:"result" protobuf:"bytes,2,opt,name=result"`
	//Mysql  mysqlOperatorV1.MysqlOperator `json:"mysql" protobuf:"bytes,3,opt,name=result,casttype"`
}

func GetResponse(data interface{}) Response {
	r := Response{
		Code:   CodeNone,
		Result: "",
	}
	return r
}
