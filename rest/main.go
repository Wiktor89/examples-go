package main

import (
	"examples-go/rest/service"
	"fmt"
)

func main() {
	//u := uuid.New().String()
	//service.Create(u)
	//fmt.Println(service.GetById(3))
	fmt.Println(service.GetAll())
}
