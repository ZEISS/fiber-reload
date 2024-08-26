// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber
package reload

import (
	"context"
	"net/http"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/google/uuid"
	"github.com/zeiss/pkg/conv"
)

// The contextKey type is unexported to prevent collisions with context keys defined in
// other packages.
type contextKey int

// The keys for the values in context
const (
	envCtx contextKey = iota
)

const (
	// Noop is a no-op function.
	Noop = "noop"
	// Environment environment
	Development = "development"
	// Testing environment
	Testing = "testing"
	// Staging environment
	Staging = "staging"
	// Production environment
	Production = "production"
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
func WithHotReload(app *fiber.App, config ...Config) {
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

// Environment is a middleware that sets the environment context.
func Environment(env string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := SetEnvironmentContext(c, env)
		if err != nil {
			return err
		}

		return c.Next()
	}
}

// SetEnvironmentContext sets the environment context.
func SetEnvironmentContext(c *fiber.Ctx, env string) error {
	userCtx := c.UserContext()

	envCtx := context.WithValue(userCtx, envCtx, env)
	c.SetUserContext(envCtx)

	return nil
}

// GetEnvironmentContext gets the environment context.
func GetEnvironmentContext(ctx context.Context) (string, error) {
	userCtx := ctx.Value(envCtx)
	if userCtx == nil {
		return Noop, nil
	}

	env, ok := userCtx.(string)
	if !ok {
		return Noop, nil
	}

	return env, nil
}

// IsDevelopment returns true if the environment is development.
func IsDevelopment(ctx context.Context) bool {
	env, err := GetEnvironmentContext(ctx)
	if err != nil {
		return false
	}

	return env == Development
}

// IsTesting returns true if the environment is testing.
func IsTesting(ctx context.Context) bool {
	env, err := GetEnvironmentContext(ctx)
	if err != nil {
		return false
	}

	return env == Testing
}

// IsStaging returns true if the environment is staging.
func IsStaging(ctx context.Context) bool {
	env, err := GetEnvironmentContext(ctx)
	if err != nil {
		return false
	}

	return env == Staging
}

// IsProduction returns true if the environment is production.
func IsProduction(ctx context.Context) bool {
	env, err := GetEnvironmentContext(ctx)
	if err != nil {
		return false
	}

	return env == Production
}
