package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ENV     string `default:"dev" envconfig:"ENV"`
	AppPort string `default:"3001" envconfig:"APP_PORT"`

	Kafka Kafka
	Mongo Mongo
	Gcs   GCS
}

type Mongo struct {
	URL string `default:"mongodb://localhost:27017" envconfig:"MONGO_CONNECTION_STRING"`
	DB  string `default:"sideproject" envconfig:"MONGO_DB_NAME"`
}

type Kafka struct {
	Brokers       string `default:"localhost:9093" envconfig:"KAFKA_BROKERS"`
	ConsumerGroup string `default:"group_notify" envconfig:"KAFKA_CONSUMER_GROUP"`
	Topic         string `default:"like_event" envconfig:"KAFKA_TOPIC"`
}

type GCS struct {
	GoogleCredFile   string `default:"gcs_creds.json"`
	GoogleBucketName string `default:"mimetic-might-148107.appspot.com"`
	GoogleHost       string `default:"https://storage.cloud.google.com"`
}

func Load() (*Config, error) {
	var conf Config
	if err := envconfig.Process("sideproject", &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func MustLoad() *Config {
	conf, err := Load()
	if err != nil {
		panic(err)
	}
	return conf
}
