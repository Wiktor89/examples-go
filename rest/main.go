package main

import (
	"encoding/json"
	er "examples-go/checks/error"
	model "examples-go/models"
	"examples-go/rabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

const (
	ApiV1 = "/api/v1"
)

func main() {
	LoadEnv()
	ListenerQ()
	AppRun()
}

func AppRun() {
	log.Println("[x] Configuration http server")
	r := gin.Default()
	r.GET(ApiV1+"/products", func(c *gin.Context) {
		c.JSONP(http.StatusOK, GetAll())
	})
	r.POST(ApiV1+"/products", func(c *gin.Context) {
		p2 := model.Product{}
		er.CheckErrorHttp(c.ShouldBindJSON(&p2))
		sendQ(p2, os.Getenv("PRODUCER_Q"))
		c.Status(http.StatusOK)
	})
	r.Run(":" + os.Getenv("SERVER_PORT"))
}

func ListenerQ() {
	log.Println("[x] Configuration listener rabbitMQ")
	rabbitmq.QueueDeclare(os.Getenv("DECLARE_Q"))
	go func() {
		for {
			message := rabbitmq.ConsumerMessage(os.Getenv("CONSUMER_Q"))
			if message != nil {
				p := model.Product{}
				err := json.Unmarshal(message, &p)
				er.CheckErrorRabbitMq(err)
				_, err = Create(p.Name)
				er.CheckErrorRabbitMq(err)
				log.Println("[x] ", p)
			} else {
				log.Println("[x] NO message from que")
			}
		}
	}()
}

func LoadEnv() {
	log.Println("[x] Load env")
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error load .env")
	}
}

func sendQ(product model.Product, nameQ string) {
	p, err := json.Marshal(product)
	er.CheckErrorRabbitMq(err)
	if err == nil {
		rabbitmq.SendMessage(p, nameQ)
	}
}
