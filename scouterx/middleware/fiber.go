package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zbum/scouter-agent-golang/scouterx/strace"
)

func FiberTracingMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		newCtx := strace.StartFastHttpService(ctx, c.Context())
		c.SetUserContext(newCtx)
		defer strace.EndFastHttpService(newCtx, c.Context())
		return c.Next()
	}
}
