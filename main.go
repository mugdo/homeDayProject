package main

import (
	"fmt"
	"backEnd"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	//(instead of default 'http' router) using Gorilla mux router
	r := mux.NewRouter()

	//just a message for ensuring that local server is running
	fmt.Println("Local Server is running...")

	//for serving perspective pages
	r.HandleFunc("/", giveservice.Index)
	r.HandleFunc("/problem", giveservice.Problem)
	r.HandleFunc("/problemView", giveservice.ProblemView)
	r.HandleFunc("/login", giveservice.Login)
	r.HandleFunc("/loginCheck", giveservice.LoginCheck)
	r.HandleFunc("/redirect", giveservice.Redirect)
	r.HandleFunc("/register", giveservice.Register)
	r.HandleFunc("/doRegister", giveservice.DoRegister)
	r.HandleFunc("/logout", giveservice.Logout)
	r.HandleFunc("/about", giveservice.About)
	r.HandleFunc("/contact", giveservice.Contact)
	r.HandleFunc("/submit", giveservice.Submit)
	r.HandleFunc("/result", giveservice.Result)
	r.HandleFunc("/scrap", giveservice.Scrap)
	r.HandleFunc("/toph", giveservice.Toph)
	
	//for serving javascripts & css files
	r.PathPrefix("/assests/").Handler(http.StripPrefix("/assests/", http.FileServer(http.Dir("assests"))))

	//A Custom Page Not Found route
	r.NotFoundHandler = http.HandlerFunc(giveservice.PageNotFound)

	//for localhost server
	http.ListenAndServe(":8080", r)
}