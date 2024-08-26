// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber
package reload

import (
	"net/http"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/google/uuid"
	"github.com/zeiss/pkg/conv"
)

var id = conv.Bytes(uuid.New().String())

// DefaultIdGenerator generates a new UUID.
func DefaultIdGenerator() []byte {
	return id
}

// Config ...
type Config struct {
	// IdGenerator
	IdGenerator func() []byte

	// Next defines a function to skip this middleware when returned true.
	Next func(c *fiber.Ctx) bool
}

// ConfigDefault is the default config.
var ConfigDefault = Config{
	IdGenerator: DefaultIdGenerator,
}

// WithHotReload is a middleware that enables a live reload of a site.
func WithHotReload(app *fiber.App, config ...Config) *fiber.App {
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Use("/static", filesystem.New(filesystem.Config{
		Root: http.FS(FS),
	}))

	app.Get("/ws/reload", Reload(config...))

	return app
}

// Reload is a middleware that enables a live reload of a site.
func Reload(config ...Config) fiber.Handler {
	cfg := configDefault(config...)

	return websocket.New(func(c *websocket.Conn) {
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				break
			}

			err = c.WriteMessage(websocket.TextMessage, cfg.IdGenerator())
			if err != nil {
				break
			}
		}
	})
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
