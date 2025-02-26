package main

import (
	"fmt"
	"html/template"
	"net/http"
)


func index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("Pages/main.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}


func handleRequest() {
	http.Handle("/Styles/", http.StripPrefix("/Styles/", http.FileServer(http.Dir("./Styles/"))))
	http.Handle("/Pages/", http.StripPrefix("/Pages/", http.FileServer(http.Dir("./Pages/"))))
	http.HandleFunc("/", index)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("connect error")
	}
}

func main() {

fmt.Println("start")


	handleRequest()
}