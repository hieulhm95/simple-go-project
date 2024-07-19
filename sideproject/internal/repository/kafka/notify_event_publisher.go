package kafka

import (
	"context"
	"fmt"
	"sideproject/config"
)

func (p *PublisherOptions) PublishLikeNotification(ctx context.Context, users []string) error {
	err := CreateTopic("localhost:9093", "like_event")
	if err != nil {
		return err
	}

	conf := config.Kafka{
		Brokers: "localhost:9093",
	}

	publisher, err := NewPublisher(PublisherOptions{
		Brokers:   []string{conf.Brokers},
		Topic:     "like_event",
		BatchSize: 1,
	})
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(users); i++ {
		err := publisher.Publish(context.TODO(), "key01", Message{
			Body: []byte(fmt.Sprintf("Notified to user: %s", users[i])),
		})
		if err != nil {
			fmt.Println("Publish Fail ", i, " ", err.Error())
		} else {
			fmt.Println("Publish Success ", i)
		}
	}
	return nil
}
