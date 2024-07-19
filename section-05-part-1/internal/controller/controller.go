package controller

import (
	"context"
	"fmt"
	echo "github.com/labstack/echo/v4"
	"net/http"
	"section-05-part-1/internal/entity"
)

type UseCase interface {
	Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error)
	Register(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResponse, error)
	Self(ctx context.Context, req *entity.SelfRequest) (*entity.SelfResponse, error)
	UploadImage(ctx context.Context, req *entity.UploadImageRequest, user *entity.SelfRequest) (*entity.UploadImageResponse, error)
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

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resp, err := h.uc.Login(context.TODO(), &req)
	if err != nil {
		return fmt.Errorf("uc.Login: %w", err)
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) Self(c echo.Context) error {
	selfReq := entity.SelfRequest{Username: c.Get("username").(string)}
	resp, err := h.uc.Self(context.TODO(), &selfReq)
	if err != nil {
		fmt.Print("error self controller")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) UploadImage(c echo.Context) error {
	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile("image")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	selfReq := entity.SelfRequest{Username: c.Get("username").(string)}
	var req = entity.UploadImageRequest{FileName: file.Filename, Content: src}
	resp, err := h.uc.UploadImage(context.TODO(), &req, &selfReq)
	if err != nil {
		return fmt.Errorf("uc.UploadImage: %w", err)
	}

	return c.JSON(http.StatusOK, resp)
}
