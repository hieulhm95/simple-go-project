package controller

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"sideproject/internal/entity"
	"strconv"
)

func (c *Controller) CreatePost(ctx echo.Context) error {
	var req entity.CreatePostRequest
	if err := ctx.Bind(&req); err != nil {
		return fmt.Errorf("bind: %w", err)
	}

	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	resp, err := c.postUC.InsertPost(context.TODO(), &req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) GetPostById(ctx echo.Context) error {
	postId := ctx.Param("postId")
	resp, err := c.postUC.GetPostById(context.TODO(), postId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) GetPostList(ctx echo.Context) error {
	selfId := ctx.Get("userId").(string)
	limitReq, _ := strconv.Atoi(ctx.QueryParam("limit"))
	offsetReq, _ := strconv.Atoi(ctx.QueryParam("offset"))
	resp, err := c.postUC.GetPostList(context.TODO(), selfId, int64(offsetReq), int64(limitReq))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusOK, resp)

}

func (c *Controller) LikePost(ctx echo.Context) error {
	postId := ctx.Param("postId")
	var req entity.LikePostRequest
	if err := ctx.Bind(&req); err != nil {
		return fmt.Errorf("bind: %w", err)
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	post, err := c.postUC.GetPostById(ctx.Request().Context(), postId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	resp, err := c.postUC.LikePost(context.TODO(), &req, post)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) UploadImage(ctx echo.Context) error {
	file, err := ctx.FormFile("image")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	var req = entity.UploadImageRequest{FileName: file.Filename, Content: src}
	resp, err := c.postUC.UploadPostImage(context.TODO(), &req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, resp)
}
