package main

// webapi that queries mariadb
// talks to an ionic frontend
import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type App struct {
	DB *sql.DB
}

func main() {
	//load form .env variables
	godotenv.Load()
	app := App{}
	//db connection setup
	dsn := os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PASSWORD") + "@tcp(" + os.Getenv("MYSQL_HOST") + ":3306)/" + os.Getenv("MYSQL_DATABASE") + "?charset=utf8mb4&parseTime=True&loc=Local"
	con, err := sql.Open("mysql", dsn)
	app.DB = con
	if err != nil {
		fmt.Println("Error connecting to database")
		fmt.Println("Error details: ", err.Error())
	}
	//test query

	//webserver
	fmt.Println("Starting server on port 8089")

	http.HandleFunc("/", app.indexHandler)
	
	http.HandleFunc("/api/v1/usernames", app.getUsernamesHandler)
	http.HandleFunc("/api/v1/createUser", app.createUserHandler)
	
	http.HandleFunc("/api/v1/animals", app.animalsHandler)
	http.HandleFunc("/api/v1/saveAnimal", app.saveAnimalHandler)

	// Apply CORS middleware
    http.ListenAndServe(":8089", corsMiddleware(http.DefaultServeMux))
}

func (app *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Wildtier API v1.0")
}
