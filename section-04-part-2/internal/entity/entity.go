package entity

import "io"

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name"`
	Address  string `json:"address"`
}

type RegisterResponse struct {
	UserId string `json:"user_id"`
}

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type SelfRequest struct {
	Token string
}

type SelfResponse struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Address  string `json:"address"`
}

type UploadImageRequest struct {
	File     io.Reader `json:"file"`
	FileName string    `json:"file_name"`
}

type UploadImageResponse struct {
	Url string `json:"url"`
}

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Address  string `json:"address"`
}

type ImageInfo struct {
	// TODO
}
