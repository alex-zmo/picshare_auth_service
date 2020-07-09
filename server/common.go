package server

import (
	"fmt"
	"log"
	"net/http"
	"smart_photos/auth_service/database"
	"smart_photos/auth_service/model"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func InitServer() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func getURLParam(r *http.Request, paramName string) string {
	keys, seen := r.URL.Query()[paramName]

	if seen && len(keys) > 0 {
		return keys[0]
	}
	return ""
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	username := getURLParam(r, "username")
	email := getURLParam(r, "email")
	firstName := getURLParam(r, "first_name")
	lastName := getURLParam(r, "last_name")
	password := getURLParam(r, "password")

	user := &model.User{Username: username, Email:email, FirstName: firstName, LastName: lastName, Password: password}
	fmt.Println(user)
	err := database.InsertUser(user)
	if err == "" {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusConflict)
	}

	_, _ = fmt.Fprintf(w, err)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

}
