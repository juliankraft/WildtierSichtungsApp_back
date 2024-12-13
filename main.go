package main

// webapi that queries mariadb
// talks to an ionic frontend
import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type App struct {
	DB  *sql.DB
	key string
}

func main() {
	//load form .env variables
	godotenv.Load()
	app := App{}
	//key := os.Getenv("KEY")
	//key = &ecdsa.PrivateKey{}
	// privatekey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	// if err != nil {
	// 	panic(err)
	// }
	// app.key = privatekey
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

	//open routes
	http.HandleFunc("/api/v1/usernames", app.getUsernamesHandler)
	http.HandleFunc("/api/v1/createUser", app.createUserHandler)
	http.HandleFunc("/api/v1/login", app.loginHandler)

	//authorized routes
	http.HandleFunc("/api/v1/animals", app.animalsHandler)
	http.HandleFunc("/api/v1/saveAnimal", app.saveAnimalHandler)

	http.HandleFunc("/api/v1/getData", app.getData)

	//serve static files in /app (the ionic app) and /static (the website)
	http.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("app"))))
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("static"))))
	// Apply CORS middleware
	err = http.ListenAndServe(":8089", corsMiddleware(http.DefaultServeMux))
	if err != nil {
		log.Println(err.Error())
	}
}
