package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Game starting")
	SetupAPI()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func SetupAPI() {
	// ctx := context.Background()
	manager := NewManager()

	http.HandleFunc("/ws", manager.ServeWS)

}
