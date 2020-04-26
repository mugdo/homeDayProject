package backEnd

import (
	"database/sql"
	"net/http"
	"strings"
	"time"
)

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	session, _ := store.Get(r, "mysession")

	if session.Values["isLogin"] == true {
		http.Redirect(w, r, "/redirect", http.StatusSeeOther)
	} else {
		Info = map[string]interface{}{
			"PageTitle": "Registration",
		}
		tpl.ExecuteTemplate(w, "register.gohtml", Info)
	}
}

func DoRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	fullName := strings.TrimSpace(r.FormValue("fullName"))
	email := strings.TrimSpace(r.FormValue("email"))
	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))
	password, _ = hashPassword(password)

	DB := dbConn()

	//do regitration
	CurrentTime := time.Now()
	token := generateToken()
	tokenSent := time.Now().Unix()

	insertQuery, err := DB.Prepare("INSERT INTO user(fullName, email, username, password, createdAt, isVerified, token, tokenSent) VALUES(?,?,?,?,?,?,?,?)")
	checkErr(err)
	insertQuery.Exec(fullName, email, username, password, CurrentTime, 0, token, tokenSent)
	println("Registration Done")

	link := "http://localhost:8080/verify-email/token=" + token
	sendMail(email, link)

	lastPage = "RegistrationDone"
	http.Redirect(w, r, "/login", http.StatusSeeOther)

	defer DB.Close()
}
func EmailVerifiation(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	runes := []rune(path)

	need := "="
	index := strings.Index(path, need)

	if index == -1 {
		errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
	} else { //url is "/verify-email/token=" something like this
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
}
