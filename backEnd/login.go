package backEnd

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	session, _ := store.Get(r, "mysession")

	if session.Values["isLogin"] == true {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		Info = map[string]interface{}{
			"PageTitle":  "Login",
			"LastPage":   lastPage,
			"PopUpCause": popUpCause,
		}
		tpl.ExecuteTemplate(w, "login.gohtml", Info)
		popUpCause = ""
	}
}
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func LoginCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))

	DB := dbConn()

	Info = map[string]interface{}{
		"Username": username,
	}

	//getting original password
	var originalPassword string
	_ = DB.QueryRow("SELECT password FROM user WHERE username=?", username).Scan(&originalPassword)

	if checkPasswordHash(password, originalPassword) { //if password matched
		session, _ := store.Get(r, "mysession")
		session.Values["username"] = username
		session.Values["password"] = password
		session.Values["isLogin"] = true
		session.Save(r, w)

		Info = map[string]interface{}{
			"Username": session.Values["username"],
			"Password": session.Values["password"],
			"IsLogged": session.Values["isLogin"],
		}
		http.Redirect(w, r, lastPage, http.StatusSeeOther)
	} else { //if password not matched
		Info["ErrPassword"] = "Invalid password"

		tpl.ExecuteTemplate(w, "login.gohtml", Info)
	}

	defer DB.Close()
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
