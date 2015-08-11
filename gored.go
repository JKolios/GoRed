package main

import (
	"fmt"
	"github.com/Jkolios/GoRed/handlers"
	"github.com/mediocregopher/radix.v2/redis"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type connectionDetails struct {
	host, port string
}

func initConnection(connDetails connectionDetails) *redis.Client {
	// Connection Establishment

	client, err := redis.Dial("tcp", connDetails.host+":"+connDetails.port)
	if err != nil {
		fmt.Println("Error while connecting:" + err.Error())
		return nil
	}
	return client
}

func main() {
	// Redis connection init
	redisConnDetails := connectionDetails{"jspi2.local", " 6379"}
	client := initConnection(redisConnDetails)
	defer client.Close()

	// HTTP server init
	httpConnDetails := connectionDetails{host: "", port: "8080"}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, os.Kill)

	var handler handlers.SimpleRedisHandler
	handler.Client = client
	go handler.StoreCount(signalChannel)
	handler.GetCurrent()
	http.Handle("/", &handler)
	log.Fatal(http.ListenAndServe(httpConnDetails.host+":"+httpConnDetails.port, nil))

}
