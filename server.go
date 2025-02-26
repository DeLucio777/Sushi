package main

import (
	"fmt"
	"html/template"
	"net/http"
)


func handleRequest() {
	http.Handle("/Styles/", http.StripPrefix("/Styles/", http.FileServer(http.Dir("./Styles/"))))
	http.Handle("/Pages/", http.StripPrefix("/Pages/", http.FileServer(http.Dir("./Pages/"))))
	http.HandleFunc("/", index)

	err := http.ListenAndServe(":8080", nil)

}

func main() {




	handleRequest()
}