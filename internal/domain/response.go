package domain

type ResponseMessage struct {
	Message string `json:"message" example:"response message"`
}

type Response struct {
	Valid bool   `json:"valid" example:"true"`
	Error *Error `json:"error"`
}
