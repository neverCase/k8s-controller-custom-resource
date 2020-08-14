package proto

func GetResponse(param Param, data []byte) ([]byte, error) {
	r := Response{
		Code:   CodeNone,
		Param:  param,
		Result: data,
	}
	return r.Marshal()
}

func ErrorResponse(param Param, errMessage string) ([]byte, error) {
	r := Response{
		Code:   CodeErr,
		Param:  param,
		Result: []byte(errMessage),
	}
	return r.Marshal()
}
