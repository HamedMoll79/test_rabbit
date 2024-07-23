package rabbit

import "test/adapter/broker/rabbit"

type Worker interface {
	StartPublish()
	StartConsume()
}

type RabbitMQWorker struct {
	client    *rabbit.Adapter
	queueName string
}

func New(client *rabbit.Adapter, queueName string) *RabbitMQWorker {
	return &RabbitMQWorker{client: client, queueName: queueName}
}
