package app

import (
	ctxPkg "context"
	"errors"
	"net/http"

	"github.com/go-playground/form"
	"github.com/inkxk/bad-bot/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type context struct {
	echo.Context
	logger *zap.Logger
}

func NewContext(c echo.Context, logger *zap.Logger) Context {
	return &context{Context: c, logger: logger}
}

func (c *context) Request() *http.Request {
	return c.Context.Request()
}

func (c *context) Bind(v any) error {
	return c.Context.Bind(v)
}

func (c *context) OK(v any) {
	err := c.Context.JSON(http.StatusOK, v)
	if err != nil {
		c.logger.Error("error while sending response", zap.Error(err))
	}
}

func (c *context) Param(name string) string {
	return c.Context.Param(name)
}

func (c *context) BindQuery(v any) error {
	values := c.QueryParams()
	return form.NewDecoder().Decode(v, values)
}

func (c *context) ErrorResponse(err error) {
	var appResponse *Response
	if errors.As(err, &appResponse) {
		err := c.Context.JSON(appResponse.HTTPStatusCode, appResponse.Response)
		if err != nil {
			c.logger.Error("error while sending response", zap.Error(err))
		}
		return
	}
	err = c.Context.JSON(http.StatusInternalServerError, response{
		Message: MsgInternalServerError,
	})
	if err != nil {
		c.logger.Error("error while sending response", zap.Error(err))
	}
}

func (c *context) GetContextValue(key ContextKey) any {
	return c.Request().Context().Value(key)
}

func (c *context) GetRequestContext() ctxPkg.Context {
	return c.Request().Context()
}

func (c *context) Logger() *zap.Logger {
	return c.logger
}

func NewEchoHandler(handler HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log, err := logger.FromContext(ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response{
				Message: MsgInternalServerError,
			})
		}

		handler(NewContext(c, log))
		return nil
	}
}

type Router struct {
	*echo.Echo
}

func NewRouter() *Router {
	r := echo.New()

	return &Router{Echo: r}
}

func (r *Router) GET(path string, handler func(Context), middleware ...echo.MiddlewareFunc) {
	r.Echo.GET(path, NewEchoHandler(handler), middleware...)
}

func (r *Router) POST(path string, handler func(Context), middleware ...echo.MiddlewareFunc) {
	r.Echo.POST(path, NewEchoHandler(handler), middleware...)
}

func (r *Router) Use(middleware ...echo.MiddlewareFunc) {
	r.Echo.Use(middleware...)
}

func (r *Router) HealthCheck(middleware ...echo.MiddlewareFunc) {
	r.GET("/health", func(c Context) {
		c.OK(map[string]string{"status": "ok"})
	})
}
