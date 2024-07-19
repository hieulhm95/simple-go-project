package mongostore

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"section-05-part-1/internal/entity"
	"testing"
)

func TestImageCollection_Save(t *testing.T) {
	coll := NewImageCollection(
		"mongodb://localhost:27017",
		"sideproject",
		"image")

	id, _ := uuid.NewUUID()
	name := fmt.Sprintf("hello_%d.png", id.ID())
	info := entity.ImageInfo{
		UserName:  "user_abc",
		ImagePath: "/user_abc/" + name,
		FileName:  name,
	}

	err := coll.Save(info)
	assert.Nil(t, err)

	// re-check
	cursor, err := coll.client.Find(context.TODO(), bson.D{{"name", name}})
	require.NoError(t, err)
	requireCursorLength(t, cursor, 1)
}

func requireCursorLength(t *testing.T, cursor *mongo.Cursor, length int) {
	i := 0
	for cursor.Next(context.Background()) {
		i++
	}

	require.NoError(t, cursor.Err())
	require.Equal(t, i, length)
}
