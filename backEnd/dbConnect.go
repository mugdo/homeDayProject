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
	
	var colName string
	matched, err := regexp.MatchString("/check/username=", path)
	checkErr(err)
	if matched {
		colName = "username"
	} else {
		matched, err = regexp.MatchString("/check/email=", path)
		checkErr(err)
		if matched {
		colName = "email"
		} else {
			errorPage(w, http.StatusNotFound) //http.StatusNotFound = 404
			return
		}
	}

	runes := []rune(path)
	need := "="
	index := strings.Index(path, need)
	colValue := string(runes[index+1:]) //username or email value

	DB := dbConn()
	defer DB.Close()

	//checking for username/email already exist or not
	id := 0
	res := DB.QueryRow("SELECT id FROM user WHERE "+colName+"=?", colValue).Scan(&id)

	if res == sql.ErrNoRows { //found no rows (username/email available)
		b := []byte("false")
		w.Write(b)
	} else { //found a row
		b := []byte("true")
		w.Write(b)
	}
}
func isAccVerifed(r *http.Request) bool {
	session, _ := store.Get(r, "mysession")
	username := session.Values["username"]

	DB := dbConn()
	defer DB.Close()

	var res bool
	var isVerified int
	_ = DB.QueryRow("SELECT isVerified FROM user WHERE username=?", username).Scan(&isVerified)
	if isVerified == 1 {
		res = true
	} else if isVerified == 0 {
		res = false
	}

	return res
}
