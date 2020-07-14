package database

import (
	"fmt"
	"github.com/gmo-personal/picshare_auth_service/model"
)

func createUserTable() {
	createUserStmt := `CREATE TABLE IF NOT EXISTS user (
		id INT AUTO_INCREMENT,
		username VARCHAR(32),
		email VARCHAR(1024),
		first_name VARCHAR(1024),
		last_name VARCHAR(1024),
		password VARCHAR(128),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	);`

	_, err := db.Exec(createUserStmt)
	if err != nil {
		fmt.Println(err)
	}
}

func InsertUser(userInfo *model.User) string {
	insertUserStmt := `INSERT INTO user (
		username, 
		email,
		first_name, 
		last_name, 
		password
	) VALUES (?, ?, ?, ?, ?);`

	if checkEmailExists(userInfo.Username) {
		return "Email already exists"
	}
	if checkUsernameExists(userInfo.Username) {
		return "Username already exists"
	}

	_, err := db.Exec(insertUserStmt, userInfo.Username, userInfo.Email, userInfo.FirstName, userInfo.LastName, userInfo.Password)
	if err != nil {
		fmt.Println(err)
	}
	return ""
}

func checkUsernameExists(username string) bool {
	selectUserStmt := `SELECT * FROM user WHERE username = ?;`

	selectUserResult, err := db.Query(selectUserStmt, username)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer closeRows(selectUserResult)

	if selectUserResult.Next() {
		return true
	}

	return false
}

func checkEmailExists(email string) bool {
	selectUserStmt := `SELECT * FROM user WHERE email = ?;`

	selectUserResult, err := db.Query(selectUserStmt, email)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer closeRows(selectUserResult)

	if selectUserResult.Next() {
		return true
	}

	return false
}


func MatchUsernameOrEmailToPassword(usernameOrEmail, password string) string {
	selectUserStmt := `SELECT username FROM user WHERE (email = ? OR username = ?) AND password = ?;`

	selectUserResult, err := db.Query(selectUserStmt, usernameOrEmail, usernameOrEmail, password)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer closeRows(selectUserResult)

	usernameOut := ""

	if selectUserResult.Next() {
		err := selectUserResult.Scan(&usernameOut)
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}
	return usernameOut
}