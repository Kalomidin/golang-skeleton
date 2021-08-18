package main

import (
	"github.com/ridebeam/go-common/config"
	"github.com/ridebeam/go-common/grpc"
	"github.com/ridebeam/go-common/kafka"
	"github.com/ridebeam/go-common/kv"
)

// Config contains all defined configuration values
type Config struct {
	config.Service

	Monitor ConfigMonitor `yaml:"monitor"`
	Kafka   ConfigKafka   `yaml:"kafka"`
	DB      ConfigDB      `yaml:"db"`
	Server  grpc.Config   `yaml:"server"`
}

// ConfigMonitor contains monitoring specific config
type ConfigMonitor struct {
	StatsD            ConfigStatsD `yaml:"statsd"`
	TracingSampleRate float64      `yaml:"sampleRate" env:"TRACE_SAMPLE_RATE" env-default:"0.001"`
}

// ConfigStatsD contains StatsD specific config
type ConfigStatsD struct {
	Address    string `yaml:"address" env:"STATSD_ADDRESS" env-default:":8125"`
	TagsFormat string `yaml:"tagsFormat" env:"STATSD_TAGS_FORMAT" env-default:"cloudwatch"`
}

// ConfigKafka contains Kafka specific config
type ConfigKafka struct {
	Connection   kafka.Config `yaml:"connection"`
	TopicExample string       `yaml:"topicExample" env:"KAFKA_TOPIC_EXAMPLE" env-default:"example"`
}

type ConfigDB struct {
	EnabledInDevelopment bool              `yaml:"enabledInDevelopment"`
	Postgres             kv.ConfigPostgres `yaml:"postgres"`
	TableExample         string            `yaml:"tableExample" env:"PG_EXAMPLE" env-default:"example"`
}

func loadConfig() (Config, error) {
	var cfg Config
	err := config.ForService(&cfg, "change-me-service-name")

	return cfg, err
}
