package backEnd

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func Problem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	OJ := r.FormValue("OJ")
	pNum := strings.TrimSpace(r.FormValue("pNum"))
	pName := strings.TrimSpace(r.FormValue("pName"))

	var pListFinal []Inner

	//searching problem from URI first
	var pURIb []byte
	if OJ == "" || OJ == "All" || OJ == "URI" {
		if pName != "" {
			pURIb = URISearch(pName)
		} else if pNum != "" {
			pURIb = URISearch(pNum)
		} else {
			if OJ == "URI" {
				pURIb = URISearch("URI Only")
			} else {
				rand.Seed(time.Now().UnixNano())
				min := 1001 //the first problem in URI is 1001
				max := 3100 //the last problem in URI is approx. 3100
				sQuery := rand.Intn(max-min+1) + min
				pURIb = URISearch(strconv.Itoa(sQuery))
			}
		}
		type pListURI struct {
			Num  string `json:"Num"`
			Name string `json:"Name"`
		}
		var pURI []pListURI
		json.Unmarshal(pURIb, &pURI)

		for i := 0; i < len(pURI); i++ {
			var temp Inner

			temp.AllowSubmit = true
			temp.ID = 0
			temp.IsFav = 0
			temp.OriginOJ = "URI"
			temp.OriginProb = pURI[i].Num
			temp.Status = 0
			temp.Title = pURI[i].Name
			temp.TriggerTime = 0

			pListFinal = append(pListFinal, temp)
		}
	}
	if OJ != "URI" {
		var length = "1000"
		body, err := pSearch(OJ, pNum, pName, length)
		checkErr(err)
		pList := getPList(body)

		for i := 0; i < len(pList.Data); i++ { //getting problem one by one
			if OJSet[pList.Data[i].OriginOJ] { //if problem come from desired OJ
				pListFinal = append(pListFinal, pList.Data[i])
			}
		}
	}

	found := true
	if len(pListFinal) == 0 { //if no problem list found
		found = false
	}

	lastPage = "problem"
	session, _ := store.Get(r, "mysession")

	Info["Username"] = session.Values["username"]
	Info["Password"] = session.Values["password"]
	Info["IsLogged"] = session.Values["isLogin"]
	Info["PList"] = pListFinal
	Info["Found"] = found
	Info["OJ"] = OJ
	Info["PNum"] = pNum
	Info["PName"] = pName
	Info["PageTitle"] = "Problem"

	tpl.ExecuteTemplate(w, "problem.gohtml", Info)
}
func ProblemView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	path := r.URL.Path
	runes := []rune(path)
	OJpNum := string(runes[13:])

	need := "-"
	index := strings.Index(OJpNum, need)
	OJ, pNum := "", ""

	if index == -1 { //url is not like this "/peoblemview/OJ-pNum"
		errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
	} else {
		runes = []rune(OJpNum)
		OJ = string(runes[0:index])
		pNum = string(runes[index+1:])

		if OJSet[OJ] == false || pNum == "" { //bad url, not OJ & pNum specified
			errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
		} else { // got something in OJ and pNum
			allowSubmit := false //just for declaring, need later
			var uvaSegment string
			problem := []map[string]interface{}{}
			URIProblem := map[string]interface{}{}

			if OJ == "URI" {
				URIProblem["Des"] = template.HTML(URIGet(pNum))
				allowSubmit = true

				if pTitle == "" {
					errorPage(w, http.StatusBadRequest)
					return
				}
			} else {
				//Finding problem Title, Time limit, Memory limit, Description source(pDesSrc)
				findPResource(OJ, pNum)

				if pDesSrc == "" { //didn't get any problem
					errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
				} else { //got a problem and resources
					// getting problem description
					pURL := "https://vjudge.net" + pDesSrc //value of pDesSrc came from findPResource(OJ, pNum) function
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
					tempP, _ := pSearch(OJ, pNum, "", "20")
					tempList := getPList(tempP)

					if tempList.Data[0].AllowSubmit == true && tempList.Data[0].Status == 0 {
						allowSubmit = true
					}

					//for uva pdf description
					if OJ == "UVA" {
						temp, _ := strconv.Atoi(pNum)

						IntSegment := temp / 100
						uvaSegment = strconv.Itoa(IntSegment)
					}
				}
			}

			lastPage = r.URL.Path
			session, _ := store.Get(r, "mysession")

			Info["Username"] = session.Values["username"]
			Info["Password"] = session.Values["password"]
			Info["IsLogged"] = session.Values["isLogin"]
			Info["PageTitle"] = pTitle
			Info["OJ"] = OJ
			Info["PNum"] = pNum
			Info["AllowSubmit"] = allowSubmit
			Info["UvaSegment"] = uvaSegment
			Info["PName"] = pTitle
			Info["TimeLimit"] = pTimeLimit
			Info["MemoryLimit"] = pMemoryLimit
			Info["Problem"] = problem
			Info["URIProblem"] = URIProblem

			tpl.ExecuteTemplate(w, "problemView.gohtml", Info)
		}
	}
}
func Origin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	path := r.URL.Path
	runes := []rune(path)
	OJpNum := string(runes[8:])

	need := "-"
	index := strings.Index(OJpNum, need)
	OJ, pNum := "", ""

	if index == -1 { //url is not like this "/origin/OJ-pNum"
		errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
	} else {
		runes = []rune(OJpNum)

		OJ = string(runes[0:index])
		pNum = string(runes[index+1:])

		if OJSet[OJ] == false || pNum == "" { //bad url, not OJ & pNum specified
			errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
		} else { // got something in OJ and pNum
			if OJ == "URI" {
				pOrigin = "https://www.urionlinejudge.com.br/judge/en/problems/view/" + pNum

				req, _ := http.NewRequest("GET", pOrigin, nil)
				req.Header.Add("Content-Type", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
				response, _ := client.Do(req)
				defer response.Body.Close()
				res, _ := ioutil.ReadAll(response.Body)
				sRes := string(res)

				need := "iframe"
				index := strings.Index(sRes, need)

				if index == -1 { //didn't get any problem
					errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
				} else {
					//redirecting to the origin site
					http.Redirect(w, r, pOrigin, http.StatusSeeOther)
				}
			} else {
				//Finding origin
				pURL := "https://vjudge.net/problem/" + OJ + "-" + pNum

				req, err := http.NewRequest("GET", pURL, nil)
				req.Header.Add("Content-Type", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
				response, err := client.Do(req)
				defer response.Body.Close()

				document, err := goquery.NewDocumentFromReader(response.Body)
				checkErr(err)

				pOrigin, _ = document.Find("span[class='origin']").Find("a").Attr("href")

				if pOrigin == "" { //didn't get any problem
					errorPage(w, http.StatusBadRequest) //http.StatusBadRequest = 400
				} else { //got a problem and origin
					//getting origin link
					pOrigin = getOriginLink("https://vjudge.net" + pOrigin)

					//redirecting to the origin site
					http.Redirect(w, r, pOrigin, http.StatusSeeOther)
				}
			}
		}
	}
}
