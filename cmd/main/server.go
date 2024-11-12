package main

import (
	"fmt"
	"net/http"
	"os"

	"encoding/json"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/mana-sg/vcs/internal/db"
	"github.com/mana-sg/vcs/internal/repository"
	"github.com/mana-sg/vcs/internal/user"
	"github.com/mana-sg/vcs/pkg/models"
	"github.com/rs/cors"
)

var VarDb db.DbHandler

func main() {
	godotenv.Load("../../.env")
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

	handler := cors.Default().Handler(r)

	http.ListenAndServe(":6969", handler)
}

func FetchRepos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	repos, err := repository.GetAllRepositoriesForUser(VarDb, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(repos)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(repos); err != nil {
		http.Error(w, "Error encoding reponse", http.StatusInternalServerError)
	}
}

func FetchCommits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoId := vars["repoId"]
	fmt.Println(repoId)

	commits, err := repository.GetAllCommitsForRepo(VarDb, repoId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(commits); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func FetchFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commitIdstr := vars["commitId"]

	commitId, err := strconv.Atoi(commitIdstr)
	if err != nil {
		http.Error(w, "Invalid commit ID", http.StatusBadRequest)
		return
	}

	files, err := repository.GetAllFilesForCommit(VarDb, commitId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(files); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}

}

type UserSignUp struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type UserLogIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignUpUser(w http.ResponseWriter, r *http.Request) {
	var varUser UserSignUp
	if err := json.NewDecoder(r.Body).Decode(&varUser); err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	err := user.CreateUser(VarDb, varUser.Name, varUser.Email, varUser.Password, varUser.ConfirmPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}

	userId, _ := models.GetActiveUser()
	response := struct {
		Message string `json:"message"`
		UserID  uint   `json:"userId"`
	}{
		Message: "User registration successful",
		UserID:  userId,
	}

	// Set the content-type to application/json
	w.Header().Set("Content-Type", "application/json")
	// Set the status code to 201 Created
	w.WriteHeader(http.StatusCreated)
	// Send the response as JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	fmt.Println(w)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var varUser UserLogIn
	if err := json.NewDecoder(r.Body).Decode(&varUser); err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	err := user.LogIn(VarDb, varUser.Email, varUser.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userId, _ := models.GetActiveUser()
	response := struct {
		Message string `json:"message"`
		UserID  uint   `json:"userId"`
	}{
		Message: "User login successful",
		UserID:  userId,
	}

	// Set the content-type to application/json
	w.Header().Set("Content-Type", "application/json")
	// Set the status code to 201 Created
	w.WriteHeader(http.StatusCreated)
	// Send the response as JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	fmt.Println(w)
}