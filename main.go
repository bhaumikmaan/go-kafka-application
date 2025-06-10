package main

import (
	"fmt"
	"go-kafka-application/consumer"
	"go-kafka-application/producer"
	"os"
)

func main() {
	arg := os.Args[1:]
	if len(arg) == 0 {
		menu()
	} else {
		action := arg[0]
		switch action {
		case "close":
			fmt.Println("Closing the applcation! Bye.")
			os.Exit(0)
		case "list-topics":
			producer.PrintTopics()
		case "start-producer":
			producer.ProduceMessages(arg[1])
		case "start-consumer":
			consumer.ConsumeMessages(arg[1])
		default:
			fmt.Println("Unknown action. Please input again.")
			menu()
		}
	}
}

func menu() {
	fmt.Println("Welcome to the Go-Kafka Application! Please enter the argument for the action you want to perform")
	fmt.Println("1. list-topics - Lists all the topics\n" +
		"2. start-producer - Reads the stream and produces in topic\n" +
		"3. start-consumer - Consumes a new topic")
}
