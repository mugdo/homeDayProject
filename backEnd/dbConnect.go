package backEnd

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "ajudge"

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	checkErr(err)

	return db
}
func CheckDB(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	runes := []rune(path)
	need := "="
	index := strings.Index(path, need)
	colValue := string(runes[index+1:]) //username or email value

	var colName string
	matched, err := regexp.MatchString("username", path)
	checkErr(err)
	if matched {
		colName = "username"
	} else {
		colName = "email"
	}

	DB := dbConn()

	//checking for username/email already exist or not
	temp := ""
	res := DB.QueryRow("SELECT id FROM user WHERE "+colName+"=?", colValue).Scan(&temp)
	//checkErr(err)
	if res == sql.ErrNoRows {
		//username/email available (found no rows)
		b := []byte("false")
		w.Write(b)
	} else { //found a row
		b := []byte("true")
		w.Write(b)
	}
	defer DB.Close()
}
func EmailVerifiation(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	runes := []rune(path)

	need := "="
	index := strings.Index(path, need)
	token := string(runes[index+1:])

	DB := dbConn()

	//checking for token sent time
	var tokenSent int64
	res := DB.QueryRow("SELECT tokenSent FROM user WHERE token=?", token).Scan(&tokenSent)

	if res == sql.ErrNoRows {
		//token not found
		lastPage = "tokenInvalid"
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else { //found a row
		//checking for already verified or not
		var isVerified int
		_ = DB.QueryRow("SELECT isVerified FROM user WHERE token=?", token).Scan(&isVerified)
		if isVerified == 1 {
			//already verified
			lastPage = "tokenAlreadyVerified"
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else if isVerified == 0 {
			//checking for token expired or not
			tokenReceived := time.Now().Unix()
			diff := tokenReceived - tokenSent

			if diff > (2 * 60 * 60) { //2 hours period
				lastPage = "tokenExpired"
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else {
				updateQuery, err := DB.Prepare("UPDATE user SET isVerified=1,tokenSent=? WHERE token=?")
				checkErr(err)
				updateQuery.Exec(tokenReceived, token) //after verified tokenSent indicates the verified time

				lastPage = "tokenVerifiedNow"
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		}
	}

	defer DB.Close()
}
