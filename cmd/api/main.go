package main

import (
	"log"
	"restapirian/config/router"
)

func main() {
	if err := router.InitializeRouter(); err != nil {
		log.Fatal(err, "err initialize apps ")
	}
	// router.InitializeRouter()
}
