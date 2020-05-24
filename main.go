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

	r.HandleFunc("/login", backEnd.Login)
	r.HandleFunc("/logout", backEnd.Logout)

	r.HandleFunc("/register", backEnd.Register)

	r.PathPrefix("/check/").HandlerFunc(backEnd.CheckDB)
	r.PathPrefix("/verify-email/token=").HandlerFunc(backEnd.EmailVerifiation)

	r.PathPrefix("/reset").HandlerFunc(backEnd.Reset)
	r.PathPrefix("/passReset/token=").HandlerFunc(backEnd.PassReset)

	r.HandleFunc("/problem", backEnd.Problem)
	r.PathPrefix("/problemView/").HandlerFunc(backEnd.ProblemView)
	r.PathPrefix("/origin/").HandlerFunc(backEnd.Origin)

	r.PathPrefix("/submit").HandlerFunc(backEnd.Submit)
	r.PathPrefix("/lang=").HandlerFunc(backEnd.GetLanguage)

	r.HandleFunc("/result", backEnd.Result)
	r.PathPrefix("/verdict/subID=").HandlerFunc(backEnd.Verdict)
	
	//URI related

	r.HandleFunc("/scrap", backEnd.Scrap)
	// r.HandleFunc("/toph", backEnd.Toph)
	// r.HandleFunc("/des", backEnd.Des)
	r.HandleFunc("/test2", backEnd.Test2)
	r.HandleFunc("/test1", backEnd.Test1)
	r.HandleFunc("/testSub", backEnd.TestSub)

	//for serving javascripts & css files
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	//A Custom Page Not Found route
	r.NotFoundHandler = http.HandlerFunc(backEnd.PageNotFound)

	//for localhost server
	http.ListenAndServe(":8080", r)
}
