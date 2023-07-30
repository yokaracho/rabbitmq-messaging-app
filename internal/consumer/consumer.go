package main

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Consumer Application")

	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		"TestQueue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			fmt.Printf("Received messages:%s\n", d.Body)
		}
	}()
	fmt.Println("successfully connectto rabbitmq")
	fmt.Println("[*] - waiting for messages")
	<-forever
}
