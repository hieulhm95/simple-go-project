package controller

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"sideproject/internal/entity"
)

func (c *Controller) CreateProfile(ctx echo.Context) error {
	var req entity.Profile
	if err := ctx.Bind(&req); err != nil {
		return fmt.Errorf("bind: %w", err)
	}

	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	resp, err := c.profileUC.InsertProfile(context.TODO(), &req)
	if err != nil {
		return ctx.JSON(http.StatusOK, err.Error())
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) GetProfile(ctx echo.Context) error {
	profileId := ctx.Param("profileId")

	resp, err := c.profileUC.GetProfileByUserId(context.TODO(), profileId)
	if err != nil {
		return ctx.JSON(http.StatusOK, err.Error())
	}
	return ctx.JSON(http.StatusOK, resp)
}
