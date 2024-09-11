package main

import (
	"fmt"

	"github.com/IBM/sarama"
)

func main() {
	// Create a new Kafka consumer
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer consumer.Close()

	// Subscribe to the topic
	topics := []string{"my_topic"}
	partitions, err := consumer.Partitions(topics[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	// Consume messages from the topic
	for _, partition := range partitions {
		pc, err := consumer.ConsumePartition(topics[0], partition, sarama.OffsetNewest)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer pc.Close()

		for msg := range pc.Messages() {
			fmt.Printf("Received message: %s\n", string(msg.Value))
		}
	}
}
