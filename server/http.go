package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
)

func (slf *server) StartHttpServer(app *fiber.App, address string) error {
	chListen := make(chan bool)
	go func() {
		defer close(chListen)

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		<-ch
		app.Shutdown()
	}()

	app.Use(recoverer.New())

	err := app.Listen(address)
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	<-chListen

	log.Println("engine turned off...")

	return nil
}
