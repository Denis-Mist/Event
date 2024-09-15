package main

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/dgrijalva/jwt-go"
)

// KafkaConsumer represents a Kafka consumer
type KafkaConsumer struct {
	brokerURLs []string
	topic      string
}

// NewKafkaConsumer returns a new Kafka consumer
func NewKafkaConsumer(brokerURLs []string, topic string) *KafkaConsumer {
	return &KafkaConsumer{
		brokerURLs: brokerURLs,
		topic:      topic,
	}
}

// Consume consumes messages from the Kafka topic
func (kc *KafkaConsumer) Consume(ctx context.Context) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	client, err := sarama.NewClient(kc.brokerURLs, config)
	if err != nil {
		log.Fatal(err)
	}

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Fatal(err)
	}

	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(ctx, kc.topic, sarama.OffsetNewest)
	if err != nil {
		log.Fatal(err)
	}

	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			// Verify JWT token
			token, err := jwt.Parse(msg.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte("your-secret-key"), nil
			})
			if err != nil {
				log.Println(err)
				continue
			}

			// Perform business logic based on the verified JWT token
			// ...

		case err := <-partitionConsumer.Errors():
			log.Println(err)
		}
	}
}

func main() {
	brokerURLs := []string{"kafka-broker-1:9092", "kafka-broker-2:9092"}
	topic := "your-kafka-topic"

	kc := NewKafkaConsumer(brokerURLs, topic)
	kc.Consume(context.Background())
}
