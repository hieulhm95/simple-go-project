package main

import (
	"context"
	"fmt"
	kafka2 "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"sideproject/config"
	"sideproject/internal/entity"
	"sideproject/internal/repository/kafka"
	"sideproject/internal/repository/mongodb"
	"strconv"
	"testing"
)

func TestKafkaProducer(t *testing.T) {

	err := kafka.CreateTopic("localhost:9093", "like_event")
	require.Nil(t, err)

	conf := config.Kafka{
		Brokers: "localhost:9093",
	}

	publisher, err := kafka.NewPublisher(kafka.PublisherOptions{
		Brokers:   []string{conf.Brokers},
		Topic:     "like_event",
		BatchSize: 1,
	})
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		err := publisher.Publish(context.TODO(), "key01", kafka.Message{
			Body: []byte("Hello " + strconv.Itoa(i)),
		})
		if err != nil {
			fmt.Println("Publish Fail ", i, " ", err.Error())
		} else {
			fmt.Println("Publish Success ", i)
		}
	}
}

func TestKafkaConsumer(t *testing.T) {
	conf := config.Kafka{
		Brokers: "localhost:9093",
	}

	consumer, err := kafka.NewConsumer(kafka.SubscribeOptions{
		Brokers: []string{
			conf.Brokers,
		},
		Topic:     "like_event",
		GroupId:   "group_notify",
		BatchSize: 1,
	})
	require.Nil(t, err)
	defer consumer.Close()

	consumer.Subscribe(context.TODO(), func(message kafka2.Message) error {
		fmt.Println("Kafka Message: ", message.Topic, " ", string(message.Value))
		return nil
	})
}

func TestMongo(t *testing.T) {
	conf := config.Mongo{
		URL: "mongodb://127.0.0.1:27017",
		DB:  "admin",
	}

	storage := mongodb.MustStorage(conf.URL, conf.DB)

	var id string
	t.Run("insert profile", func(t *testing.T) {
		res, err := storage.InsertProfile(context.TODO(), &entity.Profile{
			Name: "Hello",
		})
		require.Nil(t, err)
		id = res.Id
	})

	t.Run("get profile", func(t *testing.T) {
		res, err := storage.GetProfile(context.TODO(), id)
		require.Nil(t, err)
		require.Equal(t, res.Name, "Hello")
	})
}
