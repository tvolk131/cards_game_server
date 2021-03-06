package user

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// User .
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetByID fetches a user from the database with a given ID
func GetByID(id int, db *sql.DB) (User, error) {
	rows, err := db.Query("SELECT name, email FROM users WHERE id = " + strconv.Itoa(id))
	if err != nil {
		return User{}, err
	}
	var name string
	var email string
	rows.Next()
	if err := rows.Scan(&name, &email); err != nil {
		log.Fatal(err)
	}
	return User{ID: id, Name: name, Email: email}, nil
}

// GetByRequest gets the user that sent the HTTP request
func GetByRequest(r *http.Request, db *sql.DB) (User, error) {
	cookie, e := r.Cookie("connect.sid")
	if e != nil {
		return User{}, e
	}
	return GetByCookie(cookie.Value, db)
}

// GetByCookie fetches a user from the database associated with a given cookie
func GetByCookie(c string, db *sql.DB) (User, error) {
	uid, e := getIDByCookie(c, db)
	if e != nil {
		return User{}, e
	}
	return GetByID(uid, db)
}

func getIDByCookie(c string, db *sql.DB) (int, error) {
	if len(c) > 32 {
		c = c[4:36]
	}
	rows, err := db.Query(`SELECT data FROM "Sessions" WHERE sid = '` + c + `'`)
	if err != nil {
		fmt.Println(err)
		return -1, errors.New("Cookie is not valid")
	}
	var data string
	if rows.Next() {
		rows.Scan(&data)
		return parseUserID(data), nil
	}
	return -1, errors.New("Cookie is not valid")
}

func parseUserID(data string) int {
	split := strings.Split(data, `}`)
	data = split[len(split)-3]
	data = strings.Split(data, `"user":`)[1]
	i, _ := strconv.ParseInt(data, 10, 64)
	return int(i)
}
