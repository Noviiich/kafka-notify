package response

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

const (
	MsgOK = "OK"
)

func Error(msg string) Response {
	return Response{
		Success: false,
		Message: msg,
	}
}

func OK() Response {
	return Response{
		Success: true,
		Message: MsgOK,
	}
}

func Success(data interface{}) Response {
	return Response{
		Success: true,
		Message: MsgOK,
		Data:    data,
	}
}
