package backEnd

import (
	"database/sql"
	"net/http"
	"strings"
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
	confirmPassword := strings.TrimSpace(r.FormValue("confirmPassword"))

	db := dbConn()

	Info = map[string]interface{}{
		"fullName": fullName,
		"email":    email,
		"username": username,
	}

	//checking for email already exist or not
	var temp string
	err1 := db.QueryRow("SELECT id FROM user WHERE email=?", email).Scan(&temp)
	checkErr(err1)
	if err1 == sql.ErrNoRows {
		//email available (found no rows)
		Info["errEmail"] = ""
	} else {
		Info["errEmail"] = "Email already registered. Choose another one"
	}

	//checking for username already exist or not
	temp = ""
	err2 := db.QueryRow("SELECT id FROM user WHERE username=?", username).Scan(&temp)
	checkErr(err2)
	if err2 == sql.ErrNoRows {
		//username available (found no rows)
		Info["errUsername"] = ""
	} else {
		Info["errUsername"] = "Username already taken. Choose another one"
	}

	//checking for password & confirmPassword are same or not
	if password == confirmPassword {
		//passwords are same
		Info["errPassword"] = ""
	} else {
		Info["errPassword"] = "Password mismatched. Put cautiously"
	}

	//now do regitration
	if Info["errEmail"] != "" || Info["errUsername"] != "" || Info["errPassword"] != "" {
		if Info["errEmail"] != "" {
			Info["email"] = ""
		}
		if Info["errUsername"] != "" {
			Info["username"] = ""
		}
		tpl.ExecuteTemplate(w, "register.gohtml", Info)
	} else {
		insertQuery, err := db.Prepare("INSERT INTO user(fullName, email, username, password, confirmPassword) VALUES(?,?,?,?,?)")
		checkErr(err)
		insertQuery.Exec(fullName, email, username, password, confirmPassword)
		println("Registration Done")

		lastPage = "login"
		http.Redirect(w, r, "/redirect", http.StatusSeeOther)
	}

	defer db.Close()
}
