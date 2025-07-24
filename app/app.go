package app

import (
	"net/http"
)

type response struct {
	Message string `json:"message"`
}

type Response struct {
	HTTPStatusCode int
	Response       response
}

func (err *Response) Error() string {
	return err.Response.Message
}

func NewErrorResponse(httpCode int, code string, msg string, reason ...string) *Response {
	if len(reason) > 0 {
		return &Response{
			HTTPStatusCode: httpCode,
			Response: response{
				Message: msg,
			},
		}
	}
	return &Response{
		HTTPStatusCode: httpCode,
		Response: response{
			Message: msg,
		},
	}
}

func UnexpectedRequest() *Response {
	return NewErrorResponse(http.StatusBadRequest, CodeUnexpectedRequest, MsgUnexpectedRequest)
}

func InternalServerError(reason ...string) *Response {
	return NewErrorResponse(http.StatusInternalServerError, CodeInternalServerError, MsgInternalServerError, reason...)
}

// response code
const (
	CodeSuccess = "0000"

	CodeInvalidRequestBody      = "4000"
	CodeInvalidRequestParam     = "4001"
	CodeLoanApplicationNotFound = "4002"
	CodeUnexpectedRequest       = "4999"

	CodeInternalServerError = "5999"
)

// response msg
const (
	MsgSuccess = "Success"

	MsgUnexpectedRequest   = "Unexpected Request"
	MsgInternalServerError = "Internal Server Error"
)
