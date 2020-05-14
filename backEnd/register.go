package backEnd

import (
	"database/sql"
	"html"
	"net/http"
	"strings"
	"time"
)

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	if r.Method != "POST" {
		session, _ := store.Get(r, "mysession")
		if session.Values["isLogin"] == true {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			Info["PageTitle"] = "Registration"

			tpl.ExecuteTemplate(w, "register.gohtml", Info)
		}
	} else if r.Method == "POST" {
		fullName := html.EscapeString(strings.TrimSpace(r.FormValue("fullName")))
		email := html.EscapeString(strings.TrimSpace(r.FormValue("email")))
		username := html.EscapeString(strings.TrimSpace(r.FormValue("username")))
		password := html.EscapeString(r.FormValue("password"))
		password, _ = hashPassword(password)

		DB := dbConn()
		defer DB.Close()

		//do regitration
		CurrentTime := time.Now() //current time format like this: 2009-11-10 23:00:00 +0000 UTC m=+0.000000001
		token := generateToken()
		tokenSent := time.Now().Unix() //current time format like this: 1257894000 (time passed since 1970 in seconds)

		insertQuery, err := DB.Prepare("INSERT INTO user(fullName, email, username, password, createdAt, isVerified, token, tokenSent) VALUES(?,?,?,?,?,?,?,?)")
		checkErr(err)
		insertQuery.Exec(fullName, email, username, password, CurrentTime, 0, token, tokenSent)

		link := "http://localhost:8080/verify-email/token=" + token
		sendMail(email, link, "email")

		popUpCause = "registrationDone"
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
func EmailVerifiation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	path := r.URL.Path
	runes := []rune(path)
	need := "="
	index := strings.Index(path, need)
	token := string(runes[index+1:])
	
	DB := dbConn()
	defer DB.Close()

	//checking for token sent time
	var tokenSent int64
	res := DB.QueryRow("SELECT tokenSent FROM user WHERE token=?", token).Scan(&tokenSent)

	if res == sql.ErrNoRows { //no row found (token not found)
		popUpCause = "tokenInvalid"
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else { //found a row
		//checking for already verified or not
		var isVerified int
		_ = DB.QueryRow("SELECT isVerified FROM user WHERE token=?", token).Scan(&isVerified)

		if isVerified == 1 { //already verified
			popUpCause = "tokenAlreadyVerified"
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else if isVerified == 0 {
			//checking for token expired or not
			tokenReceived := time.Now().Unix()
			diff := tokenReceived - tokenSent

			if diff > (2 * 60 * 60) { //2 hours period
				popUpCause = "tokenExpired"
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else {
				updateQuery, err := DB.Prepare("UPDATE user SET isVerified=1,tokenSent=? WHERE token=?")
				checkErr(err)
				updateQuery.Exec(tokenReceived, token) //after verified, tokenSent indicates the time of verification

				popUpCause = "tokenVerifiedNow"
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		}
	}
}
