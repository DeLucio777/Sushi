package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"net/http"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
type Item struct {
	ID          int     `json:"id"`
	Image       string  `json:"image"`
	Description string  `json:"description"`
	Name        string  `json:"item_name"`
	Cost        float64 `json:"cost"`
}

func connectDB() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных:", err)
		return nil, err
	}
	// Проверяем соединение
	if err = db.Ping(); err != nil {
		fmt.Println("Не удалось подключиться к базе данных:", err)
		return nil, err
	}
	return db, nil
}

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

	db, err := connectDB()
	if err != nil {
		fmt.Println("Не удалось подключиться к базе данных. Завершение работы...")
		return
	}
	defer db.Close()

	handleRequest()
}
