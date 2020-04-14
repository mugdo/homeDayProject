package backEnd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/sessions"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("frontEnd/*"))
}

var store = sessions.NewCookieStore([]byte("mysession"))
var lastPage = ""
var title, timeLimit, memoryLimit, src = "-", "-", "-", "-"

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
func getTimeLimit(body string) (a, b string) {
	//fmt.Println(body)

	need := "ms"

	index := strings.Index(body, need)
	fmt.Println("got Time")
	fmt.Println(index)

	runes := []rune(body)
	sb := "-"
	if index != -1 {
		sb = string(runes[0 : index+2])
	}
	fmt.Println("Time Limit:", sb)

	// if oj == "CodeChef" {
	// 	need = "B"
	// }
	need = "B"
	index2 := strings.Index(body, need)
	fmt.Println("got Memory")
	fmt.Println(index2)

	//runes = []rune(body)
	sbm := "-"
	if index2 != -1 {
		sbm = string(runes[index+2 : index2+1])
	}
	fmt.Println("Memory Limit:", sbm)

	return sb, sbm
}

func calculateTitle(oj, pNum string) {
	fmt.Println("In GetTitle")
	client := &http.Client{
		Jar: cookieJar,
	}
	pURL := "https://vjudge.net/problem/" + oj + "-" + pNum
	// response, err := client.Get(pURL)
	// if err != nil {
	// 	log.Fatalln("Error fetching response. ", err)
	// }
	// defer response.Body.Close()

	req, err := http.NewRequest("GET", pURL, nil)
	req.Header.Add("Content-Type", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	response, err := client.Do(req)
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	timeLimitText := document.Find("dl[class='card row']").Find("dd[class='col-sm-7']").Text()
	timeLimit, memoryLimit = getTimeLimit(timeLimitText)

	title = document.Find("div[id='prob-title']").Find("h2").Text()
	src, _ = document.Find("iframe").Attr("src")

	//fmt.Println("src=", src, "iv", iv)
}
func ProblemView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
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

	fmt.Println(oj, pNum)
	//fmt.Println("Scrap started")

	calculateTitle(oj, pNum)

	//for uva & uvalive pdf description
	var segment string
	if oj == "UVA" || oj == "UVALive" {
		temp, _ := strconv.Atoi(pNum)

		IntSegment := temp / 100
		segment = strconv.Itoa(IntSegment)
	}

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
	// response, err := client.Get(pURL)
	// if err != nil {
	// 	log.Fatalln("Error fetching response. ", err)
	// }
	// defer response.Body.Close()
	req, err := http.NewRequest("GET", pURL, nil)
	req.Header.Add("Content-Type", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	response, err := client.Do(req)
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
		"Segment":     segment,
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

	type LanguagePack struct {
		LangValue int
		LangName  string
	}
	var languagePack []LanguagePack
	oj := "CodeForces"
	if oj == "CodeForces" {
		languagePack = []LanguagePack{
			LanguagePack{LangValue: 52, LangName: "Clang++17 Diagnostics"},
			LanguagePack{LangValue: 9, LangName: "C# Mono 5.18"},
			LanguagePack{LangValue: 28, LangName: "D DMD32 v2.091.0"},
			LanguagePack{LangValue: 3, LangName: "Delphi 7"},
			LanguagePack{LangValue: 4, LangName: "Free Pascal 3.0.2"},
			LanguagePack{LangValue: 43, LangName: "GNU GCC C11 5.1.0"},
			LanguagePack{LangValue: 42, LangName: "GNU G++11 5.1.0"},
			LanguagePack{LangValue: 50, LangName: "GNU G++14 6.4.0"},
			LanguagePack{LangValue: 54, LangName: "GNU G++17 7.3.0"},
			LanguagePack{LangValue: 61, LangName: "GNU G++17 9.2.0 (64 bit, msys 2)"},
			LanguagePack{LangValue: 32, LangName: "Go 1.14"},
			LanguagePack{LangValue: 12, LangName: "Haskell GHC 8.6.3"},
			LanguagePack{LangValue: 36, LangName: "Java 1.8.0_162"},
			LanguagePack{LangValue: 60, LangName: "Java 11.0.5"},
			LanguagePack{LangValue: 34, LangName: "JavaScript V8 4.8.0"},
			LanguagePack{LangValue: 48, LangName: "Kotlin 1.3.70"},
			LanguagePack{LangValue: 2, LangName: "Microsoft Visual C++ 2010"},
			LanguagePack{LangValue: 59, LangName: "Microsoft Visual C++ 2017"},
			LanguagePack{LangValue: 55, LangName: "Node.js 9.4.0"},
			LanguagePack{LangValue: 19, LangName: "OCaml 4.02.1"},
			LanguagePack{LangValue: 51, LangName: "PascalABC.NET 3.4.2"},
			LanguagePack{LangValue: 13, LangName: "Perl 5.20.1"},
			LanguagePack{LangValue: 6, LangName: "PHP 7.2.13"},
			LanguagePack{LangValue: 40, LangName: "PyPy 2.7 (7.2.0)"},
			LanguagePack{LangValue: 41, LangName: "PyPy 3.6 (7.2.0)"},
			LanguagePack{LangValue: 7, LangName: "Python 2.7.15"},
			LanguagePack{LangValue: 31, LangName: "Python 3.7.2"},
			LanguagePack{LangValue: 8, LangName: "Ruby 2.0.0p645"},
			LanguagePack{LangValue: 49, LangName: "Rust 1.42.0"},
			LanguagePack{LangValue: 20, LangName: "Scala 2.12.8"},
		}
	}
	fmt.Println(languagePack[0].LangName, languagePack[0].LangValue)
	//return languagePack
}
func Submission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	lastPage = ""

	path := r.URL.Path
	runes := []rune(path)
	path = string(runes[12:])
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
	languagePack := getLenguage(oj)
	session, _ := store.Get(r, "mysession")
	data := map[string]interface{}{
		"username":  session.Values["username"],
		"password":  session.Values["password"],
		"isLogged":  session.Values["isLogin"],
		"pageTitle": "Submission",

		"Oj":           oj,
		"PNum":         pNum,
		"PName":        title,
		"TimeLimit":    timeLimit,
		"MemoryLimit":  memoryLimit,
		"LanguagePack": languagePack,
	}
	tpl.ExecuteTemplate(w, "submission.gohtml", data)
}
