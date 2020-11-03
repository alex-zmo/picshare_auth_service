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
		admin INT DEFAULT 0,
		PRIMARY KEY (id)
	);`

	_, err := db.Exec(createUserStmt)
	if err != nil {
		fmt.Println(err)
	}
}


func createUserHistoryTable() {
	createUserStmt := `CREATE TABLE IF NOT EXISTS history (
		id INT AUTO_INCREMENT,
		user_id INT,
		username VARCHAR(32),
		email VARCHAR(1024),
		first_name VARCHAR(1024),
		last_name VARCHAR(1024),
		password VARCHAR(128),
		last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
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


func MatchUsernameOrEmailToPassword(usernameOrEmail, password string) (string, int) {
	selectUserStmt := `SELECT username, id FROM user WHERE (email = ? OR username = ?) AND password = ?;`

	selectUserResult, err := db.Query(selectUserStmt, usernameOrEmail, usernameOrEmail, password)
	if err != nil {
		fmt.Println(err)
		return "", -1
	}
	defer closeRows(selectUserResult)

	usernameOut := ""
	userIdOut := -1

	if selectUserResult.Next() {
		err := selectUserResult.Scan(&usernameOut, &userIdOut)
		if err != nil {
			fmt.Println(err)
			return "", -1
		}
	}
	return usernameOut, userIdOut
}

func GetUser(userId int) *model.User{
	getUserStmt := `SELECT username, email, first_name, last_name, created_at, admin 
	FROM user where id = ?`

	result, err := db.Query(getUserStmt, userId )
	if err != nil {
		fmt.Println(err)
	}
	defer closeRows(result)

	user := &model.User{}
	if result.Next() {
		result.Scan(&user.Username, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt, &user.Admin)
	}

	return user
}

func UpdateUserPassword(userId int, password, newpassword  string) string {
	updatePassStmt := `UPDATE user SET password = ?, updated_at = CURRENT_TIMESTAMP 
	WHERE id = ?;
;`
	insertUserHistoryStmt := `INSERT INTO history (
		user_id,
		username, 
		email,
		first_name, 
		last_name, 
		password
	) SELECT id, username, email, first_name, last_name, password 
		FROM user 
		WHERE id = ?;`

	if !MatchUserIdToPassword(userId, password) {
		return "Password Invalid"
	}

	if !CheckNotInRecentThree(userId, newpassword) {
		return "New password cannot be a recently used password"
	}
	_, err := db.Exec(insertUserHistoryStmt, userId)
	if err != nil {
		fmt.Println(err)
	}

	_, err1 := db.Exec(updatePassStmt, newpassword, userId)
	if err1 != nil {
		fmt.Println(err1)
		return "Update User Failed"
	}

	return ""
}

func CheckNotInRecentThree(userId int, password string) bool {
	selectPassStmt := `SELECT COUNT(password) FROM 
	(SELECT * FROM history WHERE user_id = ? ORDER BY last_updated DESC LIMIT 3) 
	AS TEMP
	WHERE password = ?;`

	selectPassResult, err := db.Query(selectPassStmt, userId, password)
	if err != nil {
		fmt.Println(err)
		return false
	}
	count := 0
	defer closeRows(selectPassResult)
	if selectPassResult.Next() {
		err := selectPassResult.Scan(&count)
		if err != nil {
			fmt.Println(err)
			return false
		}
	}
	return count == 0
}

func MatchUserIdToPassword(userId int, password string) bool {
	selectUserStmt := `SELECT COUNT(1) FROM user WHERE id = ? AND password = ?;`

	selectUserResult, err := db.Query(selectUserStmt, userId, password)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer closeRows(selectUserResult)
	count := 0
	if selectUserResult.Next() {
		err := selectUserResult.Scan(&count)
		if err != nil {
			fmt.Println(err)
			return false
		}
	}
	return count == 1
}

func RetrieveEmailUserId(userId int) string{
	retrieveEmailStmt := `SELECT email FROM user
	WHERE id = ?
	`
	res, err := db.Query(retrieveEmailStmt, userId)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer closeRows(res)
	email := ""
	if res.Next() {
		err := res.Scan(&email)
		if err != nil {
			fmt.Println(err)
		}
	}
	return email
}