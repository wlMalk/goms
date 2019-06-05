package errors

import "fmt"

type MethodNotImplemented struct {
	Service string
	Method  string
}

func (err *MethodNotImplemented) Error() string {
	return fmt.Sprintf("method \"%s\" is not implemented in service \"%s\"", err.Method, err.Service)
}

func ErrMethodNotImplemented(service string, method string) error {
	return &MethodNotImplemented{
		Service: service,
		Method:  method,
	}
}
