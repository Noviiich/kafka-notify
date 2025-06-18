package response

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
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
