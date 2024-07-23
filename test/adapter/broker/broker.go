package broker

type MessageBroker interface {
	Publish(message Message) error
	Consume(queueName string) (<-chan Message, error)
}
