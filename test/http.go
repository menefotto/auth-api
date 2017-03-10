package main

import (
	"log"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Test ok!"))
}

func main() {
	http.HandleFunc("/test", Test)
	err := http.ListenAndServe(":8080", nil)
	log.Println("server started")
	if err != nil {
		log.Fatalln(err)
	}
}
