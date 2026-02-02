package web

import (
	"github.com/gofiber/fiber/v3"
)

type Web interface {
	RSuccess(ctx fiber.Ctx, data interface{}, options ...*ROption) error
	RFail(ctx fiber.Ctx, err error) error
}

type web struct{}

func New() Web {
	return &web{}
}
