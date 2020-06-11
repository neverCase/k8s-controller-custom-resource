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
	Code   int    `json:"code" protobuf:"varint,1,opt,name=code"`
	Result string `json:"result" protobuf:"bytes,2,opt,name=result"`
}

func GetResponse(data interface{}) Response {
	r := Response{
		Code:   CodeNone,
		Result: data,
	}
	return r
}
