package snake

import (
	"context"
	"reflect"
	"time"
)

type Refreshable interface {
	RefreshInterval() time.Duration
}

type MiddlewareFunc func(ctx context.Context) error

type Middleware interface {
	Wrap(MiddlewareFunc) MiddlewareFunc
}

type MiddlewareProvider interface {
	Middlewares() []Middleware
}

type middlewareResolver struct {
	mw Middleware
}

func (*middlewareResolver) IsResolver() {}

func NewMiddlewareResolver(mw Middleware) Resolver {
	return &middlewareResolver{
		mw: mw,
	}
}

func (me *middlewareResolver) Ref() Method {
	return me.mw
}

func (me *middlewareResolver) RunFunc() reflect.Value {
	// this *struct{} makes the resolver "Shared"
	return reflect.ValueOf(func(ctx context.Context) (*struct{}, error) {
		return nil, nil
	})
}
