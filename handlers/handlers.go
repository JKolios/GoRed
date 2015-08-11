package handlers

import (
	"fmt"
	"github.com/mediocregopher/radix.v2/redis"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

/*SimpleRedisHandler is a  Redis page access counter*/
type SimpleRedisHandler struct {
	Client       *redis.Client
	currentValue int
}

func (handler *SimpleRedisHandler) ServeHTTP(respWriter http.ResponseWriter, request *http.Request) {
	log.Println("Handler called")
	handler.currentValue++
	log.Println(handler.currentValue)
	fmt.Fprintf(respWriter, "%d", handler.currentValue)
}

func (handler *SimpleRedisHandler) GetCurrent() {
	DBValue, err := handler.Client.Cmd("GET", "access_counter").Str()
	if err != nil {
		log.Println("Error while executing command:" + err.Error())
		return
	}
	handler.currentValue, err = strconv.Atoi(DBValue)
	if err != nil {
		log.Println("Error whileconverting response:" + err.Error())
		return
	}
}

func (handler *SimpleRedisHandler) StoreCount(signalChannel chan os.Signal) {
	// Store the state of the handler to Redis upon exiting
	<-signalChannel
	signal.Stop(signalChannel)

	err := handler.Client.Cmd("SET", "access_counter", handler.currentValue).Err
	if err != nil {
		log.Println("Error while executing command:" + err.Error())
		return
	}
}
