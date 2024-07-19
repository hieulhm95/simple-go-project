// You can edit this code!
// Click here and start typing.
package main

import (
	_ "github.com/labstack/echo-jwt/v4"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"section-05-part-2/config"
	"section-05-part-2/internal/controller"
	mongostore "section-05-part-2/internal/storage/mongo"
	"section-05-part-2/internal/usecase"
	auth "section-05-part-2/pkg/auth"
	"section-05-part-2/pkg/bucket"
	"section-05-part-2/pkg/validator"
)

func main() {
	conf := config.Config{
		Port:             "127.0.0.1:27017",
		MongoURI:         "mongodb://127.0.0.1:27017",
		MongoDB:          "admin",
		MongoCollImage:   "Images",
		MongoCollUser:    "Users",
		GoogleCredFile:   "/Users/hieulehoangminh/Desktop/Chotot/ct-backend-course-hieu-lehoang/section-05-part-1/pkg/bucket/gcs_creds.json",
		GoogleBucketName: "nddbao_bucket_test",
	}

	demoDB := mongostore.MustDatabase(conf.MongoURI, conf.MongoDB)
	userStore := mongostore.NewUserCollection(demoDB, conf.MongoCollUser)
	imageStore := mongostore.NewImageCollection(demoDB, conf.MongoCollImage)
	//imgBucket := bucket.MustNewGoogleStorageClient(context.TODO(), conf.GoogleBucketName, conf.GoogleCredFile)
	imgBucket := bucket.NewFake()

	uc := usecase.NewUseCase(imageStore, userStore, imgBucket)
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
	private.Use(auth.AuthMiddleware(), auth.ExtractUserNameFn)

	public.POST("/register", hdl.Register)
	public.POST("/login", hdl.Login)

	private.GET("/self", hdl.Self)
	private.POST("/upload", hdl.UploadImage)
	private.POST("/change-password", hdl.ChangePassword)
	return e
}
