package main

// import (
// 	"encoding/json"
// 	"net/http"
// 	"strconv"
//
// 	"github.com/gorilla/mux"
// 	"github.com/mana-sg/vcs/internal/repository"
// )

// func FetchRepos(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	userId := vars["userId"]
//
// 	repos, err := repository.GetAllRepositoriesForUser(VarDb, userId)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
//
// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(repos); err != nil {
// 		http.Error(w, "Error encoding reponse", http.StatusInternalServerError)
// 	}
// }
//
// func FetchCommits(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	repoId := vars["repoId"]
//
// 	commits, err := repository.GetAllCommitsForRepo(VarDb, repoId)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 	}
//
// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(commits); err != nil {
// 		http.Error(w, "Error encoding response", http.StatusInternalServerError)
// 	}
// }
//
// func FetchFiles(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	commitIdstr := vars["commitId"]
//
// 	commitId, err := strconv.Atoi(commitIdstr)
// 	if err != nil {
// 		http.Error(w, "Invalid commit ID", http.StatusBadRequest)
// 		return
// 	}
//
// 	files, err := repository.GetAllFilesForCommit(VarDb, commitId)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 	}
//
// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(files); err != nil {
// 		http.Error(w, "Error encoding response", http.StatusInternalServerError)
// 	}
//
// }
