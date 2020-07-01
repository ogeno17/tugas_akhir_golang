package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
)

// Menu is struct
type Menu struct {
	IDMenu   int    `json:"idMenu"`
	Nama     string `json:"nama"`
	Kategori string `json:"kategori"`
	Harga    int    `json:"harga"`
}

func connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/restaurant")

	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	http.HandleFunc("/lihat-menu", lihatMenu)
	http.HandleFunc("/cari-menu", cariMenu)
	http.HandleFunc("/edit-menu", editMenu)
	http.HandleFunc("/hapus-menu", hapusMenu)
	fmt.Println("Start Web")
	http.ListenAndServe(":8080", nil)
}

func lihatMenu(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	rows, err := db.Query("SELECT idMenu, nama, kategori, harga FROM menu")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer rows.Close()

	var result []Menu

	for rows.Next() {
		var each = Menu{}

		err = rows.Scan(&each.IDMenu, &each.Nama, &each.Kategori, &each.Harga)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = append(result, each)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}

	message, _ := json.Marshal(result)

	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

func cariMenu(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	url, err := url.Parse(r.RequestURI)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	param := url.Query()["nama"][0]

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	rows, err := db.Query("SELECT idMenu, nama, kategori, harga FROM menu WHERE nama like ?", "%"+param+"%")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer rows.Close()

	var result []Menu

	for rows.Next() {
		var each = Menu{}

		err = rows.Scan(&each.IDMenu, &each.Nama, &each.Kategori, &each.Harga)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = append(result, each)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}

	message, _ := json.Marshal(result)

	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

func editMenu(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var data Menu
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	query := "UPDATE menu SET nama = ?, kategori = ?, harga = ? WHERE idMenu = ?"
	_, err = db.Exec(query, data.Nama, data.Kategori, data.Harga, data.IDMenu)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data berhasil diperbarui."}`))
}

func hapusMenu(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var data map[string]int
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	query := "DELETE FROM menu WHERE idMenu = ?"
	_, err = db.Exec(query, data["idMenu"])

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Data berhasil dihapus."}`))
}
