package kafka

import (
	"context"
	"errors"
	"github.com/labstack/gommon/log"
	"github.com/segmentio/kafka-go"
	"net"
	"strconv"
	"time"
)

type Publisher interface {
	Publish(ctx context.Context, key string, message Message) error
	PublishMultiRaw(ctx context.Context, key string, message []Message) error
	PublishLikeNotification(ctx context.Context, users []string) error
}

type Message struct {
	Header map[string]string
	Body   []byte
}

type PublisherOptions struct {
	Brokers   []string
	Topic     string
	BatchSize int
	writer    *kafka.Writer
}

func (p *PublisherOptions) Validate() error {
	if len(p.Brokers) == 0 {
		return errors.New("cannot create a new publisher with an empty list of broker addresses")
	}

	if len(p.Topic) == 0 {
		return errors.New("cannot create a new publisher with an empty topic")
	}

	return nil
}

func MustPublisher(configs PublisherOptions) Publisher {
	pub, err := NewPublisher(configs)
	if err != nil {
		panic(err)
	}

	log.Info("init kafka publisher")
	return pub
}

func NewPublisher(configs PublisherOptions) (Publisher, error) {

	p := configs
	if p.BatchSize <= 0 {
		p.BatchSize = 1
	}
	err := p.Validate()
	if err != nil {
		return nil, err
	}
	p.newKafkaWriter()
	return &p, nil
}

func (p *PublisherOptions) Publish(ctx context.Context, key string, message Message) error {
	var kafkaHeaders []kafka.Header
	for k, v := range message.Header {
		kafkaHeaders = append(kafkaHeaders, kafka.Header{Key: k, Value: []byte(v)})
	}

	msg := kafka.Message{
		Value:   message.Body,
		Key:     []byte(key),
		Time:    time.Now(),
		Headers: kafkaHeaders,
	}
	return p.writer.WriteMessages(ctx, msg)
}

func (p *PublisherOptions) PublishMultiRaw(ctx context.Context, key string, msgs []Message) error {
	messages := make([]kafka.Message, 0, len(msgs))

	for _, msg := range msgs {
		var kafkaHeaders []kafka.Header
		for k, v := range msg.Header {
			kafkaHeaders = append(kafkaHeaders, kafka.Header{Key: k, Value: []byte(v)})
		}

		messages = append(messages, kafka.Message{
			Value:   msg.Body,
			Key:     []byte(key),
			Time:    time.Now(),
			Headers: kafkaHeaders,
		})
	}
	return p.writer.WriteMessages(ctx, messages...)
}

func (p *PublisherOptions) newKafkaWriter() {
	p.writer = &kafka.Writer{
		Addr:                   kafka.TCP(p.Brokers...),
		Topic:                  p.Topic,
		BatchSize:              p.BatchSize,
		AllowAutoTopicCreation: true,
	}
}

func CreateTopic(address, topic string) error {
	conn, err := kafka.DialContext(context.Background(), "tcp", address)
	if err != nil {
		return err
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatalf("failed to close connection: %v", err)
		}
	}()

	controller, _ := conn.Controller()

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})

	return nil
}
