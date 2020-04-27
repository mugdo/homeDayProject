package backEnd

import (
	"database/sql"
	"net/http"
	"strings"
	"time"
)

func resetCommon(w http.ResponseWriter, r *http.Request, title, action string) {
	session, _ := store.Get(r, "mysession")
	Info = map[string]interface{}{
		"Username":  session.Values["username"],
		"Password":  session.Values["password"],
		"IsLogged":  session.Values["isLogin"],
		"LastPage":  lastPage,
		"PageTitle": title,
		"Action":    action,
	}
	tpl.ExecuteTemplate(w, "reset.gohtml", Info)
}
func Reset(w http.ResponseWriter, r *http.Request) { //calling from two diff pages. 1.Request token 2.Forgot password
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	path := r.URL.Path

	var action, title string
	if path == "/resetPassword" {
		title = "Reset Password"
		action = "/DoResetP"
		resetCommon(w, r, title, action)
	} else if path == "/resetToken" {
		title = "New Token Request"
		action = "/DoResetT"
		resetCommon(w, r, title, action)
	} else {
		errorPage(w, http.StatusNotFound)
	}
}
func DoReset(w http.ResponseWriter, r *http.Request) { //calling from submit of 1.Pass reset or 2.Token reset
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	path := r.URL.Path
	email := strings.TrimSpace(r.FormValue("email"))
	DB := dbConn()
	token := generateToken()
	tokenSent := time.Now().Unix()

	if path == "/DoResetP" { //Request for Password reset
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
		sendMail(email, link)

		lastPage = "passwordRequest"
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else if path == "/DoResetT" { //Request for Token reset
		//updating DB with new token
		insertQuery, err := DB.Prepare("UPDATE user SET token=?,tokenSent=? WHERE email=?")
		checkErr(err)
		insertQuery.Exec(token, tokenSent, email)

		link := "http://localhost:8080/verify-email/token=" + token
		sendMail(email, link)

		lastPage = "tokenRequest"
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	defer DB.Close()
}
func PassReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	path := r.URL.Path
	runes := []rune(path)

	need := "token="
	index := strings.Index(path, need)

	if index == -1 {
		errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
	} else { //url is "/passReset/token=" something like this
		token := string(runes[index+6:]) //index+0 = t, on 'token='
		DB := dbConn()

		var email string
		res := DB.QueryRow("SELECT email FROM resetpassword WHERE token=?", token).Scan(&email)
		if res == sql.ErrNoRows { //Row not found
			lastPage = "passTokenInvalid"
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else { //found a row
			//checking for token expired or not
			var tokenSent int64
			_ = DB.QueryRow("SELECT tokenSent FROM resetpassword WHERE token=?", token).Scan(&tokenSent)
			tokenReceived := time.Now().Unix() //current time
			diff := tokenReceived - tokenSent

			if diff > (2 * 60 * 60) { //2 hours period
				lastPage = "passTokenExpired"
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else {

				session, _ := store.Get(r, "mysession")

				Info = map[string]interface{}{
					"Username":  session.Values["username"],
					"Password":  session.Values["password"],
					"IsLogged":  session.Values["isLogin"],
					"LastPage":  lastPage,
					"PageTitle": "Reset Password",
					"Token":     token,
				}
				tpl.ExecuteTemplate(w, "passReset.gohtml", Info)
			}
		}
		defer DB.Close()
	}
}
func DoPassReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	token := r.FormValue("token")
	password := strings.TrimSpace(r.FormValue("password"))
	password, _ = hashPassword(password) //hashing password

	DB := dbConn()

	var email string
	_ = DB.QueryRow("SELECT email FROM resetpassword WHERE token=?", token).Scan(&email)

	//updating user table with new password of perspective email/user
	query, err := DB.Prepare("UPDATE user SET password=? WHERE email=?")
	checkErr(err)
	query.Exec(password, email)

	lastPage = "passwordReset"
	http.Redirect(w, r, "/", http.StatusSeeOther)

	defer DB.Close()
}
