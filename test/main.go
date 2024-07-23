package main

import (
	"test/adapter/broker/rabbit"
	rabbit2 "test/worker/rabbit"
)

func main() {
	cfg := rabbit.Config{
		Host:        "localhost",
		Port:        5672,
		Username:    "rabbitmq",
		Password:    "rabbitmq",
		VirtualHost: "/",
	}

	adapter, err := rabbit.New(cfg, "develop", []string{"test"})
	if err != nil {
		panic(err)
	}

	rabbitWorker := rabbit2.New(adapter, "test")
	go rabbitWorker.StartPublish()
	rabbitWorker.StartConsume()
}
