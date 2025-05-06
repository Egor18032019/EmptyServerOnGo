package main

import (
	"log"
	"my-go-webserver/controllers"
	"net/http"
)

func main() {

	http.HandleFunc("/login", controllers.LoginHandler)
	http.HandleFunc("/home", controllers.HomeHandler)
	//http.HandleFunc("/logout", controllers.LogoutHandler)//todo сделать

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
