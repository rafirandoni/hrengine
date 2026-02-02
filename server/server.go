package server

import "github.com/gofiber/fiber/v3"

type Server interface {
	StartHttpServer(app *fiber.App, address string) error
}

type server struct{}

func New() Server {
	return &server{}
}
