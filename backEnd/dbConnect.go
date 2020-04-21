package backEnd

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"regexp"
	"strings"
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
}
