package accounts

import (
	"context"
	"net/http"
	"time"

	"github.com/jorgepiresg/ChallangePismo/app"
	modelAccounts "github.com/jorgepiresg/ChallangePismo/model/accounts"
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

	g.POST("", h.create)
	g.GET("/:account_id", h.getByAccountID)
}

// create godoc
// @Summary Account create
// @Description create a account
// @Tags         Account
// @Accept       json
// @Produce      json
// @Param request body modelAccounts.Create true "input"
// @Success      201  {object}  modelAccounts.CreateResponse
// @Failure      400  {object}  utils.Error
// @Router       /accounts [post]
func (h handler) create(c echo.Context) error {

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	var payload modelAccounts.Create

	if err := c.Bind(&payload); err != nil {
		return utils.NewError(http.StatusBadRequest, "payload invalid ", err.Error())
	}

	account, err := h.app.Accounts.Create(ctx, payload)
	if err != nil {
		return utils.NewError(http.StatusBadRequest, "fail to create a new account", err.Error())
	}

	res := modelAccounts.CreateResponse{
		AccountID: account.ID,
	}

	c.JSON(http.StatusCreated, res)

	return nil
}

// getByAccountID godoc
// @Summary Account
// @Description get account by id
// @Tags         Account
// @Accept       json
// @Produce      json
// @Param        account_id   path      string  true  "Account ID"
// @Success      200  {object}  modelAccounts.Account
// @Failure      400  {object}  utils.Error
// @Router       /accounts/{account_id} [get]
func (h handler) getByAccountID(c echo.Context) error {

	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	accountID := c.Param("account_id")

	res, err := h.app.Accounts.GetByAccountID(ctx, accountID)
	if err != nil {
		return utils.NewError(http.StatusBadRequest, err.Error(), nil)
	}

	c.JSON(http.StatusOK, res)

	return nil
}
