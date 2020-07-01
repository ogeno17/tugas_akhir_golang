package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var baseURL = "http://localhost:8080"

// Menu is struct
type Menu struct {
	IDMenu   int    `json:"idMenu"`
	Nama     string `json:"nama"`
	Kategori string `json:"kategori"`
	Harga    int    `json:"harga"`
}

func main() {
	// Get Menu
	menu, err := getMenu()
	if err != nil {
		fmt.Println("Error! ", err.Error())
		return
	}

	fmt.Println("Lihat Menu")
	for _, value := range menu {
		fmt.Println("ID: ", value.IDMenu, " Nama: ", value.Nama, " Harga: ", value.Harga, " Kategori: ", value.Kategori)
	}

	// Search Menu
	cari, err := cariMenu("goreng")
	if err != nil {
		fmt.Println("Error! ", err.Error())
		return
	}

	fmt.Println("\nCari Menu")

	for _, v := range cari {
		fmt.Println("ID: ", v.IDMenu, " Nama: ", v.Nama, " Harga: ", v.Harga, " Kategori: ", v.Kategori)
	}

	// Edit menu
	fmt.Println("\nEdit Menu")
	jsonStr := `{"idMenu" : 1, "nama":"Ikan Bakar", "kategori":"Makanan","harga":15000}`
	msg, err := editMenu(jsonStr)
	if err != nil {
		fmt.Println("Error! ", err.Error())
	} else {
		fmt.Println(msg)
	}

	// Hapus Menu
	fmt.Println("\nHapus Menu")
	jsonStr = `{"idMenu":5}`
	msg, err = hapusMenu(jsonStr)
	if err != nil {
		fmt.Println("Error! ", err.Error())
	} else {
		fmt.Println(msg)
	}
}

func getMenu() ([]Menu, error) {
	var err error
	var client = &http.Client{}
	var data []Menu

	request, err := http.NewRequest("GET", baseURL+"/lihat-menu", nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func cariMenu(namamenu string) ([]Menu, error) {
	var err error
	var client = &http.Client{}
	var data []Menu

	request, err := http.NewRequest("GET", baseURL+"/cari-menu?nama="+namamenu, nil)
	if err != nil {
		return data, err
	}

	request.Header.Set("Content-Type", "x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		return data, err
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return data, err
	}

	return data, nil

}

func editMenu(jsonStr string) (string, error) {
	var err error
	var client = &http.Client{}
	var message map[string]string

	payload := bytes.NewBuffer([]byte(jsonStr))

	request, err := http.NewRequest("POST", baseURL+"/edit-menu", payload)
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&message)
	if err != nil {
		return "", err
	}

	return message["message"], nil

}

func hapusMenu(jsonStr string) (string, error) {
	var err error
	var client = &http.Client{}
	var message map[string]string

	payload := bytes.NewBuffer([]byte(jsonStr))

	request, err := http.NewRequest("POST", baseURL+"/hapus-menu", payload)
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&message)
	if err != nil {
		return "", err
	}

	return message["message"], nil

}
