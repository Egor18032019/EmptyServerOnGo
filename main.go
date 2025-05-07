package main

import (
	"log"
	"my-go-webserver/controllers"
	"my-go-webserver/services"
	"net/http"
	"strconv"
)

const port = 8080

func main() {
	if err := services.InitLogging(); err != nil {
		log.Fatalf("Не удалось инициализировать логирование: %v", err)
	}

	http.HandleFunc("/login", controllers.LoginHandler)
	http.HandleFunc("/home", controllers.HomeHandler)
	//http.HandleFunc("/logout", controllers.LogoutHandler)//todo сделать

	services.LogMessage("Server starting on : " + strconv.Itoa(port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
