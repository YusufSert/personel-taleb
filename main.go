package main

import (
	"log"
	"websocketTest/api"
)

func main() {

	cRest, err := api.NewCihazlarRest()
	if err != nil {
		log.Println(err)
	}

	err = cRest.Run()
	log.Println(err)
}
