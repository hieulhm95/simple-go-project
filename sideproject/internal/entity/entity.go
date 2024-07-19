package entity

import (
	"io"
	"time"
)

type Profile struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	User     string `json:"user"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name"`
}

type RegisterResponse struct {
	UserId   string `json:"user_id"`
	Username string `json:"username" validate:"required"`
	FullName string `json:"full_name"`
}

type UserInfo struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	HashPass string
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,gt=0"`
	Password string `json:"password" validate:"required,gt=7"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	Id       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	User     string `json:"user"`
}

type UploadImageRequest struct {
	FileName string
	Content  io.Reader
}

type UploadImageResponse struct {
	ImageUrl string `json:"imageURL"`
}

type ImageInfo struct {
	UserName  string `json:"username"`
	ImagePath string `json:"image_path"`
	FileName  string `json:"file_name"`
}

type CreatePostRequest struct {
	Caption string   `json:"caption"`
	Image   string   `json:"image"`
	Likes   []string `json:"likes"`
	User    string   `json:"user"`
	Profile string   `json:"profile"`
}

type Post struct {
	Id        string    `json:"id"`
	Caption   string    `json:"caption"`
	Image     string    `json:"image"`
	Likes     []string  `json:"likes"`
	CreatedDt time.Time `json:"createdDt"`
	UpdatedDt time.Time `json:"updatedDt"`
}

type GetPostRequest struct {
	Offset string `json:"offset"`
	Limit  string `json:"limit"`
}

type LikePostRequest struct {
	UserId string `json:"userId"`
}

type GetPostListResponse struct {
	List  []Post `json:"list"`
	Count int64  `json:"count"`
	Total int64  `json:"total"`
}
