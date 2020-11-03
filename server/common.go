package server

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gmo-personal/picshare_auth_service/database"
	"github.com/gmo-personal/picshare_auth_service/model"
	"log"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Authorization")

}

func InitServer() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/signup/", signupHandler)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/validate/", validateTokenHandler)
	http.HandleFunc("/created/", createdAtHandler)
	http.HandleFunc("/username/", usernameHandler)
	http.HandleFunc("/change/", changePasswordHandler)
	http.HandleFunc("/account/", accountHandler)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func accountHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	userId := getUserIdAuth(w, r)
	if userId == -1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if userId == -2 {
		return
	}
	user := database.GetUser(userId)
	username := user.Username
	email := user.Email
	firstName := user.FirstName
	lastName := user.LastName
	admin := user.Admin

	_, _ = fmt.Fprintf(w, username + " " + email + " " + firstName + " " + lastName + " " + strconv.Itoa(admin))
}

func changePasswordHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	password := getURLParam(r, "password")
	newpassword := getURLParam(r, "newpassword")

	userId := getUserIdAuth(w, r)
	if userId == -1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if userId == -2 {
		return
	}

	err := database.UpdateUserPassword(userId, password, newpassword)
	if err == "" {
		w.WriteHeader(http.StatusCreated)
		email := database.RetrieveEmailUserId(userId)
		if len(email) != 0 {
			message := []byte("To:" + email +"\r\n" +
				"Subject: PicShare password changed successful!\r\n" +
				"\r\n" +
				"Your PicShare password has been changed!\r\n\n"  )
			SendEmail(email, message)
		}
	} else {
		w.WriteHeader(http.StatusConflict)
	}

	_, _ = fmt.Fprintf(w, err)
}

func getUserIdAuth(w http.ResponseWriter, r *http.Request) int{
	enableCors(&w)
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) == 1 {
		fmt.Println("empty")
		return - 2
	}
	reqToken = splitToken[1]

	token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		return []byte("jdnfksdmfksd"), nil
	})
	if err != nil || token.Valid == false || token.Claims.(jwt.MapClaims)["expire"].(float64) < float64(time.Now().Unix()) {
		return -1
	}
	userIdString := fmt.Sprintf("%v", token.Claims.(jwt.MapClaims)["userId"])

	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		fmt.Println(err)
	}
	return userId
}
func createdAtHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	userId := getUserIdAuth(w,r)
	if userId == -1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if userId == -2 {
		return
	}
    user := database.GetUser(userId)
    fmt.Fprintf(w, "%v", user.CreatedAt)
}

func usernameHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	userId := getUserIdAuth(w,r)
	if userId == -1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if userId == -2 {
		return
	}
	user := database.GetUser(userId)
	fmt.Fprintf(w, "%v", user.Username)
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
	err := database.InsertUser(user)
	if err == "" {
		w.WriteHeader(http.StatusCreated)
		message := []byte("To:" + email +"\r\n" +
			"Subject: Welcome to PicShare!\r\n" +
			"\r\n" + "Welcome " + firstName + " " + lastName + ". "+
			"Your Picshare account has been successfully created!\r\n\n" + "Username: " + username+"\r\n" )
		SendEmail(email, message)
	} else {
		w.WriteHeader(http.StatusConflict)
	}

	_, _ = fmt.Fprintf(w, err)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	usernameEmail := getURLParam(r, "usernameEmail")
	password := getURLParam(r, "password")

	username, userId := database.MatchUsernameOrEmailToPassword(usernameEmail, password)

	token := ""
	if userId != -1 {
		token = CreateToken(username, userId)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
	_, _ = fmt.Fprintf(w, username + " TOKEN: " + token)
}


func CreateToken(username string, userId int) string {
	var err error
	//Creating Access Token
	secret := "jdnfksdmfksd"
	atClaims := jwt.MapClaims{}
	atClaims["username"] = username
	atClaims["userId"] = userId
	atClaims["expire"] = time.Now().Add(time.Minute * 360).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return ""
	}
	return token
}

func validateTokenHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	userId := getUserIdAuth(w, r)
	if userId == -1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	_, _ = fmt.Fprintf(w, "%v",userId)
}


func SendEmail(receiver string, message []byte) {
	fmt.Println("Sending Email to: " + receiver , "message: "+ string(message))
	// Sender data.
	from := "resforge.dev@gmail.com"
	password := "!Andy***951"

	// Receiver email address.
	to := []string{
		receiver,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}