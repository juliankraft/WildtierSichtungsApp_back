package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Sichtung struct {
	SichtungenID   int     `json:"sichtungen_id"`
	UserName       string  `json:"user_name"`
	Tierart        string  `json:"tierart"`
	AnzahlMaennlich int    `json:"anzahl_maennlich"`
	AnzahlWeiblich  int    `json:"anzahl_weiblich"`
	AnzahlUnbekannt int    `json:"anzahl_unbekannt"`
	Bemerkung      string  `json:"sichtung_bemerkung"`
	SichtungDate   string  `json:"sichtung_date"`
	Location       struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"sichtung_location"`
}


func (app *App) getData(w http.ResponseWriter, r *http.Request) {
	query := `SELECT 
                s.sichtungen_id, 
                u.user_name, 
                t.tierart, 
                s.anzahl_maennlich, 
                s.anzahl_weiblich, 
                s.anzahl_unbekannt, 
                s.sichtung_bemerkung, 
                s.sichtung_date, 
                ST_X(s.sichtung_location) AS lng, 
                ST_Y(s.sichtung_location) AS lat
              FROM sichtungen s
              INNER JOIN users u ON s.user_id = u.user_id
              INNER JOIN tierarten t ON s.tierart_id = t.tierart_id`

	rows, err := app.DB.Query(query)
	if err != nil {
		fmt.Println("Error querying database")
		fmt.Println("Error details: ", err.Error())
	}
	defer rows.Close()

	var sichtungen []Sichtung
	for rows.Next() {
		var s Sichtung
		err := rows.Scan(
			&s.SichtungenID, 
			&s.UserName, 
			&s.Tierart, 
			&s.AnzahlMaennlich, 
			&s.AnzahlWeiblich, 
			&s.AnzahlUnbekannt, 
			&s.Bemerkung, 
			&s.SichtungDate, 
			&s.Location.Lng, 
			&s.Location.Lat)
			
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning row: %v", err), http.StatusInternalServerError)
			return
		}
		sichtungen = append(sichtungen, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sichtungen)
}