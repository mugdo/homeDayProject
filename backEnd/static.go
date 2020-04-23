package backEnd

import (
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	
	if lastPage != "tokenInvalid" && lastPage != "tokenAlreadyVerified" && lastPage != "tokenExpired" && lastPage != "tokenVerifiedNow" {
		lastPage = "index"
	}

	session, _ := store.Get(r, "mysession")
	if session.Values["isLogin"] == nil {
		session.Values["isLogin"] = false
	}

	Info = map[string]interface{}{
		"username":  session.Values["username"],
		"password":  session.Values["password"],
		"isLogged":  session.Values["isLogin"],
		"LastPage":  lastPage,
		"pageTitle": "Homepage",
	}

	tpl.ExecuteTemplate(w, "index.gohtml", Info)
	lastPage = "index"
}

func About(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	lastPage = "about"

	session, _ := store.Get(r, "mysession")

	if session.Values["isLogin"] == nil {
		session.Values["isLogin"] = false
	}

	Info = map[string]interface{}{
		"username":  session.Values["username"],
		"password":  session.Values["password"],
		"isLogged":  session.Values["isLogin"],
		"pageTitle": "About",
	}

	tpl.ExecuteTemplate(w, "about.gohtml", Info)
}
func Contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	lastPage = "contact"

	session, _ := store.Get(r, "mysession")

	if session.Values["isLogin"] == nil {
		session.Values["isLogin"] = false
	}

	Info = map[string]interface{}{
		"username":  session.Values["username"],
		"password":  session.Values["password"],
		"isLogged":  session.Values["isLogin"],
		"pageTitle": "Contact",
	}

	tpl.ExecuteTemplate(w, "contact.gohtml", Info)
}
func Redirect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	if lastPage == "problem" {
		http.Redirect(w, r, "/problem", http.StatusSeeOther)
	} else if lastPage == "about" {
		http.Redirect(w, r, "/about", http.StatusSeeOther)
	} else if lastPage == "contact" {
		http.Redirect(w, r, "/contact", http.StatusSeeOther)
	} else if lastPage == "login" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func PageNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	Info = map[string]interface{}{
		"pageTitle": "Page Not Found",
	}
	tpl.ExecuteTemplate(w, "404.gohtml", Info)
}
