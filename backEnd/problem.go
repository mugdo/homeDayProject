package backEnd

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Problem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
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
		"Username":  session.Values["username"],
		"Password":  session.Values["password"],
		"IsLogged":  session.Values["isLogin"],
		"PList":     pList,
		"Found":     found,
		"OJ":        OJ,
		"PNum":      pNum,
		"PName":     pName,
		"PageTitle": "Problem",
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

	//Finding problem Title, Time limit, Memory limit, Description source(pDesSrc)
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
		"Username":    session.Values["username"],
		"Password":    session.Values["password"],
		"IsLogged":    session.Values["isLogin"],
		"PageTitle":   pTitle,
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
