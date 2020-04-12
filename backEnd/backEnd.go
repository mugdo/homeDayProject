package backEnd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
	"strings"
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
		"username":  session.Values["username"],
		"password":  session.Values["password"],
		"isLogged":  session.Values["isLogin"],
		"pageTitle": "Homepage",
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

	oj := r.FormValue("oj")
	pNum := r.FormValue("pNum")
	pName := r.FormValue("pName")

	body, err := pSearch(oj, pNum, pName)
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(body))

	pList := getPList(body)
	found := true

	if len(pList.Data) == 0 {
		found = false
	}

	data := map[string]interface{}{
		"username":  session.Values["username"],
		"password":  session.Values["password"],
		"isLogged":  session.Values["isLogin"],
		"PList":     pList,
		"Found":     found,
		"Oj":        oj,
		"PNum":      pNum,
		"PName":     pName,
		"pageTitle": "Problem",
	}

	tpl.ExecuteTemplate(w, "problem.gohtml", data)
}
func ProblemView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	lastPage = "problem"

	path := r.URL.Path
	runes := []rune(path)
	path = string(runes[13:])
	fmt.Println(path)

	need := "-"
	index := strings.Index(path, need)
	// fmt.Println("got Dash")
	// fmt.Println(index)
	var oj, pNum string
	runes = []rune(path)
	oj = string(runes[0:3])

	if oj == "计蒜客" {
		//oj = string(runes[0:3])
		pNum = string(runes[4:])
	} else {
		oj = string(runes[0:index])
		pNum = string(runes[index+1:])
	}

	fmt.Println(oj,pNum)
	//fmt.Println("Scrap started")

	title, timeLimit, memoryLimit, src := getTitle(oj, pNum)

	// fmt.Println("After getting title ", title)
	// fmt.Println("After getting TL ", timeLimit)
	// fmt.Println("After getting ML ", memoryLimit)
	// fmt.Println("After getting src ", src)
	// fmt.Println("Finished Scrap")

	// getting problem description///////////////
	url := "https://vjudge.net" + src
	// body, _ := rGET(url)
	// fmt.Printf(string(body))
	client := &http.Client{
		Jar: cookieJar,
	}
	pURL := url
	response, err := client.Get(pURL)
	if err != nil {
		log.Fatalln("Error fetching response. ", err)
	}
	defer response.Body.Close()
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}
	textArea := document.Find("textarea").Text()
	fmt.Println("Textarea:::::::::::", textArea)
	fmt.Println("Textarea:::::::::::")

	b := []byte(textArea)
	//fmt.Println("Byte:",b)

	type Inner2 struct {
		Format  string `json:"format"`
		Content string `json:"content"`
	}
	type Inner struct {
		Title string `json:"title"`
		Value Inner2 `json:"value"`
	}
	type Res struct {
		Trustable bool    `json:"trustable"`
		Sections  []Inner `json:"sections"`
	}
	var res Res
	json.Unmarshal(b, &res)

	//Eliminating Default CSS on Example-Input-Output
	// need = "<script"

	// index = strings.Index(res.Sections[0].Value.Content, need)
	// fmt.Println("got Script")
	// fmt.Println(index)

	// if index != -1 {
	// 	runes = []rune(res.Sections[0].Value.Content)
	// 	res.Sections[0].Value.Content = string(runes[index:])
	// }

	problem := []map[string]interface{}{}

	for i := 0; i < len(res.Sections); i++ {

		mp := map[string]interface{}{
			"Title":   template.HTML(res.Sections[i].Title),
			"Content": template.HTML(res.Sections[i].Value.Content),
		}
		problem = append(problem, mp)
	}

	session, _ := store.Get(r, "mysession")
	data := map[string]interface{}{
		"username":  session.Values["username"],
		"password":  session.Values["password"],
		"isLogged":  session.Values["isLogin"],
		"pageTitle": title,

		"Oj":          oj,
		"PNum":        pNum,
		"PName":       title,
		"TimeLimit":   timeLimit,
		"MemoryLimit": memoryLimit,
		"Problem":     problem,
		//"ProblemContent": problemContent,
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
		"username":  session.Values["username"],
		"password":  session.Values["password"],
		"isLogged":  session.Values["isLogin"],
		"pageTitle": "About",
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
		"username":  session.Values["username"],
		"password":  session.Values["password"],
		"isLogged":  session.Values["isLogin"],
		"pageTitle": "Contact",
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
			"pageTitle": "Login",
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
		"username": username,
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
			"pageTitle": "Registration",
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
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	db := dbConn()

	data := map[string]interface{}{
		"fullName": fullName,
		"email":    email,
		"username": username,
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
	id = ""
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
	if data["errEmail"] != "" || data["errUsername"] != "" || data["errPassword"] != "" {
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
func PageNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	data := map[string]interface{}{
		"pageTitle": "Page Not Found",
	}
	tpl.ExecuteTemplate(w, "404.gohtml", data)
}
func Result(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	data := map[string]interface{}{
		"pageTitle": "Verdict",
	}
	tpl.ExecuteTemplate(w, "result.gohtml", data)
}
func TestPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	body, err := rLogin("ajudgebd", "aj199273", "https://vj.z180.cn")
	if err != nil{
		panic(err)
	}
	res := string(body)
	fmt.Fprint(w, res)
}