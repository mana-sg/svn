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
	"github.com/mana-sg/vcs/internal/utils"
	"github.com/mana-sg/vcs/pkg/models"
	"github.com/mana-sg/vcs/pkg/types"
	"github.com/rs/cors"
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
	r.HandleFunc("/api/createRepo", CreateRepository).Methods("POST")
	r.HandleFunc("/api/fetchLatestCommitId/{repoId}", FetchLatestCommitId).Methods("GET")
	r.HandleFunc("/api/commit/{repoId}", CreateCommit).Methods("POST")
	r.HandleFunc("/api/numberOfRepos/{userId}", GetNumRepos).Methods("GET")

	handler := cors.Default().Handler(r)

	http.ListenAndServe(":6969", handler)
}

func GetNumRepos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	numberOfRepos, err := repository.GetNumberOfrepositories(VarDb, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(struct {
		ReposNum int `json:"numberOfRepos"`
	}{
		ReposNum: numberOfRepos,
	}); err != nil {
		http.Error(w, "Error encoding reponse", http.StatusInternalServerError)
	}
}

func FetchRepos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	repos, err := repository.GetAllRepositoriesForUser(VarDb, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(repos); err != nil {
		http.Error(w, "Error encoding reponse", http.StatusInternalServerError)
	}
}

func FetchCommits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoId := vars["repoId"]

	commits, err := repository.GetAllCommitsForRepo(VarDb, repoId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(commits); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func FetchLatestCommitId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoId := vars["repoId"]

	repoIdNum, err := strconv.Atoi(repoId)
	if err != nil {
		http.Error(w, "Ivalid repository id", http.StatusBadRequest)
		return
	}

	userId, _ := models.GetActiveUser()
	commitId, err := repository.GetLatestCommit(VarDb, int(userId), repoIdNum)
	fmt.Println(commitId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(struct {
		CommitId int `json:"commitId"`
	}{CommitId: commitId}); err != nil {
		http.Error(w, "Error encoding data", http.StatusInternalServerError)
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
	w.WriteHeader(http.StatusOK)
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
}

type RepoCreation struct {
	UserId   int    `json:"userId"`
	RepoName string `json:"repoName"`
}

func CreateRepository(w http.ResponseWriter, r *http.Request) {
	var repo RepoCreation

	if err := json.NewDecoder(r.Body).Decode(&repo); err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	err := repository.CreateRepo(VarDb, repo.RepoName, fmt.Sprintf("%d", repo.UserId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	repoId, _ := models.GetActiveRepo()

	response := struct {
		RepoId  uint   `json:"repoId"`
		Message string `json:"message"`
	}{
		RepoId:  repoId,
		Message: "Repository created succesfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

type Updates struct {
	Id           int    `json:"id"`
	Modification string `json:"modification"`
}

type CreateCommitStruct struct {
	Commit string           `json:"commit"`
	Files  []types.FileNode `json:"files"`
}

func CreateCommit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoId, err := strconv.Atoi(vars["repoId"])
	if err != nil {
		http.Error(w, "Invalid repository id", http.StatusBadRequest)
		return
	}
	var files []types.FileNode
	var commit string
	var createCommitVar CreateCommitStruct

	if err := json.NewDecoder(r.Body).Decode(&createCommitVar); err != nil {
		http.Error(w, "Input Invalid", http.StatusBadRequest)
		return
	}
	files = createCommitVar.Files
	commit = createCommitVar.Commit

	commitId, err := repository.CreateCommit(VarDb, commit, repoId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rootDirHash, err := utils.HashDirectoryContents(files)
	if err != nil {
		http.Error(w, "Error hashing directory", http.StatusInternalServerError)
		return
	}
	err = repository.AddTree(VarDb, string(rootDirHash), files, commitId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{
		Message: "Insertion Succesful",
	}); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
	}
}
