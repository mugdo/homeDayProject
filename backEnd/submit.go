package backEnd

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
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
func Submit(w http.ResponseWriter, r *http.Request) {
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
		if session.Values["isLogin"] == nil {
			session.Values["isLogin"] = false
		}

		Info = map[string]interface{}{
			"username":  session.Values["username"],
			"password":  session.Values["password"],
			"isLogged":  session.Values["isLogin"],
			"pageTitle": "Verdict",
			"Res":       res,
		}

		http.Redirect(w, r, "/verdict", http.StatusSeeOther)
		//tpl.ExecuteTemplate(w, "result.gohtml", data)
	}
}
func Verdict(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "verdict.gohtml", Info)
}
