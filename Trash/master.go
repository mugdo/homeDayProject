package giveservice

import (
	"html/template"
	"net/http"
	"github.com/gorilla/sessions"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("frontEnd/*.html"))
}

var store = sessions.NewCookieStore([]byte("mysession"))
var lastPage = ""

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	lastPage = "index"

	session, _ := store.Get(r, "mysession")

	if session.Values["isLogin"] == nil {
		session.Values["isLogin"] = false
	}

	data := map[string]interface{}{
		"username": session.Values["username"],
		"password": session.Values["password"],
		"isLogged": session.Values["isLogin"],
		"title": "Home",
	}

	tpl.ExecuteTemplate(w, "index.html", data)
}
func Problem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	lastPage = "problem"

	session, _ := store.Get(r, "mysession")

	if session.Values["isLogin"] == nil {
		session.Values["isLogin"] = false
	}

	data := map[string]interface{}{
		"username": session.Values["username"],
		"password": session.Values["password"],
		"isLogged": session.Values["isLogin"],
		"title": "Problem",
	}

	tpl.ExecuteTemplate(w, "index.html", data)
}
func About(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	lastPage = "about"

	session, _ := store.Get(r, "mysession")

	if session.Values["isLogin"] == nil {
		session.Values["isLogin"] = false
	}

	data := map[string]interface{}{
		"username": session.Values["username"],
		"password": session.Values["password"],
		"isLogged": session.Values["isLogin"],
		"title": "About",
	}

	tpl.ExecuteTemplate(w, "index.html", data)
}
func Contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	lastPage = "contact"

	session, _ := store.Get(r, "mysession")

	if session.Values["isLogin"] == nil {
		session.Values["isLogin"] = false
	}

	data := map[string]interface{}{
		"username": session.Values["username"],
		"password": session.Values["password"],
		"isLogged": session.Values["isLogin"],
		"title": "Contact",
	}

	tpl.ExecuteTemplate(w, "index.html", data)
}
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	session, _ := store.Get(r, "mysession")

	data := map[string]interface{}{
		"title": "Login",
	}
	
	if session.Values["isLogin"] == true {
		http.Redirect(w, r, "/redirect", http.StatusSeeOther)
	} else {
		tpl.ExecuteTemplate(w, "index.html", data)
	}
}
func LoginCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "abc" && password == "123" {
		session, _ := store.Get(r, "mysession")
		session.Values["username"] = username
		session.Values["password"] = password
		session.Values["isLogin"] = true
		session.Save(r, w)

		http.Redirect(w, r, "/redirect", http.StatusSeeOther)
	} else {
		data := map[string]interface{}{
			"err": "Invalid username or password",
			"title": "Login",
		}
		tpl.ExecuteTemplate(w, "index.html", data)
	}
}
func Redirect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	if lastPage == "problem" {
		http.Redirect(w, r, "/problem", http.StatusSeeOther)
	} else if lastPage == "about" {
		http.Redirect(w, r, "/about", http.StatusSeeOther)
	} else if lastPage == "contact" {
		http.Redirect(w, r, "/contact", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	session, _ := store.Get(r, "mysession")
	session.Values["username"] = ""
	session.Values["password"] = ""
	session.Values["isLogin"] = false

	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
