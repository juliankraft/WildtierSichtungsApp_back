package main

import (
	"crypto/ecdsa"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	//jwt-go
	"github.com/dgrijalva/jwt-go"
)

func hashPassword(password string) string {
	hasher := sha1.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
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
	ID        int    `json:"user_id"`
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

	// Insert the user into the database
	query := `INSERT INTO users (user_name, first_name, last_name, email, pwd) VALUES (?, ?, ?, ?, ?)`
	_, err = app.DB.Exec(query, user.Username, user.FirstName, user.LastName, user.Email, hashPassword(user.Password))
	if err != nil {
		errormessage := fmt.Sprintf("Error inserting user into database: %s", err)
		http.Error(w, errormessage, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

var (
	key *ecdsa.PrivateKey
	t   *jwt.Token
	s   string
)

func (app *App) createNewJwtToken(u User) string {
	key = app.key
	t = jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"iss":      "tokenprovider",
			"sub":      u.ID,
			"username": u.Username,
		})
	s, err := t.SignedString(key)
	if err != nil {
		fmt.Println(err)
	}
	return s
}

func (app *App) loginHandler(w http.ResponseWriter, r *http.Request) {
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

	// Query to check if the user exists
	query := `SELECT user_id,user_name FROM users WHERE user_name = ? AND pwd = ?`
	row := app.DB.QueryRow(query, user.Username, hashPassword(user.Password))

	var id int
	var username string
	err = row.Scan(&id, &username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	user.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success",
		"token":     app.createNewJwtToken(user),
		"user_name": user.Username,
		"user_id":   fmt.Sprintf("%d", user.ID),
	})
}
