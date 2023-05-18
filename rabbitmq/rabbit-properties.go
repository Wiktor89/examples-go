package rabbitmq

import (
	"context"
	er "examples-go/checks/error"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

func createConnection() *amqp091.Connection {
	connection, err := amqp091.Dial(os.Getenv("RABBIT"))
	er.CheckErrorRabbitMq(err)
	return connection
}

func SendMessage(msg []byte, nameQ string) {
	connection := createConnection()
	defer connection.Close()
	channel, err := connection.Channel()
	defer channel.Close()
	er.CheckErrorRabbitMq(err)
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
	er.CheckErrorRabbitMq(err)
	defer channel.Close()
	consume, err := channel.Consume(nameQ,
		"",
		true,
		false,
		false,
		false,
		nil)
	er.CheckErrorRabbitMq(err)
	for msg := range consume {
		return msg.Body
	}
	return nil
}

func QueueDeclare(nameQ string) {
	connection := createConnection()
	defer func(connection *amqp091.Connection) {
		err := connection.Close()
		er.CheckErrorRabbitMq(err)
	}(connection)
	channel, err := connection.Channel()
	defer func(channel *amqp091.Channel) {
		err := channel.Close()
		er.CheckErrorRabbitMq(err)
	}(channel)
	er.CheckErrorRabbitMq(err)
	_, err = channel.QueueDeclare(nameQ, false, false, false, false, nil)
	er.CheckErrorRabbitMq(err)
}
