package response

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Envelope struct {
	Data  any        `json:"data,omitempty"`
	Error *ErrorBody `json:"error,omitempty"`
}

func Ok(data any) Envelope {
	return Envelope{Data: data}
}

func Fail(code, message string) Envelope {
	return Envelope{Error: &ErrorBody{Code: code, Message: message}}
}
