package middlewares

import (
	zapMiddleware "github.com/gofiber/contrib/v3/zap"
	"github.com/gofiber/fiber/v3"
)

type LoggerConfig struct {
	Config zapMiddleware.Config
}

func (c *LoggerConfig) Init() fiber.Handler {
	return zapMiddleware.New(c.Config)
}
