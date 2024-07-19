package mongodb

import (
	"github.com/google/uuid"
	"sideproject/internal/entity"
	"sideproject/pkg/hashpass"
	"time"
)

var (
	_ IDocument = &ProfileDoc{}
)

type IDocument interface {
	GetDocId() string
}

func NewDoc() Doc {
	docId := uuid.New().String()
	return Doc{
		DocId:     docId,
		Version:   1,
		CreatedDt: time.Now(),
		UpdatedDt: time.Now(),
	}
}

type Doc struct {
	DocId     string    `bson:"docId"`
	Version   int64     `bson:"version"`
	CreatedDt time.Time `bson:"createdDt"`
	UpdatedDt time.Time `bson:"updatedDt"`
}

func (d *Doc) GetDocId() string {
	return d.DocId
}

type ProfileDoc struct {
	Doc      `bson:",inline"`
	Id       string `bson:"id"`
	Name     string `bson:"name"`
	User     string `bson:"user"`
	Username string `bson:"username"`
}

type UserDoc struct {
	Doc      `bson:",inline"`
	Id       string `bson:"id"`
	Username string `bson:"username"`
	HashPass string `bson:"hashPass"`
	FullName string `bson:"fullName"`
}

type PostDoc struct {
	Doc     `bson:",inline"`
	Id      string   `bson:"id"`
	Caption string   `bson:"caption"`
	Image   string   `bson:"image"`
	Likes   []string `bson:"likes"`
	User    string   `bson:"user"`
	Profile string   `bson:"profile"`
}

func NewProfileDoc(profile *entity.Profile) ProfileDoc {
	return ProfileDoc{
		Doc:      NewDoc(),
		Id:       uuid.NewString(),
		Name:     profile.Name,
		User:     profile.User,
		Username: profile.Username,
	}
}

func NewUserDocument(info *entity.RegisterRequest) *UserDoc {
	return &UserDoc{
		Doc:      NewDoc(),
		Id:       uuid.NewString(),
		Username: info.Username,
		HashPass: hashpass.HashPassword(info.Password),
		FullName: info.FullName,
	}
}

func NewPostDoc(post *entity.CreatePostRequest) *PostDoc {
	return &PostDoc{
		Doc:     NewDoc(),
		Id:      uuid.NewString(),
		Caption: post.Caption,
		Image:   post.Image,
		Likes:   post.Likes,
		User:    post.User,
		Profile: post.Profile,
	}
}
