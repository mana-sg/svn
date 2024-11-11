package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
//
// 	"github.com/mana-sg/vcs/internal/user"
// )
//
// type UserSignUp struct {
// 	Name            string `json:"name"`
// 	Email           string `json:"email"`
// 	Password        string `json:"password"`
// 	ConfirmPassword string `json:"confirmPassword"`
// }
//
// type UserLogIn struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }
//
// func SignUpUser(w http.ResponseWriter, r *http.Request) {
// 	var varUser UserSignUp
// 	if err := json.NewDecoder(r.Body).Decode(&varUser); err != nil {
// 		http.Error(w, "Invalid Input", http.StatusBadRequest)
// 		return
// 	}
//
// 	err := user.CreateUser(VarDb, varUser.Name, varUser.Email, varUser.Password, varUser.ConfirmPassword)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
//
// 	w.WriteHeader(http.StatusCreated)
// 	fmt.Fprintln(w, "User registration succesfull")
// }
//
// func LoginUser(w http.ResponseWriter, r *http.Request) {
// 	var varUser UserLogIn
// 	if err := json.NewDecoder(r.Body).Decode(&varUser); err != nil {
// 		http.Error(w, "Invalid Input", http.StatusBadRequest)
// 		return
// 	}
//
// 	err := user.LogIn(VarDb, varUser.Email, varUser.Password)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
//
// 	w.WriteHeader(http.StatusCreated)
// 	fmt.Fprintln(w, "User log in succesfull")
// }
