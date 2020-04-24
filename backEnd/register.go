package backEnd

import (
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
