package main

import (
	"log"
	"net/http"
	"github.com/Hossam16/go-chat-creation-api/internal/router"
)

func main() {
	r := router.InitRouter()

	log.Println("Listening on 8080 ......")
	log.Fatal(http.ListenAndServe(":8080", r))
}
