package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/mana-sg/vcs/internal/db"
)

var VarDb db.DbHandler

func main() {
	godotenv.Load("../.env")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	host := os.Getenv("DATABASE_HOST")
	dbname := os.Getenv("DATABASE_NAME")

	VarDb = db.DbHandler{}

	err := VarDb.ConfigDB(user, password, host, dbname)
	if err != nil {
		fmt.Errorf("Error connecting to database: %v", err)
	}

	err = VarDb.PrepareDb()
	if err != nil {
		fmt.Errorf("error preparing the database: %v", err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/api/login", LoginUser).Methods("POST")
	r.HandleFunc("/api/signup", SignUpUser).Methods("POST")
	r.HandleFunc("/api/fetchRepos/{userId}", FetchRepos).Methods("GET")
	r.HandleFunc("/api/fetchCommits/{repoId}", FetchCommits).Methods("GET")
	r.HandleFunc("/api/fetchFiles/{commitId}", FetchFiles).Methods("GET")

	http.ListenAndServe(":6969", r)
}
