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
	ItemType    string  `json:"item_type"`
}

type UsersToItems struct {
	Id     int    `json:"id"`
	UserId int    `json:"user_id"`
	Items  []Item `json:"items"`
}

func userHandel(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUserHandel(w, r)

	case http.MethodPost:
		setUserHandel(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func itemsHandel(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getItemsHandel(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func userToItemsHandel(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getItemsToUsersHandel(w, r)
	case http.MethodPost:
		addItemToUserHandel(w, r)
	case http.MethodDelete:
		fmt.Println("s")
		deleteItemFromUserHandel(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func addItemToUserHandel(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")
	itemId := r.URL.Query().Get("item_id")
	fmt.Println(userId, itemId)
	if userId == "" || itemId == "" {
		http.Error(w, "user_id и item_id обязательны", http.StatusBadRequest)
		return
	}

	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO tbl_users_to_items (user_id, item_id) VALUES ($1, $2)", userId, itemId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func deleteItemFromUserHandel(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")
	itemId := r.URL.Query().Get("item_id")
	fmt.Println("do")

	if userId == "" || itemId == "" {
		http.Error(w, "user_id и item_id обязательны", http.StatusBadRequest)
		return
	}

	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec(`
    DELETE FROM tbl_users_to_items 
    WHERE ctid IN (
        SELECT ctid 
        FROM tbl_users_to_items 
        WHERE user_id = $1 AND item_id = $2 
        LIMIT 1
    )
`, userId, itemId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func deleteAllItemFromUserHandel(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")
	fmt.Println("do")

	if userId == "" {
		http.Error(w, "user_id обязательны", http.StatusBadRequest)
		return
	}

	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec(`
    DELETE FROM tbl_users_to_items 
    WHERE ctid IN (
        SELECT ctid 
        FROM tbl_users_to_items 
        WHERE user_id = $1
    )
`, userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func getItemsHandel(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM tbl_item")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	items := []Item{}
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Image, &item.Description, &item.Name, &item.Cost, &item.ItemType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)

}

func getUserHandel(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	name := r.URL.Query().Get("name")
	password := r.URL.Query().Get("password")

	var id int
	err = db.QueryRow("SELECT ID FROM tbl_users WHERE name = $1 AND password = $2", name, password).Scan(&id)
	exists := id != 0
	if exists {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"id": id})
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"id": -1})
	}

}

func setUserHandel(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	defer db.Close()
	name := r.URL.Query().Get("name")
	password := r.URL.Query().Get("password")

	fmt.Println(name + " " + password)
	var id int
	err = db.QueryRow("SELECT ID FROM tbl_users WHERE name = $1 AND password = $2", name, password).Scan(&id)
	exists := id != 0
	if exists {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Query(`INSERT INTO tbl_users (name, password) VALUES ($1, $2)`, name, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nil)
}

func getItemsToUsersHandel(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT ID, user_id, item_id FROM tbl_users_to_items")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var usersToItems []UsersToItems
	for rows.Next() {
		var itemId int
		var userToItem UsersToItems
		err := rows.Scan(&userToItem.Id, &userToItem.UserId, &itemId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var items []Item
		itemsRows, err := db.Query("SELECT ID, image, description, item_name, cost, item_type FROM tbl_item where ID = $1", itemId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for itemsRows.Next() {
			var item Item
			err := itemsRows.Scan(&item.ID, &item.Image, &item.Description, &item.Name, &item.Cost, &item.ItemType)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			items = append(items, item)
		}
		userToItem.Items = items
		usersToItems = append(usersToItems, userToItem)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usersToItems)
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
	http.HandleFunc("/api/items", itemsHandel)
	http.HandleFunc("/api/users", userHandel)
	http.HandleFunc("/api/items_to_users", userToItemsHandel)
	http.HandleFunc("/api/items_to_users_delete", deleteItemFromUserHandel)
	http.HandleFunc("/api/items_to_users_delete_all", deleteAllItemFromUserHandel)
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
