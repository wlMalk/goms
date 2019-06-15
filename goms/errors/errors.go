package errors

import "fmt"

type ErrMethodNotImplemented struct {
	Service string
	Method  string
}

func (err *ErrMethodNotImplemented) Error() string {
	return fmt.Sprintf("method '%s' is not implemented in service '%s'", err.Method, err.Service)
}

type ErrInvalidRequest struct {
	Service string
	Method  string
}

func (err *ErrInvalidRequest) Error() string {
	return fmt.Sprintf("invalid request received for '%s' method in '%s' service", err.Method, err.Service)
}

type ErrInvalidResponse struct {
	Service string
	Method  string
}

func (err *ErrInvalidResponse) Error() string {
	return fmt.Sprintf("invalid response returned from '%s' method in '%s' service", err.Method, err.Service)
}

func MethodNotImplemented(service string, method string) error {
	return &ErrMethodNotImplemented{
		Service: service,
		Method:  method,
	}
}

func InvalidRequest(service string, method string) error {
	return &ErrInvalidRequest{
		Service: service,
		Method:  method,
	}
}

func InvalidResponse(service string, method string) error {
	return &ErrInvalidResponse{
		Service: service,
		Method:  method,
	}
}
