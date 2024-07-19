// You can edit this code!
// Click here and start typing.
package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"section-04-part-2/internal/controller"
	inmemory "section-04-part-2/internal/storage/in-memory"
	"section-04-part-2/internal/usecase"
	"section-04-part-2/pkg/bucket"
	"section-04-part-2/pkg/validator"
)

func main() {
	userStore := inmemory.NewUserStore()
	imgBucket := bucket.NewFake()
	uc := usecase.NewUseCase(userStore, imgBucket)
	hdl := controller.NewHandler(uc)

	srv := newServer(hdl)
	if err := srv.Start(":8090"); err != nil {
		log.Error(err)
	}
}

func newServer(hdl *controller.Handler) *echo.Echo {
	e := echo.New()
	e.Validator = validator.NewCustomValidator()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	public := e.Group("/api/public")
	private := e.Group("/api/private")

	public.POST("/register", hdl.Register)
	public.POST("/login", hdl.Login)

	private.GET("/self", hdl.Self)
	private.POST("/self", hdl.UploadImage)

	return e
}
