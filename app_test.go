package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	err := a.Initialise(DBUser, DBPassword, "test")
	if err != nil {
		log.Fatal("Error occurred while initialising database")
	}
	createTable()
	m.Run()
}

func createTable() {
	createTableQuery := "create table if not exists products (id int not null auto_increment,name varchar(20) not null, quantity int, price float(10,7), primary key(id));"

	_, err := a.DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("delete from products")
	a.DB.Exec("alter table products auto_increment=1")
	log.Println("clear table")
}

func addProducts(name string, quantity int, price float32) {
	query := fmt.Sprintf("insert into products(name, quantity, price) values('%v',%v,%v)", name, quantity, price)
	_, err := a.DB.Exec(query)
	if err != nil {
		return
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProducts("pen", 12, 1.99)
	request, _ := http.NewRequest("GET", "/products/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)
}

func checkStatusCode(t *testing.T, expectedStatusCode int, actualStatusCode int) {
	if expectedStatusCode != actualStatusCode {
		t.Errorf("Expected status code: %v, Recieved status code: %v", expectedStatusCode, actualStatusCode)
	}
}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}

func TestCreateProduct(t *testing.T) {
	clearTable()
	var product = []byte(`{"name":"bed", "quantity":1, "price":32.99}`)
	request, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(product))
	request.Header.Set("Content-Type", "application/json")

	response := sendRequest(request)
	checkStatusCode(t, http.StatusCreated, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "bed" {
		t.Errorf("Expected name: %v, Recieved name: %v", "bed", m["name"])
	}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProducts("charger", 1, 1.99)

	request, _ := http.NewRequest("GET", "/products/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

	request, _ = http.NewRequest("DELETE", "/products/1", nil)
	response = sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

	request, _ = http.NewRequest("GET", "/products/1", nil)
	response = sendRequest(request)
	checkStatusCode(t, http.StatusNotFound, response.Code)
}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProducts("charger", 1, 1.99)

	request, _ := http.NewRequest("GET", "/products/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

	var oldValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &oldValue)

	var product = []byte(`{"name":"charger", "quantity":5, "price":2.99}`)
	request, _ = http.NewRequest("PUT", "/products/1", bytes.NewBuffer(product))
	request.Header.Set("Content-Type", "application/json")

	response = sendRequest(request)
	var newValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &newValue)

	if oldValue["id"] != newValue["id"] {
		t.Errorf("Expected ID: %v, Received ID: %v", newValue["id"], oldValue["id"])
	}
	if oldValue["name"] != newValue["name"] {
		t.Errorf("Expected name: %v, Received name: %v", newValue["id"], oldValue["id"])
	}
}
