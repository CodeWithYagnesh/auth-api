package main

import (
	"gin_jwt/app"
	"log"
)

func main() {
	r := app.SetupRouter()

	log.Println("Server started on port 80!")
	r.Run()
}
