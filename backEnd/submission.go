package backEnd

import (
	"encoding/json"
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
func Submission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	lastPage = "submission"
	session, _ := store.Get(r, "mysession")
	if r.URL.Path == "/submission" {
		Info = map[string]interface{}{
			"Username":    session.Values["username"],
			"Password":    session.Values["password"],
			"IsLogged":    session.Values["isLogin"],
			"Lastpage":    lastPage,
			"PageTitle":   "Submission",
			"OJ":          "",
			"PNum":        "",
			"PName":       "",
			"TimeLimit":   "-",
			"MemoryLimit": "-",
		}
		tpl.ExecuteTemplate(w, "submission.gohtml", Info)
	} else {
		path := r.URL.Path
		runes := []rune(path)
		var OJpNum, OJ, pNum string

		if len(path) > 11 { //if only "/submission" is received (without OJ & pNum)
			OJpNum = string(runes[12:])

			need := "-"
			index := strings.Index(OJpNum, need)

			runes = []rune(OJpNum)
			OJ = string(runes[0:3])
			if OJ == "计蒜客" {
				pNum = string(runes[4:])
			} else {
				OJ = string(runes[0:index])
				pNum = string(runes[index+1:])
			}
		}

		Info = map[string]interface{}{
			"Username":    session.Values["username"],
			"Password":    session.Values["password"],
			"IsLogged":    session.Values["isLogin"],
			"Lastpage":    lastPage,
			"PageTitle":   "Submission",
			"OJ":          OJ,
			"PNum":        pNum,
			"PName":       pTitle,
			"TimeLimit":   pTimeLimit,
			"MemoryLimit": pMemoryLimit,
		}
		tpl.ExecuteTemplate(w, "submission.gohtml", Info)
	}
}
func GetLanguage(w http.ResponseWriter, r *http.Request) {
	getLanguage(w, r)
}
func Submit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	if r.Method != "POST" {
		http.Redirect(w, r, "/submission", http.StatusSeeOther)
		return
	}
	//do login first
	body, err := rLogin("ajudgebd", "aj199273", "https://vjudge.net/user/login")
	checkErr(err)
	result := string(body)

	if result == "success" {
		//submit the code
		body, err := rSubmit(r)
		checkErr(err)

		type result struct {
			RunID int64 `json:"runId"`
			SubID string
			URL   string
		}
		var res result
		json.Unmarshal(body, &res)

		submissionID := strconv.FormatInt(res.RunID, 10)

		//sending submission id to frontend for getting the verdict with ajax call
		res.URL = "https://vjudge.net/solution/data/"
		res.SubID = submissionID

		session, _ := store.Get(r, "mysession")
		Info = map[string]interface{}{
			"Username":  session.Values["username"],
			"Password":  session.Values["password"],
			"IsLogged":  session.Values["isLogin"],
			"PageTitle": "Verdict",
			"Res":       res,
		}

		http.Redirect(w, r, "/verdict", http.StatusSeeOther)
	}
}
func Verdict(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	tpl.ExecuteTemplate(w, "verdict.gohtml", Info)
}
