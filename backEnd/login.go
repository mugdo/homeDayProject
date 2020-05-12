package backEnd

import (
	"golang.org/x/crypto/bcrypt"
	"html"
	"net/http"
	"strings"
)

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	session, _ := store.Get(r, "mysession")

	if r.Method != "POST" {
		if session.Values["isLogin"] == true {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			Info["PageTitle"] = "Login"
			Info["LastPage"] = lastPage
			Info["PopUpCause"] = popUpCause

			tpl.ExecuteTemplate(w, "login.gohtml", Info)
			Info["Username"] = ""
			Info["ErrPassword"] = ""
			popUpCause = ""
			Info["PopUpCause"] = popUpCause
		}
	} else if r.Method == "POST" {
		username := html.EscapeString(strings.TrimSpace(r.FormValue("username")))
		password := html.EscapeString(r.FormValue("password"))

		DB := dbConn()
		defer DB.Close()

		//getting original password from DB
		var originalPassword string
		_ = DB.QueryRow("SELECT password FROM user WHERE username=?", username).Scan(&originalPassword)

		if checkPasswordHash(password, originalPassword) { //if password matched
			session.Values["username"] = username
			session.Values["password"] = password
			session.Values["isLogin"] = true
			session.Save(r, w)

			Info["Username"] = session.Values["username"]
			Info["Password"] = session.Values["password"]
			Info["IsLogged"] = session.Values["isLogin"]

			http.Redirect(w, r, lastPage, http.StatusSeeOther)
		} else { //if password not matched
			Info["Username"] = username
			Info["ErrPassword"] = "Invalid password"

			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}
}
func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	session, _ := store.Get(r, "mysession")
	session.Values["username"] = ""
	session.Values["password"] = ""
	session.Values["isLogin"] = false

	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
