package main

import (
	"fmt"

	"github.com/IBM/sarama"
)

func main() {
	// Create a new Kafka producer
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer producer.Close()

	// Create a new message
	msg := &sarama.ProducerMessage{
		Topic: "my_topic",
		Value: sarama.StringEncoder("Hello, World!"),
	}

	// Send the message to Kafka
	_, _, err = producer.SendMessage(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Message sent successfully!")
}
