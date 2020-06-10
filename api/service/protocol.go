package service

const (
	SvcPing   = "ping"
	SvcList   = "list"
	SvcWatch  = "watch"
	SvcAdd    = "add"
	SvcUpdate = "update"
	SvcDelete = "delete"
	SvcResource = ""
)

const (
	CodeNone = iota
	CodeErr
)

type Request struct {
	Service string      `json:"service"`
	Data    interface{} `json:"data"`
}

type Response struct {
	Code   int         `json:"code"`
	Result interface{} `json:"result"`
}

func GetResponse(data interface{}) Response {
	r := Response{
		Code:   CodeNone,
		Result: data,
	}
	return r
}
