package backEnd

import (
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	lastPage = "/"
	session, _ := store.Get(r, "mysession")

	Info["Username"] = session.Values["username"]
	Info["Password"] = session.Values["password"]
	Info["IsLogged"] = session.Values["isLogin"]
	Info["LastPage"] = lastPage
	Info["PopUpCause"] = popUpCause
	Info["PageTitle"] = "Homepage"

	tpl.ExecuteTemplate(w, "index.gohtml", Info)
	popUpCause = ""
	Info["PopUpCause"]=popUpCause
}

func About(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	lastPage = "about"
	session, _ := store.Get(r, "mysession")

	Info["Username"] = session.Values["username"]
	Info["Password"] = session.Values["password"]
	Info["IsLogged"] = session.Values["isLogin"]
	Info["PageTitle"] = "About"

	tpl.ExecuteTemplate(w, "about.gohtml", Info)
}
func Contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	lastPage = "contact"
	session, _ := store.Get(r, "mysession")

	Info["Username"] = session.Values["username"]
	Info["Password"] = session.Values["password"]
	Info["IsLogged"] = session.Values["isLogin"]
	Info["PageTitle"] = "Contact"

	tpl.ExecuteTemplate(w, "contact.gohtml", Info)
}
func PageNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	errorPage(w, http.StatusNotFound)
}
func errorPage(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode) //status code such as: 400, 404 etc.

	Info["StatusCode"] = statusCode

	tpl.ExecuteTemplate(w, "404.gohtml", Info)
}
