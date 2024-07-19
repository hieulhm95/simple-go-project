package controller

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"sideproject/internal/entity"
	"sideproject/pkg/auth"
	"time"
)

func (c *Controller) Register(ctx echo.Context) error {
	var req entity.RegisterRequest
	if err := ctx.Bind(&req); err != nil {
		return fmt.Errorf("bind: %w", err)
	}

	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	resp, err := c.userUC.CreateUser(context.TODO(), &req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) Login(ctx echo.Context) error {
	var req entity.LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return fmt.Errorf("bind: %w", err)
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	resp, err := c.userUC.GetUser(context.TODO(), &entity.UserInfo{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	profile, err := c.createProfile(context.TODO(), resp.Id, resp.Username, resp.FullName)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	token, err := auth.GenerateToken(resp.Id, 60*time.Minute)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusOK, &entity.LoginResponse{
		Token:    token,
		Id:       profile.Id,
		Name:     profile.Name,
		Username: profile.Username,
		User:     profile.User,
	})
}

func (c *Controller) createProfile(ctx context.Context, userId string, username string, name string) (*entity.Profile, error) {
	profile, err := c.profileUC.GetProfileByUserId(ctx, userId)
	if err != nil {
		resp, profileErr := c.profileUC.InsertProfile(ctx, &entity.Profile{
			Name:     name,
			Username: username,
			User:     userId,
		})
		if profileErr != nil {
			return nil, profileErr
		}
		return resp, nil
	}
	return profile, nil
}

func (c *Controller) ChangePassword(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "Change Password")
}
