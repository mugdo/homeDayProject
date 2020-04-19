package backEnd

import (
	"database/sql"
	"net/http"
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

	fullName := r.FormValue("fullName")
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	db := dbConn()

	Info = map[string]interface{}{
		"fullName": fullName,
		"email":    email,
		"username": username,
	}

	//checking for email already exist or not
	var id string
	err1 := db.QueryRow("SELECT id FROM user WHERE email=?", email).Scan(&id)
	if err1 == sql.ErrNoRows {
		//email available (found no rows)
		Info["errEmail"] = ""
	} else if err1 != nil {
		panic(err1.Error())
	} else {
		Info["errEmail"] = "Email already registered. Choose another one"
	}

	//checking for username already exist or not
	id = ""
	err2 := db.QueryRow("SELECT id FROM user WHERE username=?", username).Scan(&id)
	if err2 == sql.ErrNoRows {
		//username available (found no rows)
		Info["errUsername"] = ""
	} else if err2 != nil {
		panic(err2.Error())
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
		if err != nil {
			panic(err.Error())
		}
		insertQuery.Exec(fullName, email, username, password, confirmPassword)
		println("Registration Done")

		lastPage = "login"
		http.Redirect(w, r, "/redirect", http.StatusSeeOther)
	}

	defer db.Close()
}