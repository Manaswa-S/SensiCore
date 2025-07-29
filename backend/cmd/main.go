package main

import (
	"fmt"
	"os"
	"os/signal"
	"sensicore/cmd/db"
	"sensicore/cmd/server"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {

	fmt.Println("Starting Server...")

	flowChan := make(chan os.Signal, 1)
	signal.Notify(flowChan, syscall.SIGINT, syscall.SIGTERM)

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	dataStore, err := db.NewDataStore()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = server.InitHTTPServer(dataStore)
	if err != nil {
		fmt.Println(err)
		return
	}

	<-flowChan

	fmt.Println("Shutting server down...")

	if err := db.Close(dataStore); err != nil {
		fmt.Println(err)
	}
}
