package backEnd

import (
	"database/sql"
	"html"
	"net/http"
	"strings"
	"time"
)

func resetCommon(w http.ResponseWriter, r *http.Request, title string) {
	session, _ := store.Get(r, "mysession")
	Info["Username"] = session.Values["username"]
	Info["Password"] = session.Values["password"]
	Info["IsLogged"] = session.Values["isLogin"]
	Info["LastPage"] = lastPage
	Info["PageTitle"] = title

	tpl.ExecuteTemplate(w, "reset.gohtml", Info)
}
func Reset(w http.ResponseWriter, r *http.Request) { //calling from submit of 1.Pass reset or 2.Token reset
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	path := r.URL.Path

	if r.Method != "POST" {
		var title string
		if path == "/resetPassword" {
			title = "Reset Password"
			resetCommon(w, r, title)
		} else if path == "/resetToken" {
			title = "New Token Request"
			resetCommon(w, r, title)
		} else {
			errorPage(w, http.StatusNotFound)
		}
	} else if r.Method == "POST" {
		email := html.EscapeString(strings.TrimSpace(r.FormValue("email")))

		DB := dbConn()
		defer DB.Close()

		token := generateToken()
		tokenSent := time.Now().Unix()

		if path == "/resetPassword" { //Request for Password reset
			//cheching for email already exist or not in the resetpassword table in DB
			var ID int
			res := DB.QueryRow("SELECT id FROM resetpassword WHERE email=?", email).Scan(&ID)
			if res == sql.ErrNoRows { //Row not found
				//Inserting in Reset Table
				insertQuery, err := DB.Prepare("INSERT INTO resetpassword (email,token,tokenSent) VALUES(?,?,?)")
				checkErr(err)
				insertQuery.Exec(email, token, tokenSent)
			} else { //found a row
				//already exist the email //only updating now
				updateQuery, err := DB.Prepare("UPDATE resetpassword SET token=?,tokenSent=? WHERE email=?")
				checkErr(err)
				updateQuery.Exec(token, tokenSent, email)
			}
			link := "http://localhost:8080/passReset/token=" + token //sending link to mail
			sendMail(email, link, "password")

			popUpCause = "passwordRequest"
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else if path == "/resetToken" { //Request for Token reset
			//updating DB with new token
			insertQuery, err := DB.Prepare("UPDATE user SET token=?,tokenSent=? WHERE email=?")
			checkErr(err)
			insertQuery.Exec(token, tokenSent, email)

			link := "http://localhost:8080/verify-email/token=" + token
			sendMail(email, link, "email")

			popUpCause = "tokenRequest"
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			errorPage(w, http.StatusNotFound)
		}
	}
}
func PassReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	DB := dbConn()
	defer DB.Close()

	if r.Method != "POST" {
		path := r.URL.Path
		runes := []rune(path)
		need := "="
		index := strings.Index(path, need)
		token := string(runes[index+1:])

		var email string
		res := DB.QueryRow("SELECT email FROM resetpassword WHERE token=?", token).Scan(&email)
		if res == sql.ErrNoRows { //Row not found
			popUpCause = "passTokenInvalid"
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else { //found a row
			//checking for token expired or not
			var tokenSent int64
			_ = DB.QueryRow("SELECT tokenSent FROM resetpassword WHERE token=?", token).Scan(&tokenSent)
			tokenReceived := time.Now().Unix() //current time
			diff := tokenReceived - tokenSent

			if diff > (2 * 60 * 60) { //2 hours period
				popUpCause = "passTokenExpired"
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else {
				session, _ := store.Get(r, "mysession")

				Info["Username"] = session.Values["username"]
				Info["Password"] = session.Values["password"]
				Info["IsLogged"] = session.Values["isLogin"]
				Info["LastPage"] = lastPage
				Info["PageTitle"] = "Reset Password"
				Info["Token"] = token

				tpl.ExecuteTemplate(w, "passReset.gohtml", Info)
			}
		}
	} else if r.Method == "POST" {
		token := r.FormValue("token")
		password := strings.TrimSpace(r.FormValue("password"))
		password, _ = hashPassword(password) //hashing password

		var email string
		_ = DB.QueryRow("SELECT email FROM resetpassword WHERE token=?", token).Scan(&email)

		//updating user table with new password of perspective email/user
		query, err := DB.Prepare("UPDATE user SET password=? WHERE email=?")
		checkErr(err)
		query.Exec(password, email)

		popUpCause = "passwordReset"
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
