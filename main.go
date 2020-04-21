package main

import (
	"backEnd"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	//(instead of default 'http' router) using Gorilla mux router
	r := mux.NewRouter()

	//just a message for ensuring that local server is running
	fmt.Println("Local Server is running...")

	//for serving perspective pages
	r.HandleFunc("/", backEnd.Index)
	r.HandleFunc("/about", backEnd.About)
	r.HandleFunc("/contact", backEnd.Contact)
	r.HandleFunc("/redirect", backEnd.Redirect)

	r.HandleFunc("/login", backEnd.Login)
	r.HandleFunc("/loginCheck", backEnd.LoginCheck)
	r.HandleFunc("/logout", backEnd.Logout)
	r.HandleFunc("/register", backEnd.Register)
	r.HandleFunc("/doRegister", backEnd.DoRegister)

	r.PathPrefix("/check").HandlerFunc(backEnd.CheckDB)

	r.HandleFunc("/problem", backEnd.Problem)
	r.PathPrefix("/problemView").HandlerFunc(backEnd.ProblemView)

	r.HandleFunc("/submit", backEnd.Submit)
	r.PathPrefix("/submission").HandlerFunc(backEnd.Submission)
	r.PathPrefix("/lang").HandlerFunc(backEnd.GetLanguage)
	r.HandleFunc("/verdict", backEnd.Verdict)

	// r.HandleFunc("/scrap", backEnd.Scrap)
	// r.HandleFunc("/toph", backEnd.Toph)
	// r.HandleFunc("/des", backEnd.Des)
	r.HandleFunc("/test2", backEnd.Test2)
	r.HandleFunc("/test1", backEnd.Test1)
	r.HandleFunc("/testSub", backEnd.TestSub)

	//for serving javascripts & css files
	r.PathPrefix("/assests/").Handler(http.StripPrefix("/assests/", http.FileServer(http.Dir("assests"))))

	//A Custom Page Not Found route
	r.NotFoundHandler = http.HandlerFunc(backEnd.PageNotFound)

	//for localhost server
	http.ListenAndServe(":8080", r)
}
