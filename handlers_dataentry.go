package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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
	fmt.Println("saveAnimalHandler called") // Debugging statement

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Print the request method and headers
	fmt.Println("Request Method:", r.Method)
	fmt.Println("Request Headers:", r.Header)

	// Read the request body
	var dataset map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&dataset)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		fmt.Println("Error decoding JSON:", err) // Debugging statement
		return
	}

	// Print the received dataset
	fmt.Println("Received dataset:", dataset)

	// Convert latitude and longitude to WKT format for GEOMETRY type
	location := fmt.Sprintf("POINT(%f %f)", dataset["longitude"].(float64), dataset["latitude"].(float64))

	// Parse the date string to time.Time
	sichtungDate, err := time.Parse("2006-01-02T15:04:05", dataset["date"].(string))
	if err != nil {
		http.Error(w, "Error parsing date", http.StatusBadRequest)
		fmt.Println("Error parsing date:", err) // Debugging statement
		return
	}
	fmt.Println(sichtungDate)
	sichtungDateString := sichtungDate.Format("2006-01-02 15:04:05")
	fmt.Println(sichtungDateString)

	// Prepare the query to insert data into the database
	query := `
		INSERT INTO sichtungen (
			tierart_id, 
			anzahl_maennlich, 
			anzahl_weiblich, 
			anzahl_unbekannt, 
			sichtung_date, 
			sichtung_location,
			phone_location,
			sichtung_bemerkung
		) VALUES (?, ?, ?, ?, ?, ST_GeomFromText(?), ?, ?)
	`

	// Execute the query
	_, err = app.DB.Exec(query,
		dataset["animal_id"],
		dataset["count_male"],
		dataset["count_female"],
		dataset["count_unknown"],
		sichtungDateString,
		location,
		dataset["phone_location"],
		dataset["notes"],
	)
	if err != nil {
		http.Error(w, "Error inserting data into database", http.StatusInternalServerError)
		fmt.Println("Error inserting data into database:", err) // Debugging statement
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
