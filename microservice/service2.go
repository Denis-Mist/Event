package main

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)

type Request struct {
	Text string `json:"text"`
}

type Response struct {
	Text string `json:"text"`
}

func main() {
	// Create a new Kafka consumer
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer consumer.Close()

	// Subscribe to the request topic
	topics := []string{"request_topic"}
	partitions, err := consumer.Partitions(topics[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	// Consume requests from the topic
	for _, partition := range partitions {
		pc, err := consumer.ConsumePartition(topics[0], partition, sarama.OffsetNewest)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer pc.Close()

		for msg := range pc.Messages() {
			var request Request
			err := json.Unmarshal(msg.Value, &request)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("Received request: %s\n", request.Text)

			// Process the request and send a response
			response := Response{"good"}
			jsonResponse, err := json.Marshal(response)
			if err != nil {
				fmt.Println(err)
				return
			}
			kafkaMsg := &sarama.ProducerMessage{
				Topic: "response_topic",
				Value: sarama.ByteEncoder(jsonResponse),
			}
			producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer producer.Close()
			_, _, err = producer.SendMessage(kafkaMsg)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
