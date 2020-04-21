package backEnd

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/smtp"
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
			"pageTitle": "Registration",
		}
		tpl.ExecuteTemplate(w, "register.gohtml", Info)
	}
}

func sendMail(email, link string) {
	// Choose auth method and set it up
	auth := smtp.PlainAuth("", "ajudge.bd", "aj199273", "smtp.gmail.com")

	// Here we do it all: connect to our server, set up a message and send it
	to := []string{email}
	msg := []byte("To: " + email + "\r\n" +
		"Subject: Ajudge Email verification\r\n" +
		"\r\n" +
		"Here’s the link of account activation. Click on the:\r\n" +
		link)
	err := smtp.SendMail("smtp.gmail.com:587", auth, "ajudge.bd@gmail.com", to, msg)
	checkErr(err)
}
func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)

	hasher := md5.New()
	hasher.Write(b)
	return hex.EncodeToString(hasher.Sum(nil))
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

	DB := dbConn()

	//do regitration
	CurrentTime := time.Now()
	token := generateToken()
	tokenPeriod := time.Now().Unix()

	insertQuery, err := DB.Prepare("INSERT INTO user(fullName, email, username, password, createdAt, isVerified, token, tokenPeriod) VALUES(?,?,?,?,?,?,?,?)")
	checkErr(err)
	insertQuery.Exec(fullName, email, username, password, CurrentTime, 0, token, tokenPeriod)
	println("Registration Done")

	link := "http://localhost:8080/verify-email/token=" + token
	sendMail(email, link)

	lastPage = "login"
	http.Redirect(w, r, "/redirect", http.StatusSeeOther)

	defer DB.Close()
}
