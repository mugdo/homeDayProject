package backEnd

import (
	"database/sql"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	lastPage = "index"

	session, _ := store.Get(r, "mysession")
	if session.Values["isLogin"] == nil {
		session.Values["isLogin"] = false
	}

	Info = map[string]interface{}{
		"username":  session.Values["username"],
		"password":  session.Values["password"],
		"isLogged":  session.Values["isLogin"],
		"pageTitle": "Homepage",
	}

	tpl.ExecuteTemplate(w, "index.gohtml", Info)
}
func Problem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	lastPage = "problem"

	session, _ := store.Get(r, "mysession")
	if session.Values["isLogin"] == nil {
		session.Values["isLogin"] = false
	}

	OJ := r.FormValue("OJ")
	pNum := r.FormValue("pNum")
	pName := r.FormValue("pName")

	body, err := pSearch(OJ, pNum, pName)
	checkErr(err)

	pList := getPList(body)
	found := true

	if len(pList.Data) == 0 {
		found = false
	}

	Info = map[string]interface{}{
		"username":  session.Values["username"],
		"password":  session.Values["password"],
		"isLogged":  session.Values["isLogin"],
		"PList":     pList,
		"Found":     found,
		"OJ":        OJ,
		"PNum":      pNum,
		"PName":     pName,
		"pageTitle": "Problem",
	}

	tpl.ExecuteTemplate(w, "problem.gohtml", Info)
}
func ProblemView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	lastPage = "problem"

	path := r.URL.Path
	runes := []rune(path)
	OJpNum := string(runes[13:])

	need := "-"
	index := strings.Index(OJpNum, need)
	var OJ, pNum string
	runes = []rune(OJpNum)
	OJ = string(runes[0:3])

	if OJ == "计蒜客" {
		pNum = string(runes[4:])
	} else {
		OJ = string(runes[0:index])
		pNum = string(runes[index+1:])
	}

	//Finding problem Title, Time limit, Memory limit, Description source
	findPResource(OJ, pNum)

	//for uva & uvalive pdf description
	var uvaSegment string
	if OJ == "UVA" || OJ == "UVALive" {
		temp, _ := strconv.Atoi(pNum)

		IntSegment := temp / 100
		uvaSegment = strconv.Itoa(IntSegment)
	}

	// getting problem description///////////////
	pURL := "https://vjudge.net" + pDesSrc
	req, err := http.NewRequest("GET", pURL, nil)
	req.Header.Add("Content-Type", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	response, err := client.Do(req)
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	checkErr(err)

	textArea := document.Find("textarea").Text()
	b := []byte(textArea)

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

	problem := []map[string]interface{}{}

	for i := 0; i < len(res.Sections); i++ {
		//Eliminating Default CSS on Example-Input-Output
		styleBody := res.Sections[i].Value.Content
		content := removeStyle(styleBody)

		mp := map[string]interface{}{
			"Title":   template.HTML(res.Sections[i].Title),
			"Content": template.HTML(content),
		}
		problem = append(problem, mp)
	}

	//checking whether problem submission allowed or not
	tempP, _ := pSearch(OJ, pNum, "")
	tempList := getPList(tempP)

	allowSubmit := false
	if tempList.Data[0].AllowSubmit == true && tempList.Data[0].Status == 0 {
		allowSubmit = true
	}

	session, _ := store.Get(r, "mysession")
	Info = map[string]interface{}{
		"username":    session.Values["username"],
		"password":    session.Values["password"],
		"isLogged":    session.Values["isLogin"],
		"pageTitle":   pTitle,
		"OJ":          OJ,
		"PNum":        pNum,
		"AllowSubmit": allowSubmit,
		"UvaSegment":  uvaSegment,
		"PName":       pTitle,
		"TimeLimit":   pTimeLimit,
		"MemoryLimit": pMemoryLimit,
		"Problem":     problem,
	}

	tpl.ExecuteTemplate(w, "problemView.gohtml", Info)
}
func About(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
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
	w.Header().Set("Content-Type", "text/html")
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
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	session, _ := store.Get(r, "mysession")

	if session.Values["isLogin"] == true {
		http.Redirect(w, r, "/redirect", http.StatusSeeOther)
	} else {
		Info = map[string]interface{}{
			"pageTitle": "Login",
		}
		tpl.ExecuteTemplate(w, "login.gohtml", Info)
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

	Info = map[string]interface{}{
		"username": username,
	}

	//checking for username really exist or not
	var originalPassword string
	err := db.QueryRow("SELECT password FROM user WHERE username=?", username).Scan(&originalPassword)
	if err == sql.ErrNoRows {
		//username not found (found no rows)
		Info["errUsername"] = "No Account found with this username. Try again"
		Info["username"] = ""

		tpl.ExecuteTemplate(w, "login.gohtml", Info)
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
			Info["errPassword"] = "Invalid password"

			tpl.ExecuteTemplate(w, "login.gohtml", Info)
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
		Info = map[string]interface{}{
			"pageTitle": "Registration",
		}
		tpl.ExecuteTemplate(w, "register.gohtml", Info)
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

	Info = map[string]interface{}{
		"fullName": fullName,
		"email":    email,
		"username": username,
	}

	//checking for email already exist or not
	var id string
	err1 := db.QueryRow("SELECT id FROM user WHERE email=?", email).Scan(&id)
	if err1 == sql.ErrNoRows {
		//email available (found no rows)
		Info["errEmail"] = ""
	} else if err1 != nil {
		panic(err1.Error())
	} else {
		Info["errEmail"] = "Email already registered. Choose another one"
	}

	//checking for username already exist or not
	id = ""
	err2 := db.QueryRow("SELECT id FROM user WHERE username=?", username).Scan(&id)
	if err2 == sql.ErrNoRows {
		//username available (found no rows)
		Info["errUsername"] = ""
	} else if err2 != nil {
		panic(err2.Error())
	} else {
		Info["errUsername"] = "Username already taken. Choose another one"
	}
	//checking for password & confirmPassword are same or not
	if password == confirmPassword {
		//passwords are same
		Info["errPassword"] = ""
	} else {
		Info["errPassword"] = "Password mismatched. Put cautiously"
	}

	//now do regitration
	if Info["errEmail"] != "" || Info["errUsername"] != "" || Info["errPassword"] != "" {
		if Info["errEmail"] != "" {
			Info["email"] = ""
		}
		if Info["errUsername"] != "" {
			Info["username"] = ""
		}
		tpl.ExecuteTemplate(w, "register.gohtml", Info)
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

	Info = map[string]interface{}{
		"pageTitle": "Page Not Found",
	}
	tpl.ExecuteTemplate(w, "404.gohtml", Info)
}
func Result(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	Info = map[string]interface{}{
		"pageTitle": "Verdict",
	}
	tpl.ExecuteTemplate(w, "result.gohtml", Info)
}
func Submission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	lastPage = ""

	path := r.URL.Path
	runes := []rune(path)
	OJpNum := string(runes[12:])

	need := "-"
	index := strings.Index(OJpNum, need)

	var OJ, pNum string
	runes = []rune(OJpNum)
	OJ = string(runes[0:3])
	if OJ == "计蒜客" {
		pNum = string(runes[4:])
	} else {
		OJ = string(runes[0:index])
		pNum = string(runes[index+1:])
	}

	session, _ := store.Get(r, "mysession")
	Info = map[string]interface{}{
		"username":    session.Values["username"],
		"password":    session.Values["password"],
		"isLogged":    session.Values["isLogin"],
		"pageTitle":   "Submission",
		"OJ":          OJ,
		"PNum":        pNum,
		"PName":       pTitle,
		"TimeLimit":   pTimeLimit,
		"MemoryLimit": pMemoryLimit,
	}
	tpl.ExecuteTemplate(w, "submission.gohtml", Info)
}
func GetLanguage(w http.ResponseWriter, r *http.Request) {
	getLanguage(w, r)
}
