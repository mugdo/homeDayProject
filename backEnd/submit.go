package backEnd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func rLogin(username, password, apiURL string) ([]byte, error) {
	postData := url.Values{}

	if apiURL == "https://toph.co/login" {
		postData.Set("handle", username)
	} else if apiURL == "https://vjudge.net/user/login" {
		postData.Set("username", username)
	}
	postData.Set("password", password)

	return rPOST(apiURL, postData)
}
func rSubmit(r *http.Request) ([]byte, error) {
	OJ := r.FormValue("OJ")
	pNum := r.FormValue("pNum")
	language := r.FormValue("language")
	source := r.FormValue("source")

	apiURL := "https://vjudge.net/problem/submit"

	postData := url.Values{
		"language": {language},
		"share":    {"0"},
		"source":   {source},
		"captcha":  {""},
		"oj":       {OJ},
		"probNum":  {pNum},
	}

	return rPOST(apiURL, postData)
}
func GetLanguage(w http.ResponseWriter, r *http.Request) {
	getLanguage(w, r)
}
func Submit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	session, _ := store.Get(r, "mysession")

	if r.Method != "POST" {
		lastPage = r.URL.Path

		if session.Values["isLogin"] == true {
			if isAccVerifed(r) {
				if r.URL.Path == "/submit" { //if only "/submit" is received (without OJ & pNum)
					Info["Username"] = session.Values["username"]
					Info["Password"] = session.Values["password"]
					Info["IsLogged"] = session.Values["isLogin"]
					Info["Lastpage"] = lastPage
					Info["PopUpCause"] = popUpCause
					Info["ErrorType"] = errorType
					Info["PageTitle"] = "Submission"
					Info["OJ"] = "AtCoder"
					Info["PNum"] = ""
					Info["PName"] = ""
					Info["TimeLimit"] = "-"
					Info["MemoryLimit"] = "-"

					tpl.ExecuteTemplate(w, "submit.gohtml", Info)
					popUpCause = ""
				} else if len(r.URL.Path) > 7 { //url is like "/submit/..."
					path := r.URL.Path
					runes := []rune(path)
					OJpNum := string(runes[8:])

					need := "-"
					index := strings.Index(OJpNum, need)
					OJ, pNum := "", ""

					if index == -1 { //url is not like this "/submit/OJ-pNum"
						pTitle, pTimeLimit, pMemoryLimit, pDesSrc, pOrigin = "", "", "", "", ""
						http.Redirect(w, r, "/submit", http.StatusSeeOther)
					} else { //url is something like this "/submit/OJ-pNum"
						runes = []rune(OJpNum)

						OJ = string(runes[0:index])
						pNum = string(runes[index+1:])

						if OJSet[OJ] == false || pNum == "" { //bad url, not OJ & pNum specified
							pTitle, pTimeLimit, pMemoryLimit, pDesSrc, pOrigin = "", "", "", "", ""
							http.Redirect(w, r, "/submit", http.StatusSeeOther)
						} else { // got something in OJ and pNum
							pDesSrc = "" //resetting for now
							//(Finding problem) Verifying that problem exist with this OJ & pNum
							findPResource(OJ, pNum)

							if pDesSrc == "" { //didn't get any problem
								pTitle, pTimeLimit, pMemoryLimit, pDesSrc, pOrigin = "", "", "", "", ""
								http.Redirect(w, r, "/submit", http.StatusSeeOther)
							} else { //got a problem with this OJ & pNum
								//checking whether problem submit allowed or not
								tempP, _ := pSearch(OJ, pNum, "", "", "20")
								tempList := getPList(tempP)

								allowSubmit := false
								if tempList.Data[0].AllowSubmit == true && tempList.Data[0].Status == 0 {
									allowSubmit = true
								}

								if allowSubmit == true {
									Info["Username"] = session.Values["username"]
									Info["Password"] = session.Values["password"]
									Info["IsLogged"] = session.Values["isLogin"]
									Info["Lastpage"] = lastPage
									Info["PopUpCause"] = popUpCause
									Info["ErrorType"] = errorType
									Info["PageTitle"] = "Submission"
									Info["OJ"] = OJ
									Info["PNum"] = pNum
									Info["PName"] = pTitle
									Info["TimeLimit"] = pTimeLimit
									Info["MemoryLimit"] = pMemoryLimit

									tpl.ExecuteTemplate(w, "submit.gohtml", Info)
									popUpCause = ""
								} else if allowSubmit == false {
									link := "/problemView/" + OJ + "-" + pNum
									http.Redirect(w, r, link, http.StatusSeeOther)
								}
							}
						}
					}
				}
			} else {
				lastPage = r.URL.Path
				popUpCause = "verifyRequired"
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		} else {
			lastPage = r.URL.Path
			popUpCause = "loginRequired"
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	} else if r.Method == "POST" {
		//do login first
		body, err := rLogin("ajudgebd", "aj199273", "https://vjudge.net/user/login")
		checkErr(err)
		result := string(body)

		if result == "success" { //if login success
			//submit the code
			body, err := rSubmit(r)
			checkErr(err)
			//fmt.Println("body", string(body))
			type result struct { //json reply gives either error or runID
				RunID int64  `json:"runId"`
				Error string `json:"error"`
				SubID string
				URL   string
			}
			var res result
			json.Unmarshal(body, &res)
			fmt.Println("1", res.Error, lastPage, Info["Res"])

			if res.Error != "" {
				//OJ := r.FormValue("OJ")
				//pNum := r.FormValue("pNum")
				//lastPage = "/submit/" + OJ + "-" + pNum
				//fmt.Println("lastpage=",lastPage)
				errorType = res.Error
				popUpCause = "submissionError"
				http.Redirect(w, r, lastPage, http.StatusSeeOther)
			} else {
				submissionID := strconv.FormatInt(res.RunID, 10)
				//sending submit id to frontend for getting the verdict with ajax call
				res.URL = "https://vjudge.net/solution/data/"
				res.SubID = submissionID

				Info["Username"] = session.Values["username"]
				Info["Password"] = session.Values["password"]
				Info["IsLogged"] = session.Values["isLogin"]
				Info["PageTitle"] = "Verdict"
				Info["Res"] = res

				http.Redirect(w, r, "/verdict", http.StatusSeeOther)
			}
		}
	}
}
func Verdict(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	tpl.ExecuteTemplate(w, "verdict.gohtml", Info)
}
