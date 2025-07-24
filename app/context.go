package app

import (
	ctxPkg "context"
	"net/http"

	"go.uber.org/zap"
)

type ContextKey int

const (
	KeyUserID ContextKey = iota
)

type Context interface {
	Request() *http.Request
	Bind(v any) error
	OK(v any)
	ErrorResponse(err error)
	GetContextValue(key ContextKey) any
	Logger() *zap.Logger
	GetRequestContext() ctxPkg.Context
	Param(name string) string
	BindQuery(v any) error
}

type HandlerFunc func(Context)
