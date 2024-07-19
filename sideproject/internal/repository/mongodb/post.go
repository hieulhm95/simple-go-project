package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sideproject/internal/entity"
)

func (s *Storage) GetPostById(ctx context.Context, id string) (*entity.Post, error) {
	var doc PostDoc
	ctx, cancelFn := context.WithTimeout(context.Background(), s.timeout)
	defer cancelFn()

	filter := bson.M{"id": id}
	if err := s.db.Collection(CollectionPost).FindOne(ctx, filter).Decode(&doc); err != nil {
		return nil, err
	}

	return &entity.Post{
		Id:      doc.Id,
		Caption: doc.Caption,
		Image:   doc.Image,
		Likes:   doc.Likes,
	}, nil
}

func (s *Storage) count(ctx context.Context) (int64, error) {
	countOption := options.Count().SetHint("_id_")
	filter := bson.M{}

	count, err := s.db.Collection(CollectionPost).CountDocuments(ctx, filter, countOption)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Storage) GetPostList(ctx context.Context, id string, offset int64, limit int64) (*entity.GetPostListResponse, error) {
	var doc PostDoc
	var results []entity.Post
	ctx, cancelFn := context.WithTimeout(context.Background(), s.timeout)
	defer cancelFn()
	count, err := s.count(ctx)
	if err != nil {
		return nil, err
	}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "createdDt", Value: -1}})
	findOptions.SetLimit(limit)
	findOptions.SetSkip(offset)
	filter := bson.M{}
	cur, err := s.db.Collection(CollectionPost).Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		err := cur.Decode(&doc)
		if err != nil {
			return nil, err
		}

		results = append(results, entity.Post{
			Id:        doc.Id,
			Caption:   doc.Caption,
			Image:     doc.Image,
			Likes:     doc.Likes,
			CreatedDt: doc.CreatedDt,
			UpdatedDt: doc.UpdatedDt,
		})
	}
	total := len(results)
	return &entity.GetPostListResponse{
		List:  results,
		Total: int64(total),
		Count: count,
	}, nil
}

func (s *Storage) InsertPost(ctx context.Context, req *entity.CreatePostRequest) (*entity.Post, error) {
	postDoc := NewPostDoc(req)
	ctx, cancelFn := context.WithTimeout(context.Background(), s.timeout)
	defer cancelFn()

	_, err := s.db.Collection(CollectionPost).InsertOne(ctx, postDoc)
	if err != nil {
		return nil, err
	}
	return &entity.Post{
		Id:      postDoc.Id,
		Caption: postDoc.Caption,
		Image:   postDoc.Image,
	}, nil
}

func (s *Storage) UploadPostImage(ctx context.Context, req *entity.UploadImageRequest) (*entity.UploadImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) LikePost(ctx context.Context, req *entity.LikePostRequest, post *entity.Post) (*entity.Post, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), s.timeout)
	defer cancelFn()

	filter := bson.M{"id": post.Id}
	update := bson.M{"$set": bson.M{"likes": post.Likes}}
	_, err := s.db.Collection(CollectionPost).UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return &entity.Post{
		Id:      post.Id,
		Caption: post.Caption,
		Image:   post.Image,
		Likes:   post.Likes,
	}, nil
}
