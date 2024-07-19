package main

import (
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	RegisterRequest struct {
		UserName string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required,min=8"`
		FullName string `json:"full_name" validate:"required"`
		Address  string `json:"address" validate:"required"`
	}

	LoginRequest struct {
		UserName string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	CustomValidator struct {
		validator *validator.Validate
	}

	jwtLoginClaims struct {
		UserName string `json:"username"`
		jwt.RegisteredClaims
	}
)

var secretKey = []byte("ct-secret-key")

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func main() {
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, time:${time_rfc3339_nano}, latency:${latency_human}\n",
	}))
	e.Validator = &CustomValidator{validator: validator.New()}

	prefixPath := e.Group("/api")

	publicGroup := prefixPath.Group("/public")

	privateGroup := prefixPath.Group("/private")

	jwtConfig := echojwt.Config{
		ContextKey: "user", // same as default
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtLoginClaims)
		},
		SigningKey: secretKey,
	}

	privateGroup.Use(echojwt.WithConfig(jwtConfig))

	publicGroup.POST("/register", register)

	publicGroup.POST("/login", login)

	privateGroup.GET("/self", self)

	e.Logger.Fatal(e.Start(":8080"))
	// http.HandleFunc("/api/public/register", register)
	// http.HandleFunc("/api/public/login", login)
	// http.HandleFunc("/api/private/self", self)

	// http.HandleFunc("/api/public/log/register", LogWrapper(register))
	// http.HandleFunc("/api/public/log/login", LogWrapper(login))
	// http.HandleFunc("/api/private/log/self", LogWrapper(self))

	// http.ListenAndServe(":8090", nil)
}

/*
		TODO #2:
		- implement the logic to register a new user (username, password, full_name, address)
	  	- Validate username (not empty and unique)
	  	- Validate password (length should at least 8)
*/
func register(c echo.Context) error {
	regReq := new(RegisterRequest)
	if err := c.Bind(regReq); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	if err := c.Validate(regReq); err != nil {
		return err
	}

	user := UserInfo{
		regReq.UserName,
		regReq.Password,
		regReq.FullName,
		regReq.Address,
	}

	if err := userStore.Save(user); err != nil {
		return c.String(http.StatusConflict, "Username already exists")
	}

	return c.JSON(http.StatusOK, user)
}

/*
		TODO #3:
		- implement the logic to login
		- validate the user's credentials (username, password)
	  	- Return JWT token to client
*/
func login(c echo.Context) error {
	loginReq := new(LoginRequest)

	if err := c.Bind(loginReq); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	if err := c.Validate(loginReq); err != nil {
		return err
	}
	user, err := userStore.Get(loginReq.UserName)

	if err != nil {
		return echo.ErrUnauthorized
	}

	if user.Password != loginReq.Password {
		return echo.ErrUnauthorized
	}

	//token, err := GenerateToken(user.UserName, 24*time.Minute)
	// Same as GenerateToken in jwt.go

	claims := &jwtLoginClaims{
		user.UserName,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(secretKey)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

/*
TODO #4:
- implement the logic to get user info
- Extract the JWT token from the header
- Validate Token
- Return user info`
*/
func self(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtLoginClaims)
	userName := claims.UserName

	userInfo, err := userStore.Get(userName)
	if err != nil {
		return c.String(http.StatusNotFound, "Not found "+userName)
	}

	return c.JSON(http.StatusOK, userInfo)
}

/*
	TODO #1: implement in-memory user store
	TODO #2: implement register handler
	TODO #3: implement login handler
	TODO #4: implement self handler

	Extra: implement log handler
*/
