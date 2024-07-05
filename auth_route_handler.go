package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

type authController struct {
}

func newAuthController() *authController {
	return &authController{}
}

func (a *authController) loginRouteHandler(c echo.Context) error {
	url, err := gothic.GetAuthURL(c.Response().Writer, c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, &apiResponse{Data: map[string]string{"login_url": url}})
}

func (a *authController) callbackRouteHandler(c echo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response().Writer, c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, &user)
}
