package backEnd

import (
	"database/sql"
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
			"pageTitle": "Login",
			"LastPage":  lastPage,
		}
		tpl.ExecuteTemplate(w, "login.gohtml", Info)

		lastPage = "login"
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
		"username": username,
	}

	//checking for username really exist or not
	var originalPassword string
	err := DB.QueryRow("SELECT password FROM user WHERE username=?", username).Scan(&originalPassword)

	if err == sql.ErrNoRows {
		//username not found (found no rows)
		Info["errUsername"] = "No Account found with this username. Try again"
		Info["username"] = ""

		tpl.ExecuteTemplate(w, "login.gohtml", Info)
	} else {
		//no error on db.QueryRow (username found & original password achieved)
		if checkPasswordHash(password, originalPassword) {
			session, _ := store.Get(r, "mysession")
			session.Values["username"] = username
			session.Values["password"] = password
			session.Values["isLogin"] = true
			session.Save(r, w)

			Info = map[string]interface{}{
				"username": session.Values["username"],
				"password": session.Values["password"],
				"isLogged": session.Values["isLogin"],
			}

			http.Redirect(w, r, "/redirect", http.StatusSeeOther)
		} else {
			Info["errPassword"] = "Invalid password"

			tpl.ExecuteTemplate(w, "login.gohtml", Info)
		}
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
