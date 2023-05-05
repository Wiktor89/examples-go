package service

import (
	"context"
	"examples-go/rest/service"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

func createConnection() *amqp091.Connection {
	connection, err := amqp091.Dial(os.Getenv("RABBIT"))
	service.CheckError(err)
	return connection
}

func SendMessage(msg []byte, nameQ string) {
	connection := createConnection()
	defer connection.Close()
	channel, err := connection.Channel()
	defer channel.Close()
	service.CheckErrorRabbitMq(err)
	err = channel.PublishWithContext(context.Background(),
		"",
		nameQ,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
	if err != nil {
		log.Println("Error while send message =", err)
	}
}

func ConsumerMessage(nameQ string) []byte {
	connection := createConnection()
	defer connection.Close()
	channel, err := connection.Channel()
	service.CheckErrorRabbitMq(err)
	defer channel.Close()
	consume, err := channel.Consume(nameQ,
		"",
		true,
		false,
		false,
		false,
		nil)
	service.CheckErrorRabbitMq(err)
	for msg := range consume {
		return msg.Body
	}
	return nil
}

func QueueDeclare(nameQ string) {
	connection := createConnection()
	defer func(connection *amqp091.Connection) {
		err := connection.Close()
		service.CheckError(err)
	}(connection)
	channel, err := connection.Channel()
	defer func(channel *amqp091.Channel) {
		err := channel.Close()
		service.CheckError(err)
	}(channel)
	service.CheckError(err)
	_, err = channel.QueueDeclare(nameQ, false, false, false, false, nil)
	service.CheckError(err)
}
