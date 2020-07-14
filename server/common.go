package server

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gmo-personal/picshare_auth_service/database"
	"github.com/gmo-personal/picshare_auth_service/model"
	"log"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func InitServer() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/signup/", signupHandler)
	http.HandleFunc("/login/", loginHandler)

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
	enableCors(&w)
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
	enableCors(&w)
	usernameEmail := getURLParam(r, "usernameEmail")
	password := getURLParam(r, "password")

	username := database.MatchUsernameOrEmailToPassword(usernameEmail, password)
	token := CreateToken(username)

	_, _ = fmt.Fprintf(w, username + " TOKEN: " + token)
}


func CreateToken(username string) string {
	var err error
	//Creating Access Token
	secret := "jdnfksdmfksd"
	atClaims := jwt.MapClaims{}
	atClaims["username"] = username
	atClaims["expire"] = time.Now().Add(time.Minute * 360).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return ""
	}
	return token
}