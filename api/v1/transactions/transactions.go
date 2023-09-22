package transactions

import (
	"context"
	"net/http"
	"time"

	"github.com/jorgepiresg/ChallangePismo/app"
	modelTransactions "github.com/jorgepiresg/ChallangePismo/model/transactions"
	"github.com/jorgepiresg/ChallangePismo/utils"
	"github.com/labstack/echo/v4"
)

type handler struct {
	app app.App
}

func Register(g *echo.Group, app app.App) {
	h := handler{
		app: app,
	}

	g.POST("", h.make)
}

// get godoc
// @Summary Make transaction
// @Description make a transaction from an account.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param request body modelTransactions.MakeTransaction true "input"
// @Success      201
// @Failure      400  {object}  utils.Error
// @Router       /transactions [post]
func (h handler) make(c echo.Context) error {

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	var payload modelTransactions.MakeTransaction

	if err := c.Bind(&payload); err != nil {
		return utils.NewError(http.StatusBadRequest, "payload invalid ", nil)
	}

	err := h.app.Transactions.Make(ctx, payload)
	if err != nil {
		return utils.NewError(http.StatusBadRequest, err.Error(), nil)
	}

	c.NoContent(http.StatusCreated)
	return nil
}
