package rabbit

import (
	"encoding/json"
	"fmt"
	"test/adapter/broker"
	"time"
)

type Body struct {
	Size  string `json:"size"`
	Color string `json:"color"`
}

func (w *RabbitMQWorker) StartPublish() {
	for {
		messageBody := Body{
			Size:  "15",
			Color: "red",
		}
		marsh, err := json.Marshal(messageBody)
		if err != nil {
			fmt.Printf("cant marshal messageBody: %s\n", err)
		}
		message := broker.Message{
			ID:          time.Now().String(),
			Body:        marsh,
			MessageType: "test",
			QueueName:   w.queueName,
		}
		err = w.client.Publish(message)
		if err != nil {
			fmt.Printf("cant publish message: %s\n", err)
		}

		time.Sleep(5 * time.Second)
	}
}
