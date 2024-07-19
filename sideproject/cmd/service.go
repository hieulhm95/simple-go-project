package cmd

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"net/http"
	"sideproject/config"
	"sideproject/internal/controller"
	"sideproject/internal/repository/kafka"
	"sideproject/internal/repository/mongodb"
	"sideproject/internal/usecase"
	"sideproject/pkg/auth"
	"sideproject/pkg/bucket"
	"sideproject/pkg/validator"
)

var service = &cobra.Command{
	Use:   "service",
	Short: "API Command of service",
	Long:  "API Command of service",
	Run: func(cmd *cobra.Command, args []string) {
		RunHTTPServer()
	},
}

func RunHTTPServer() {
	conf := config.MustLoad()

	fmt.Println(conf)

	e := SetupEcho(conf)
	e.Validator = validator.NewCustomValidator()

	mongoStore := mongodb.MustStorage(conf.Mongo.URL, conf.Mongo.DB)
	gcsStore := bucket.MustNewGoogleStorageClient(context.TODO(), conf.Gcs.GoogleBucketName, conf.Gcs.GoogleCredFile, conf.Gcs.GoogleHost)

	if err := kafka.CreateTopic(conf.Kafka.Brokers, conf.Kafka.Topic); err != nil {
		panic(err)
	}

	publisher := kafka.MustPublisher(kafka.PublisherOptions{
		Brokers:   []string{conf.Kafka.Brokers},
		Topic:     conf.Kafka.Topic,
		BatchSize: 1,
	})
	profileUC := usecase.NewProfileUC(mongoStore, publisher)
	userUC := usecase.NewUserUC(mongoStore, publisher)
	postUC := usecase.NewPostUC(mongoStore, publisher, gcsStore, profileUC)
	api := controller.NewController(profileUC, userUC, postUC)

	internal := e.Group("/api/v1/internal")
	public := e.Group("/api/v1/public")
	internal.Use(auth.AuthMiddleware(), auth.ExtractUserNameFn)
	internal.GET("/profile/:profileId", api.GetProfile)
	internal.POST("/profile", api.CreateProfile)
	public.POST("/register", api.Register)
	public.POST("/login", api.Login)
	internal.POST("/post", api.CreatePost)
	internal.GET("/post", api.GetPostList)
	internal.GET("/post/:postId", api.GetPostById)
	internal.POST("/upload", api.UploadImage)
	internal.POST("/like/:postId", api.LikePost)
	e.Logger.Fatal(e.Start(":" + conf.AppPort))
}

func SetupEcho(conf *config.Config) *echo.Echo {
	e := echo.New()
	// Debug mode
	if conf.ENV != "production" {
		e.Debug = true
	}

	// add health check for echo
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	return e
}
