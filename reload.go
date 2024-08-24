// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber
package reload

import (
	"github.com/gofiber/fiber/v2"
)

// Config ...
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	Next func(c *fiber.Ctx) bool
}

// ConfigDefault is the default config.
var ConfigDefault = Config{}

// Reload is a middleware that enables a live reload of a site.
func Authenticate(handler fiber.Handler, config ...Config) fiber.Handler {
	cfg := configDefault(config...)

	return func(c *fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		return nil
	}
}

// Helper function to set default values
func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	return cfg
}
