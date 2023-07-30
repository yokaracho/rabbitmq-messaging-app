package main

import (
	"bufio"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"time"
)

type Message struct {
	Content  string
	Author   string
	SentTime time.Time
}

var messageHistory []Message
var totalMessagesSent int
var totalResponseTime time.Duration

func historyMessage(displayMessage []Message) {
	fmt.Println("Hisory messages:")
	for i, message := range displayMessage {
		fmt.Printf("%d. %s\n", i+1, message)
	}
	fmt.Println("Input message (to exit input - 'exit' and 'history' to output history messages):")
}

func main() {
	fmt.Println("GO  RabbitMQ")

	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()
	fmt.Println("Susccesful connect to Rabbitmq")

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"TestQueue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(q)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Input message (to exit input - 'exit' and 'history' to output history messages):")

	for scanner.Scan() {
		message := scanner.Text()
		if message == "exit" {
			break
		} else if message == "history" {
			historyMessage(messageHistory)
			continue
		}

		startTime := time.Now()

		err = ch.Publish(
			"",
			"TestQueue",
			false,
			false,
			amqp091.Publishing{
				ContentType: "text/plain",
				Body:        []byte(message),
			},
		)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		//messageHistory = append(messageHistory, Message{
		//	Content:  message,
		//	Author:   "Volodya",
		//	SentTime: time.Now(),
		//})

		messageHistory = append(messageHistory, Message{
			Content:  message,
			Author:   "Volodya",
			SentTime: startTime,
		})

		elapsedTime := time.Since(startTime)
		totalResponseTime += elapsedTime
		totalMessagesSent++
	}

	fmt.Printf("Total messages sent: %d\n", totalMessagesSent)
	if totalMessagesSent > 0 {
		averageResponseTime := totalResponseTime / time.Duration(totalMessagesSent)
		fmt.Printf("Average response time from RabbitMQ: %s\n", averageResponseTime)
	} else {
		fmt.Println("Total messages sent < 0")
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading standard input:", err)
	}
}
