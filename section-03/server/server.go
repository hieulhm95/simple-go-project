package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

func main() {
	http.HandleFunc("/api/public/register", register)
	http.HandleFunc("/api/public/login", login)
	http.HandleFunc("/api/private/self", self)

	// http.HandleFunc("/api/public/log/register", LogWrapper(register))
	// http.HandleFunc("/api/public/log/login", LogWrapper(login))
	// http.HandleFunc("/api/private/log/self", LogWrapper(self))
	fmt.Println("Successfully run on localhost:8090 ")
	http.ListenAndServe(":8090", nil)
}

/*
		TODO #2:
		- implement the logic to register a new user (username, password, full_name, address)
	  	- Validate username (not empty and unique)
	  	- Validate password (length should at least 8)
*/
func register(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req RegisterRequest
	if err := json.Unmarshal(b, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(req)
	if len(req.UserName) == 0 {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 8 {
		http.Error(w, "Password not valid ! At least 8 characters", http.StatusBadRequest)
		return
	}

	userData := UserInfo{
		req.UserName,
		req.Password,
		req.FullName,
		req.Address,
	}

	if err := userStore.Save(userData); err != nil {
		return
	}

	b, err = json.Marshal(userData)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
	return
}

type RegisterRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Address  string `json:"address"`
}

/*
		TODO #3:
		- implement the logic to login
		- validate the user's credentials (username, password)
	  	- Return JWT token to client
*/
func login(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}
	var req LoginRequest
	if err := json.Unmarshal(b, &req); err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}

	user, err := userStore.Get(req.UserName)
	if err != nil {
		http.Error(w, "Username not found", http.StatusNotFound)
		return
	}

	if user.Password != req.Password {
		http.Error(w, "Password not correct", http.StatusNotFound)
		return
	}

	token, err := GenerateToken("", 24*time.Second)
	if err != nil {
		return
	}

	resp := LoginResponse{Token: token}
	res, err := json.Marshal(resp)

	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return
}

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string
}

/*
TODO #4:
- implement the logic to get user info
- Extract the JWT token from the header
- Validate Token
- Return user info`
*/
func self(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	extractUserNameFn := func(authenticationHeader string) (string, error) {
		tokenString := authenticationHeader[len("Bearer "):]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("cannot parse")
			}
			secret := []byte("ct-secret-key")
			fmt.Println(secret)
			return secret, nil
		})
		fmt.Println(token)
		if err != nil {
			return "", errors.New(fmt.Sprintf("invalid token: %v", err))
		}
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			username := claims["sub"].(string)
			return username, nil
		}
		return "", errors.New("invalid token")
	}

	username, err := extractUserNameFn(authHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := userStore.Get(username)
	if err != nil {
		// Handle the error, e.g., return an error response to client
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(user)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return
}

/*
TODO: extra wrapper
Print some logs to console
  - Path
  - Http Status code
  - Time start, Duration
*/
func LogWrapper(handler http.HandlerFunc) http.HandlerFunc {
	panic("TODO implement me")
}

/*
	TODO #1: implement in-memory user store
	TODO #2: implement register handler
	TODO #3: implement login handler
	TODO #4: implement self handler

	Extra: implement log handler
*/
