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
		lastPage = "/"
	}

	session, _ := store.Get(r, "mysession")
	Info = map[string]interface{}{
		"Username":  session.Values["username"],
		"Password":  session.Values["password"],
		"IsLogged":  session.Values["isLogin"],
		"LastPage":  lastPage,
		"PageTitle": "Homepage",
	}

	tpl.ExecuteTemplate(w, "index.gohtml", Info)
	lastPage = "/"
}

func About(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	
	lastPage = "about"
	session, _ := store.Get(r, "mysession")
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
	Info = map[string]interface{}{
		"Username":  session.Values["username"],
		"Password":  session.Values["password"],
		"IsLogged":  session.Values["isLogin"],
		"PageTitle": "Contact",
	}

	tpl.ExecuteTemplate(w, "contact.gohtml", Info)
}
func PageNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	errorPage(w, http.StatusNotFound)
}
func errorPage(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode) //http.StatusNotFound = 404 //

	Info = map[string]interface{}{
		"StatusCode": statusCode,
	}
	tpl.ExecuteTemplate(w, "404.gohtml", Info)
}
