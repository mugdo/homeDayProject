package backEnd

import (
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	if lastPage != "tokenInvalid" &&
		lastPage != "tokenAlreadyVerified" &&
		lastPage != "tokenExpired" &&
		lastPage != "tokenVerifiedNow" &&
		lastPage != "tokenRequest" &&
		lastPage != "passwordRequest" &&
		lastPage != "passTokenInvalid" &&
		lastPage != "passTokenExpired" &&
		lastPage != "passwordReset" {
		lastPage = "index"
	}

	session, _ := store.Get(r, "mysession")
	if session.Values["isLogin"] == nil {
		session.Values["isLogin"] = false
	}

	Info = map[string]interface{}{
		"Username":  session.Values["username"],
		"Password":  session.Values["password"],
		"IsLogged":  session.Values["isLogin"],
		"LastPage":  lastPage,
		"PageTitle": "Homepage",
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
		"Username":  session.Values["username"],
		"Password":  session.Values["password"],
		"IsLogged":  session.Values["isLogin"],
		"PageTitle": "About",
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
		"Username":  session.Values["username"],
		"Password":  session.Values["password"],
		"IsLogged":  session.Values["isLogin"],
		"PageTitle": "Contact",
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
	errorPage(w,http.StatusNotFound)
}
func errorPage(w http.ResponseWriter,statusCode int){
	w.WriteHeader(statusCode) //http.StatusNotFound = 404 // 

	Info = map[string]interface{}{
		"StatusCode": statusCode,
	}
	tpl.ExecuteTemplate(w, "404.gohtml", Info)
}