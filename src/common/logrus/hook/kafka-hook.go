package hook

import (
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// KafkaHook implements Hook interface of logrus
type KafkaHook struct {
	// Id of the hook
	id string

	// Log levels allowed
	levels []logrus.Level

	// Log entry formatter
	formatter logrus.Formatter

	// sarama.AsyncProducer
	producer sarama.AsyncProducer

	// default topics
	defaultTopics []string
}

// Id returns hook's ID
func (hook *KafkaHook) Id() string {
	return hook.id
}

// Levels returns applicable levels
func (hook *KafkaHook) Levels() []logrus.Level {
	return hook.levels
}

// Fire is triggered when data is to be pushed
func (hook *KafkaHook) Fire(entry *logrus.Entry) error {
	// Check time for partition key
	var partitionKey sarama.ByteEncoder

	// Get field time, should work fine
	t, _ := entry.Data["time"].(time.Time)

	// Convert it to bytes
	b, err := t.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "Error Marshaling Binary")
	}

	partitionKey = sarama.ByteEncoder(b)

	// Check topics
	var topics []string

	if ts, ok := entry.Data["topics"]; ok {
		if topics, ok = ts.([]string); !ok {
			return errors.New("Field topics must be []string")
		}
	} else {
		// Load default topics, if not passed as argument
		topics = hook.defaultTopics
	}

	// Format before writing
	b, err = hook.formatter.Format(entry)

	if err != nil {
		return err
	}

	value := sarama.ByteEncoder(b)

	for _, topic := range topics {
		hook.producer.Input() <- &sarama.ProducerMessage{
			Key:   partitionKey,
			Topic: topic,
			Value: value,
		}
	}

	return nil
}

// NewKafkaHook creates a new hook for Kafka
func NewKafkaHook(
	id string,
	levels []logrus.Level,
	formatter logrus.Formatter,
	brokers []string,
	topics []string,
) (*KafkaHook, error) {
	kafkaConfig := sarama.NewConfig()
	// Only wait for the leader to ack
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForLocal

	// Compress messages
	kafkaConfig.Producer.Compression = sarama.CompressionSnappy

	// Flush batches every 500ms
	kafkaConfig.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(brokers, kafkaConfig)
	if err != nil {
		return nil, err
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			log.Printf(
				"Error sending message to Kafka: %v\n",
				err,
			)
		}
	}()

	return &KafkaHook{
		id,
		levels,
		formatter,
		producer,
		topics,
	}, nil
}
