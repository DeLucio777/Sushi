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

func userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Println("do1")
		getUserByInputParamsHandler(w, r)
	case http.MethodPost:
		setUserByInputParamsHandler(w, r)
	default:
		fmt.Println("do2")

		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func getUserByInputParamsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	sName := r.URL.Query().Get("user_name")
	if sName == "" {
		http.Error(w, "Параметр имя обязателен", http.StatusBadRequest)
		return
	}
	sPassword := r.URL.Query().Get("user_password")
	if sPassword == "" {
		http.Error(w, "Параметр пароль обязателен", http.StatusBadRequest)
		return
	}
	fmt.Println(sName, sPassword)
	res, err := getUserByInputParams(db, sName, sPassword)
	if err != nil {
		http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
		return
	}
	if res == nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func getUserByInputParams(db *sql.DB, sName string, sPassword string) (*User, error) {
	query := `SELECT * FROM tbl_users WHERE name = $1 AND password = $2`

	rows, err := db.Query(query, sName, sPassword)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var user User

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Name, &user.Password)

		if err != nil {
			return nil, err
		}

		return &user, nil
	}
	return nil, fmt.Errorf("пользователь не найден")
}

func setUserByInputParamsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var input struct {
		UserName     string `json:"user_name"`
		UserPassword string `json:"user_password"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Ошибка декодирования данных", http.StatusBadRequest)
		return
	}

	sName := input.UserName
	sPassword := input.UserPassword

	if sName == "" {
		http.Error(w, "Параметр имя обязателен", http.StatusBadRequest)
		return
	}
	if sPassword == "" {
		http.Error(w, "Параметр пароль обязателен", http.StatusBadRequest)
		return
	}

	exists, err := userExists(db, sName, sPassword)
	if err != nil {
		http.Error(w, "Ошибка проверки пользователя", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Пользователь уже существует", http.StatusConflict)
		return
	}

	err = insertUser(db, sName, sPassword)
	if err != nil {
		http.Error(w, "Ошибка добавления пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated) // Успешное создание
	json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь успешно создан"})
}

func userExists(db *sql.DB, sName string, sPassword string) (bool, error) {
	query := `SELECT COUNT(*) FROM tbl_users WHERE name = $1 and password = $2`
	var count int
	err := db.QueryRow(query, sName, sPassword).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func insertUser(db *sql.DB, sName string, sPassword string) error {
	query := `INSERT INTO tbl_users (name, password) VALUES ($1, $2)`
	_, err := db.Exec(query, sName, sPassword)
	return err
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getItemsHandler(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
func getItemsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	result, err := getItems(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func getItems(db *sql.DB) ([]Item, error) {
	query := `SELECT * FROM tbl_item`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Item
	for rows.Next() {
		var item Item
		err = rows.Scan(&item.ID, &item.Image, &item.Description, &item.Name, &item.Cost)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func treshHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTreshHandler(w, r)
	case http.MethodPost:
		setThreshHandler(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func setThreshHandler(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	sName := r.URL.Query().Get("user_name")
	if sName == "" {
		http.Error(w, "Параметр имя обязателен", http.StatusBadRequest)
		return
	}
	res, err := getTresh(db, sName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func getTreshHandler(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	sName := r.URL.Query().Get("user_name")
	if sName == "" {
		http.Error(w, "Параметр имя обязателен", http.StatusBadRequest)
		return
	}
	res, err := getTresh(db, sName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)

}

func getTresh(db *sql.DB, sName string) ([]Item, error) {

	user_id, err := db.Query("SELECT tbl_item.id, tbl_item.image, tbl_item.description, tbl_item.item_name, tbl_item.cost FROM tbl_item JOIN tbl_users_to_items ON (tbl_item.id = tbl_users_to_items.item_id) JOIN tbl_users ON tbl_users_to_items.user_id = tbl_users.id WHERE tbl_users.name=$1", sName)
	if err != nil {
		return nil, err
	}
	defer user_id.Close()
	var items []Item
	for user_id.Next() {
		var item Item
		err = user_id.Scan(&item.ID, &item.Image, &item.Description, &item.Name, &item.Cost)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil

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
	http.HandleFunc("/api/users", userHandler)
	http.HandleFunc("/api/items", itemsHandler)
	http.HandleFunc("/api/tresh", treshHandler)
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
