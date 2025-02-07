package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tuannamnguyen/go-chat-app/internal/models"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

type AuthService struct {
	config *oauth2.Config
	db     AuthRepository
}

type AuthRepository interface {
	GetUserInfo(ctx context.Context, userID string) string
	SetUserInfo(ctx context.Context, userID string, userName string) error
}

func NewAuthService(config *oauth2.Config, db AuthRepository) *AuthService {
	return &AuthService{
		config: config,
		db:     db,
	}
}

func (a *AuthService) LoginHandler(c echo.Context) error {
	url := a.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	return c.JSON(http.StatusOK, &models.ApiResponse{
		Data: map[string]any{
			"auth_url": url,
		},
	})
}

func getUserInformation(ctx context.Context, client *http.Client) (*people.Person, error) {
	srv, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("error creating people api service: %v", err)
	}

	res, err := srv.People.Get("people/me").PersonFields("names").Do()
	if err != nil {
		return nil, fmt.Errorf("error getting call api get user info: %v", err)
	}

	return res, nil
}

func (a *AuthService) CallbackHandler(c echo.Context) error {
	requestCtx := c.Request().Context()

	authCode := c.QueryParam("code")
	token, err := a.config.Exchange(requestCtx, authCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error when exchange auth code: %v", err))
	}

	apiClient := a.config.Client(requestCtx, token)
	userInfo, err := getUserInformation(requestCtx, apiClient)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error when get user information: %v", err))
	}

	peopleID := strings.Split(userInfo.ResourceName, "/")[1]
	err = a.db.SetUserInfo(requestCtx, peopleID, userInfo.Names[0].DisplayName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error when saving user info: %v", err))
	}

	return c.JSON(http.StatusOK, models.ApiResponse{Data: map[string]any{"user_id": peopleID}})
}

func (a *AuthService) GetUserName(c echo.Context) error {
	requestCtx := c.Request().Context()
	userID := c.Param("user_id")

	userName := a.db.GetUserInfo(requestCtx, userID)

	return c.JSON(http.StatusOK, models.ApiResponse{Data: map[string]any{"user_name": userName}})
}
