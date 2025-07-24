package app

import (
	"encoding/json"
	"net/http"
)

type response struct {
	ApplicationId string `json:"applicationId,omitempty"`
	Timestamp     string `json:"timestamp,omitempty"`
	Eligible      *bool  `json:"eligible,omitempty"`
	Message       string `json:"message"`
	Reason        string `json:"reason,omitempty"`
}

type Response struct {
	HTTPStatusCode int
	Response       response
}

func (err *Response) Error() string {
	return err.Response.Message
}

func NewErrorResponse(httpCode int, code string, msg string, eligible *bool, reason ...string) *Response {
	if len(reason) > 0 {
		return &Response{
			HTTPStatusCode: httpCode,
			Response: response{
				Eligible: eligible,
				Message:  msg,
				Reason:   reason[0],
			},
		}
	}
	return &Response{
		HTTPStatusCode: httpCode,
		Response: response{
			Eligible: eligible,
			Message:  msg,
		},
	}
}

func Success(data any) map[string]any {
	// convert struct -> map[string]any
	var dataMap map[string]any
	tmp, _ := json.Marshal(data)
	_ = json.Unmarshal(tmp, &dataMap)

	result := map[string]any{
		"code":    CodeSuccess,
		"message": MsgSuccess,
	}

	// merge data
	for k, v := range dataMap {
		result[k] = v
	}

	return result
}

func InvalidRequestBody(reason string) *Response {
	return NewErrorResponse(http.StatusBadRequest, CodeInvalidRequestBody, MsgInvalidRequestBody, nil, reason)
}

func InvalidRequestParam(reason string) *Response {
	return NewErrorResponse(http.StatusBadRequest, CodeInvalidRequestParam, MsgInvalidRequestParam, nil, reason)
}

func LoanApplicationNotFound(reason string) *Response {
	return NewErrorResponse(http.StatusNotFound, CodeLoanApplicationNotFound, MsgLoanApplicationNotFound, nil, reason)
}

func UnexpectedRequest() *Response {
	return NewErrorResponse(http.StatusBadRequest, CodeUnexpectedRequest, MsgUnexpectedRequest, nil)
}

func InternalServerError(reason ...string) *Response {
	return NewErrorResponse(http.StatusInternalServerError, CodeInternalServerError, MsgInternalServerError, nil, reason...)
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
	MsgSuccess                     = "Success"
	MsgInvalidRequestBody          = "Invalid request body"
	MsgInvalidRequestParam         = "Invalid request param"
	MsgLoanApplicationNotFound     = "Loan application not found"
	MsgMonthlyIncomeIsInsufficient = "Monthly income is insufficient"
	MsgAgeNotInRange               = "Age not in range (must be between 20-60)"
	MsgBusinessLoansNotSupported   = "Business loans not supported"
	MsgLoanAmountExceedsLimit      = "Loan amount cannot exceed 12 months of income"
	MsgLoanPurposeInvalid          = "Loan purpose must be one of: education, home, car, business, personal"
	MsgUnexpectedRequest           = "Unexpected Request"
	MsgInternalServerError         = "Internal Server Error"
)

// helper
func BoolPtr(b bool) *bool {
	return &b
}
