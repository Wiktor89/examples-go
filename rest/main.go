package main

import (
	"encoding/json"
	"examples-go/common/model"
	_ "examples-go/rabbitmq/service"
	service2 "examples-go/rabbitmq/service"
	"examples-go/rest/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
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
	r := gin.Default()
	r.GET(ApiV1+"/products", func(c *gin.Context) {
		c.JSONP(http.StatusOK, service.GetAll())
	})
	r.POST(ApiV1+"/products", func(c *gin.Context) {
		p2 := model.Product{}
		service.CheckError(c.ShouldBindJSON(&p2))
		sendQ(p2, os.Getenv("PRODUCER_Q"))
		c.Status(http.StatusOK)
	})
	r.Run(":8080")
}

func ListenerQ() {
	service2.QueueDeclare(os.Getenv("DECLARE_Q"))
	go func() {
		for {
			message := service2.ConsumerMessage(os.Getenv("CONSUMER_Q"))
			if message != nil {
				p := model.Product{}
				err := json.Unmarshal(message, &p)
				service.CheckError(err)
				_, err = service.Create(p.Name)
				service.CheckError(err)
				log.Println("[x] ", p)
			} else {
				log.Println("NO message from que")
			}
			log.Println("Before")
			time.Sleep(time.Minute / 6)
			log.Println("After")
		}
	}()
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error load .env")
	}
}

func sendQ(product model.Product, nameQ string) {
	p, err := json.Marshal(product)
	service.CheckErrorRabbitMq(err)
	if err == nil {
		service2.SendMessage(p, nameQ)
	}
}
