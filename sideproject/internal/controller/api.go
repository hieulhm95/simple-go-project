package controller

import (
	uc "sideproject/internal/usecase"
)

type Controller struct {
	profileUC uc.ProfileUC
	userUC    uc.UserUC
	postUC    uc.PostUC
}

func NewController(profileUC uc.ProfileUC, userUC uc.UserUC, postUC uc.PostUC) *Controller {

	h := &Controller{profileUC: profileUC, userUC: userUC, postUC: postUC}
	return h
}
