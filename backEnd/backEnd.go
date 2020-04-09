package backEnd

import (
	"database/sql"
	"html/template"
	"net/http" 
	"github.com/gorilla/sessions"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("frontEnd/*"))
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
		"username"	: session.Values["username"],
		"password"	: session.Values["password"],
		"isLogged"	: session.Values["isLogin"],
		"pageTitle"	: "Homepage",
	}

	tpl.ExecuteTemplate(w, "index.gohtml", data)
}

func Problem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	lastPage = "problem"

	session, _ := store.Get(r, "mysession")

	if session.Values["isLogin"] == nil {
		session.Values["isLogin"] = false
	}

	oj		:= r.FormValue("oj")
	pNum	:= r.FormValue("pNum")
	pName	:= r.FormValue("pName")

	body, err := pSearch(oj, pNum, pName)
	if(err!=nil){
		panic(err)
	}
	//fmt.Println(string(body))

	pList := getPList(body);
	found := true

	if len(pList.Data)==0 {
		found = false
	}

	data := map[string]interface{}{
		"username"	: session.Values["username"],
		"password"	: session.Values["password"],
		"isLogged"	: session.Values["isLogin"],
		"PList"		: pList,
		"Found"		: found,
		"Oj"		: oj,
		"PNum"		: pNum,
		"PName"		: pName,
		"pageTitle"	: "Problem",
	}

	tpl.ExecuteTemplate(w, "problem.gohtml", data)
}
func ProblemView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	lastPage = "problem"

	session, _ := store.Get(r, "mysession")

	if session.Values["isLogin"] == nil {
		session.Values["isLogin"] = false
	}

	body, err := pSearch("","","")
	if(err!=nil){
		panic(err)
	}
	//fmt.Println(string(body))

	pList := getPList(body);

	data := map[string]interface{}{
		"username"	: session.Values["username"],
		"password"	: session.Values["password"],
		"isLogged"	: session.Values["isLogin"],
		"PList"		: pList,
		"pageTitle"	: "",
	}

	tpl.ExecuteTemplate(w, "problemView.gohtml", data)
}
func About(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	lastPage = "about"

	session, _ := store.Get(r, "mysession")

	if session.Values["isLogin"] == nil {
		session.Values["isLogin"] = false
	}

	data := map[string]interface{}{
		"username"	: session.Values["username"],
		"password"	: session.Values["password"],
		"isLogged"	: session.Values["isLogin"],
		"pageTitle"	: "About",
	}

	tpl.ExecuteTemplate(w, "about.gohtml", data)
}
func Contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	lastPage = "contact"

	session, _ := store.Get(r, "mysession")

	if session.Values["isLogin"] == nil {
		session.Values["isLogin"] = false
	}

	data := map[string]interface{}{
		"username"	: session.Values["username"],
		"password"	: session.Values["password"],
		"isLogged"	: session.Values["isLogin"],
		"pageTitle"	: "Contact",
	}

	tpl.ExecuteTemplate(w, "contact.gohtml", data)
}
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	session, _ := store.Get(r, "mysession")
	
	if session.Values["isLogin"] == true {
		http.Redirect(w, r, "/redirect", http.StatusSeeOther)
	} else {
		data := map[string]interface{}{
			"pageTitle"	: "Login",
		}
		tpl.ExecuteTemplate(w, "login.gohtml", data)
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

	db := dbConn()

	data := map[string]interface{}{
		"username" : username,
	}

	//checking for username really exist or not
	var originalPassword string
	err := db.QueryRow("SELECT password FROM user WHERE username=?", username).Scan(&originalPassword)
	if err == sql.ErrNoRows {
		//username not found (found no rows)
		data["errUsername"] = "No Account found with this username. Try again"
		data["username"] = ""

		tpl.ExecuteTemplate(w, "login.gohtml", data)
	} else if err != nil {
		panic(err.Error())
	} else {
		//no error on db.QueryRow (username found & original password achieved)
		if password == originalPassword {
			session, _ := store.Get(r, "mysession")
			session.Values["username"] = username
			session.Values["password"] = password
			session.Values["isLogin"] = true
			session.Save(r, w)
	
			http.Redirect(w, r, "/redirect", http.StatusSeeOther)
		} else {
			data["errPassword"] = "Invalid password"

			tpl.ExecuteTemplate(w, "login.gohtml", data)
		}
	}

	defer db.Close()
}
func Redirect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

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
func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	session, _ := store.Get(r, "mysession")

	if session.Values["isLogin"] == true {
		http.Redirect(w, r, "/redirect", http.StatusSeeOther)
	} else {
		data := map[string]interface{}{
			"pageTitle"	: "Registration",
		}
		tpl.ExecuteTemplate(w, "register.gohtml", data)
	}
}
func DoRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	fullName := r.FormValue("fullName")
	email	 := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	db := dbConn()

	data := map[string]interface{}{
		"fullName" : fullName,
		"email" : email,
		"username" : username,
	}
	
	//checking for email already exist or not
	var id string
	err1 := db.QueryRow("SELECT id FROM user WHERE email=?", email).Scan(&id)
	if err1 == sql.ErrNoRows {
		//email available (found no rows)
		data["errEmail"] = ""
	} else if err1 != nil {
		panic(err1.Error())
	} else {
		data["errEmail"] = "Email already registered. Choose another one"
	}

	//checking for username already exist or not
	id=""
	err2 := db.QueryRow("SELECT id FROM user WHERE username=?", username).Scan(&id)
	if err2 == sql.ErrNoRows {
		//username available (found no rows)
		data["errUsername"] = ""
	} else if err2 != nil {
		panic(err2.Error())
	} else {
		data["errUsername"] = "Username already taken. Choose another one"
	}
	//checking for password & confirmPassword are same or not
	if password == confirmPassword {
		//passwords are same
		data["errPassword"] = ""
	} else {
		data["errPassword"] = "Password mismatched. Put cautiously"
	}

	//now do regitration
	if data["errEmail"] != "" || data["errUsername"] != "" || data["errPassword"] != ""{
		if data["errEmail"] != "" {
			data["email"] = ""
		}
		if data["errUsername"] != "" {
			data["username"] = ""
		}
		tpl.ExecuteTemplate(w, "register.gohtml", data)
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
func PageNotFound(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html")

	data := map[string]interface{}{
		"pageTitle"	: "Page Not Found",
	}
	tpl.ExecuteTemplate(w, "404.gohtml", data)
}
func Result(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	data := map[string]interface{}{
		"pageTitle"	: "Verdict",
	}
	tpl.ExecuteTemplate(w, "result.gohtml", data)
}