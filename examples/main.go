package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/spf13/cobra"
	reload "github.com/zeiss/fiber-reload"
	"github.com/zeiss/pkg/server"
)

// Config ...
type Config struct {
	Flags *Flags
}

// Flags ...
type Flags struct {
	Addr string
}

// DemoRow ...
type DemoRow struct {
	ID   int
	Name string
}

var cfg = &Config{
	Flags: &Flags{},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfg.Flags.Addr, "addr", ":3000", "addr")
	rootCmd.SilenceUsage = true
}

var rootCmd = &cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(cmd.Context())
	},
}

type webSrv struct{}

func (w *webSrv) Start(ctx context.Context, ready server.ReadyFunc, run server.RunFunc) func() error {
	return func() error {
		app := fiber.New()
		app.Use(requestid.New())
		app.Use(logger.New())
		app.Use(recover.New())

		reload.WithHotReload(app)
		app.Static("/", ".")

		err := app.Listen(cfg.Flags.Addr)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	}
}

func run(ctx context.Context) error {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	webSrv := &webSrv{}

	s, _ := server.WithContext(ctx)
	s.Listen(webSrv, true)

	err := s.Wait()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
