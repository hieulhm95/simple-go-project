package kafka

import (
	"context"
	"fmt"
	kafka2 "github.com/segmentio/kafka-go"
	"sideproject/config"
)

func (s *SubscribeOptions) ConsumeLikeNotification(ctx context.Context) error {
	conf := config.Kafka{
		Brokers: "localhost:9093",
	}

	consumer, err := NewConsumer(SubscribeOptions{
		Brokers: []string{
			conf.Brokers,
		},
		Topic:     "like_event",
		GroupId:   "group_notify",
		BatchSize: 1,
	})
	if err != nil {
		return err
	}
	defer consumer.Close()

	consumer.Subscribe(context.TODO(), func(message kafka2.Message) error {
		fmt.Println("Kafka Message: ", message.Topic, " ", string(message.Value))
		return nil
	})
	return nil
}
