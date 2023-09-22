package v1

import (
	"github.com/jorgepiresg/ChallangePismo/api/v1/accounts"
	"github.com/jorgepiresg/ChallangePismo/api/v1/transactions"
	"github.com/jorgepiresg/ChallangePismo/app"
	"github.com/labstack/echo/v4"
)

type handler struct {
	app app.App
}

func Register(e *echo.Group, app app.App) {

	v1 := e.Group("/v1")

	accounts.Register(v1.Group("/accounts"), app)
	transactions.Register(v1.Group("/transactions"), app)
}
