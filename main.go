package main

import (
	"fmt"
	"log"
	"net/http"
	"stocks-api/router"
)

func main() {
	r := router.Router()
	fmt.Println("Starting server on port 5600...")
	log.Fatal(http.ListenAndServe(":5600", r))
}
