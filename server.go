package main

import (
	"database/sql"
	"encoding/json"
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

func getUserByInputParamsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	sPassword := r.URL.Query().Get("u_name")
	if sPassword == "" {
		http.Error(w, "Параметр имя обязателен", http.StatusBadRequest)
		return
	}
	sName := r.URL.Query().Get("u_password")
	if sName == "" {
		http.Error(w, "Параметр пароль обязателен", http.StatusBadRequest)
		return
	}

	res, err := getUserByInputParams(db, sName, sPassword)
	if err != nil {
		http.Error(w, "Ошибка получения данных занятий", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func getUserByInputParams(db *sql.DB, sName string, sPassword string) (bool, error) {

	query := `SELECT * FROM users WHERE name = $1 AND password = $2`

	rows, err := db.Query(query, sName, sPassword)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if !rows.Next() {
		return false, nil 
	}

	for rows.Next() {

		var student User
		err := rows.Scan(&student.ID, &student.Name, &student.Password)
		if err != nil {
			return false, err
		}

	}

	return true, nil
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
	http.HandleFunc("/api/user_by_params", getUserByInputParamsHandler)

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
