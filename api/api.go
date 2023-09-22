package api

import (
	"log"

	v1 "github.com/jorgepiresg/ChallangePismo/api/v1"
	"github.com/jorgepiresg/ChallangePismo/app"
	"github.com/labstack/echo/v4"
)

type Options struct {
	Group *echo.Group
	App   app.App
}

func New(opts Options) {

	v1.Register(opts.Group, opts.App)

	log.Println("API Created")
}
