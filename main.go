package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Customer struct {
	Id        uuid.UUID `json:"id,omitempty"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Address   *Address  `json:"address,omitempty"`
}

type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

var customers = make([]Customer, 0)

func GetCustomers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(customers)
}

func GetCustomerById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range customers {
		id, err := uuid.Parse(params["id"])
		if err != nil {
			fmt.Println("Erro ao converter string para UUID", err)
		}

		if item.Id == id {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Customer{})
}

func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer Customer
	_ = json.NewDecoder(r.Body).Decode(&customer)
	customer.Id = uuid.New()
	customers = append(customers, customer)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range customers {
		id, err := uuid.Parse(params["id"])
		if err != nil {
			fmt.Println("Erro ao converter string para UUID", err)
		}

		if item.Id == id {
			customers = append(customers[:index], customers[index+1:]...)
			break
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	router := mux.NewRouter()

	customers = append(customers, Customer{Id: uuid.New(), FirstName: "John", LastName: "Doe",
		Address: &Address{City: "City 1", State: "State 1"}})

	customers = append(customers, Customer{Id: uuid.New(), FirstName: "Koko", LastName: "Doe",
		Address: &Address{City: "City 2", State: "State 2"}})

	router.HandleFunc("/customers", GetCustomers).Methods("GET")
	router.HandleFunc("/customers/{id}", GetCustomerById).Methods("GET")
	router.HandleFunc("/customers", CreateCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", DeleteCustomer).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
