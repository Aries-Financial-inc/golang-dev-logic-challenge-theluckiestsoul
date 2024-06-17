package main

import (
	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-theluckiestsoul/routes"
)

func main() {
	router := routes.SetupRouter()
	router.Run(":8080")
}
