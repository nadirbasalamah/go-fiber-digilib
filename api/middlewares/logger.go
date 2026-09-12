package middlewares

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type LoggerConfig struct {
	Config middleware.RequestLoggerConfig
}

func (c *LoggerConfig) Init() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(c.Config)
}
