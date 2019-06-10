package service

import "context"

type Method struct {
	Name    string
	Service Service
}

type Service struct {
	Name string
}

type contextKeyType string

const contextMethodKey contextKeyType = "method"

func NewMethod(service string, method string) Method {
	return Method{
		Name: method,
		Service: Service{
			Name: service,
		},
	}
}

func SetMethod(ctx context.Context, method Method) context.Context {
	return context.WithValue(ctx, contextMethodKey, method)
}

func GetMethod(ctx context.Context) Method {
	method := ctx.Value(contextMethodKey)
	if method == nil {
		return Method{}
	}
	return method.(Method)
}
