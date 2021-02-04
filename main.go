package main

import (
	"dbserver"
	. "helper"
	"log"
	"net/http"
	. "routes"
)

func init() {
	dbserver.DBInit()
}

func main() {
	defer Db.Close()
	log.Println("Server running....")

	http.HandleFunc("/", HandleUser)
	http.HandleFunc("/save-note", SaveNote)
	http.HandleFunc("/delete-note", DeleteNote)
	http.HandleFunc("/update-note", UpdateNote)

	http.ListenAndServe(":8080", nil)
}
