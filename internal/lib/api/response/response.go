package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Alias  string `json:"alias,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK(alias string) Response {
	return Response{
		Status: StatusOK,
		Alias:  alias,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a URL", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}
	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, "; "),
	}
}
