package error

import (
	"database/sql"
	"log"
)

func CheckErrorHttp(err error) {
	if err != nil {
		log.Println("Http error =", err)
	}
}

func CheckErrorRabbitMq(err error) {
	if err != nil {
		log.Println("RabbitMq =", err)
	}
}

func CheckErrorDb(err error) {
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			log.Println("No rows were returned!")
			break
		case nil:
			break
		default:
			log.Println("Unable to scan the row. ", err)
		}
	}
}
