package rabbit

import (
	"encoding/base64"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"log/slog"
	"test/adapter/broker"
)

const RabbitMQNS = "rabbitmq"

// TODO: add ack and publish metrics
type Config struct {
	Host        string `koanf:"host"`
	Port        int    `koanf:"port"`
	Username    string `koanf:"username"`
	Password    string `koanf:"password"`
	VirtualHost string `koanf:"virtual_host"`
}

type Adapter struct {
	cfg     Config
	conn    *amqp.Connection
	channel *amqp.Channel
	runMod  string
}

func New(config Config, runMod string, queueNames []string) (*Adapter, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.VirtualHost))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}
	for _, queueName := range queueNames {
		fmt.Printf("queue declare %s \n", queueName)
		_, err := channel.QueueDeclare(
			queueName,
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			return nil, fmt.Errorf("queue declare %s %v", queueName, slog.String("error", err.Error()))
		}
	}
	adapter := &Adapter{
		conn:    conn,
		channel: channel,
		cfg:     config,
		runMod:  runMod,
	}
	return adapter, nil
}

func (a *Adapter) Publish(message broker.Message) error {

	base64EncodedMessage := base64.StdEncoding.EncodeToString(message.Body)

	err := a.channel.Publish(
		"",                // exchange
		message.QueueName, // routing key (same as queue name)
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			MessageId:    message.ID,
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/protobuf",
			Type:         message.MessageType,
			Body:         []byte(base64EncodedMessage),
		})
	// If there is an error publishing the message, a log will be displayed in the terminal.
	if err != nil {
		log.Println("publish rabbit - publish", slog.String("error", err.Error()))
		return err
	}

	log.Println("Sending message to Rabbitmq ...")
	return nil
}

func (a *Adapter) Consume(queueName string) (<-chan broker.Message, error) {
	log.Println("Consumeq ...", queueName)
	messages, err := a.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Println("Consumeq err...", err)
		return nil, err
	}

	messageChan := make(chan broker.Message)
	//TODO: Add communication with  supervisor
	go func() {
		for delivery := range messages {
			decodedMessage, err := base64.StdEncoding.DecodeString(string(delivery.Body))
			if err != nil {
				// Handle error
				log.Println("Error decoding base64 message:", err.Error())
				return
			}
			message := broker.Message{
				ID:          delivery.MessageId,
				Body:        decodedMessage,
				MessageType: delivery.Type,
				QueueName:   delivery.RoutingKey,
			}
			messageChan <- message
			// Acknowledge the message to remove it from the queue
			err = a.Ack(delivery.DeliveryTag)
			if err != nil {
				log.Printf("Error acknowledging message: %v\n", err)
			}
		}
	}()

	return messageChan, nil
}

func (a *Adapter) Ack(deliveryTag uint64) error {
	return a.channel.Ack(deliveryTag, false)
}
