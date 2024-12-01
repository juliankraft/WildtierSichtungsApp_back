package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Animal struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (app *App) animalsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := app.DB.Query("SELECT tierart_id,tierart FROM tierarten order by tierart")
	if err != nil {
		fmt.Println("Error querying database")
		fmt.Println("Error details: ", err.Error())
	}
	defer rows.Close()

	var animals []Animal

	for rows.Next() {
		var animal Animal
		err := rows.Scan(&animal.ID, &animal.Name)
		animals = append(animals, animal)
		if err != nil {
			fmt.Println("Error scanning rows")
			fmt.Println("Error details: ", err.Error())
		}
	}
	//json response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animals)
}

func (app *App) saveAnimalHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var dataset map[string]interface{}
    err := json.NewDecoder(r.Body).Decode(&dataset)
    if err != nil {
        http.Error(w, "Error decoding JSON", http.StatusBadRequest)
        return
    }

    fmt.Println("Received dataset:", dataset)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
