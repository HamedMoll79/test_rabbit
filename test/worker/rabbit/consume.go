package rabbit

import "fmt"

func (w *RabbitMQWorker) StartConsume() {
	fmt.Printf("start consume ...")
	ch, err := w.client.Consume(w.queueName)
	if err != nil {
		fmt.Printf("cant get client consume error\n")
		panic(err)
	}

	for msg := range ch {
		fmt.Printf("got message: %s\n", msg.Body)
	}
}
