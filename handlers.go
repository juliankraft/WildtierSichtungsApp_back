package main

import (
	"crypto/sha1"
    "encoding/hex"
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
			sichtung_bemerkung
		) VALUES (?, ?, ?, ?, ?, ST_GeomFromText(?), ?)
	`

	// Execute the query
	_, err = app.DB.Exec(query,
		dataset["animal_id"],
		dataset["count_male"],
		dataset["count_female"],
		dataset["count_unknown"],
		sichtungDateString,
		location,
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

func (app *App) getUsernamesHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Query to get all usernames
    query := `SELECT user_name FROM users`
    rows, err := app.DB.Query(query)
    if err != nil {
        http.Error(w, "Error querying database", http.StatusInternalServerError)
        fmt.Println("Error querying database:", err)
        return
    }
    defer rows.Close()

    var usernames []string
    for rows.Next() {
        var username string
        if err := rows.Scan(&username); err != nil {
            http.Error(w, "Error scanning rows", http.StatusInternalServerError)
            fmt.Println("Error scanning rows:", err)
            return
        }
        usernames = append(usernames, username)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(usernames)
}

type User struct {
    Username  string `json:"user_name"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string `json:"email"`
    Password  string `json:"pwd"`
}

func (app *App) createUserHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Error decoding JSON", http.StatusBadRequest)
        return
    }

    // Hash the password using SHA-1
    hasher := sha1.New()
    hasher.Write([]byte(user.Password))
    hashedPassword := hex.EncodeToString(hasher.Sum(nil))

    // Insert the user into the database
    query := `INSERT INTO users (user_name, first_name, last_name, email, pwd) VALUES (?, ?, ?, ?, ?)`
    _, err = app.DB.Exec(query, user.Username, user.FirstName, user.LastName, user.Email, hashedPassword)
    if err != nil {
		errormessage := fmt.Sprintf("Error inserting user into database: %s", err)
        http.Error(w, errormessage, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}