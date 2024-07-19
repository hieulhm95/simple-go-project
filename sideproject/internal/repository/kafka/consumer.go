package kafka

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer interface {
	Subscribe(ctx context.Context, handler HandlerFunc) error
	Close() error
	ConsumeLikeNotification(ctx context.Context) error
}

type SubscribeOptions struct {
	Brokers     []string
	Topic       string
	GroupId     string
	BatchSize   int
	StartOffset int64
	Reader      *kafka.Reader
}

func (s *SubscribeOptions) Subscribe(ctx context.Context, handler HandlerFunc) error {
	for {
		m, err := s.Reader.ReadMessage(ctx)
		if err != nil {
			return errors.New(fmt.Sprintf("subscribe %s error: %s", s.Topic, err.Error()))
		}
		handler(m)
	}
}

func (s *SubscribeOptions) Close() error {
	return s.Reader.Close()
}

type HandlerFunc func(message kafka.Message) error

type SubscribeOption func(*SubscribeOptions)

func (s *SubscribeOptions) newKafkaReader() {
	startOffset := s.StartOffset
	if startOffset == 0 {
		startOffset = kafka.FirstOffset
	}
	s.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:         s.Brokers,
		GroupID:         s.GroupId,
		Topic:           s.Topic,
		MinBytes:        10,   // 10B
		MaxBytes:        10e6, // 10MB
		ReadLagInterval: -1 * time.Second,
		MaxWait:         100 * time.Millisecond,
		RetentionTime:   7 * 24 * time.Hour,
		StartOffset:     startOffset,
	})
}

func (s *SubscribeOptions) Validate() error {
	if len(s.Brokers) == 0 {
		return errors.New("cannot create a new subscriber with an empty list of broker addresses")
	}

	if len(s.Topic) == 0 {
		return errors.New("cannot create a new subscriber with an empty topic")
	}

	if len(s.GroupId) == 0 {
		return errors.New("cannot create a new subscriber with an empty group")
	}

	return nil
}

func NewConsumer(opts SubscribeOptions) (Consumer, error) {
	s := opts
	err := s.Validate()
	if err != nil {
		return nil, err
	}
	s.newKafkaReader()
	return &s, nil
}
