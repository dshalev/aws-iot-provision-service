package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"github.com/dshalev2/aws-iot-provision-service/handlers"
)


func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/provision/{thingName}", handlers.HandleProvision)
	log.Fatal(http.ListenAndServe(":8081", router))
}
