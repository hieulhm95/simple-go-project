package controller

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"section-04-part-2/internal/entity"
	"section-04-part-2/pkg/auth"
)

type UseCase interface {
	Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error)
	Register(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResponse, error)
	Self(ctx context.Context, userName string) (*entity.SelfResponse, error)
	UploadImage(ctx context.Context, req *entity.UploadImageRequest) (*entity.UploadImageResponse, error)

	// TODO: implement more
}

func NewHandler(uc UseCase) *Handler {
	return &Handler{uc: uc}
}

type Handler struct {
	uc UseCase
}

func (h *Handler) Register(c echo.Context) error {
	var req entity.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return fmt.Errorf("bind: %w", err)
	}

	resp, err := h.uc.Register(context.TODO(), &req)
	if err != nil {
		return fmt.Errorf("uc.Register: %w", err)
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) Login(c echo.Context) error {
	var req entity.LoginRequest
	if err := c.Bind(&req); err != nil {
		return fmt.Errorf("bind: %w", err)
	}

	resp, err := h.uc.Login(context.TODO(), &req)
	if err != nil {
		return fmt.Errorf("uc.Login: %w", err)
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) Self(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	username, err := auth.ParseToken(authHeader)
	if err != nil {
		return fmt.Errorf("uc.Self: %w", err)
	}
	user, err := h.uc.Self(context.TODO(), username)
	if err != nil {
		return fmt.Errorf("uc.Self: %w", err)
	}
	return c.JSON(http.StatusOK, user)

}

func (h *Handler) UploadImage(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return fmt.Errorf("uc.UploadImage", err)
	}
	source, err := file.Open()
	if err != nil {
		return fmt.Errorf("uc.UploadImage ", err)
	}
	defer source.Close()
	resp, err := h.uc.UploadImage(context.TODO(), &entity.UploadImageRequest{
		File:     source,
		FileName: file.Filename,
	})
	if err != nil {
		return fmt.Errorf("uc.UploadImage ", err)
	}
	return c.JSON(http.StatusOK, resp)
}
