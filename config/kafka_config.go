package config

var (
	// KafkaConfig var
	// manage configs for kafka
	KafkaConfig kafkaInterface = &kafka{}
)

type kafkaInterface interface{}

type kafka struct{}
