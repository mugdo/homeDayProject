package backEnd

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
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
							allowSubmit := false

							if OJ == "URI" {
								_ = URIGet(pNum)

								if pTitle == "" { //didn't get any problem
									pTitle, pTimeLimit, pMemoryLimit, pDesSrc, pOrigin = "", "", "", "", ""
									http.Redirect(w, r, "/submit", http.StatusSeeOther)
									return
								}
								//else got a problem with this OJ & pNum
								allowSubmit = true
							} else {
								//(Finding problem) Verifying that problem exist with this OJ & pNum
								findPResource(OJ, pNum)

								if pDesSrc == "" { //didn't get any problem
									pTitle, pTimeLimit, pMemoryLimit, pDesSrc, pOrigin = "", "", "", "", ""
									http.Redirect(w, r, "/submit", http.StatusSeeOther)
								} else { //got a problem with this OJ & pNum
									//checking whether problem submit allowed or not
									tempP, _ := pSearch(OJ, pNum, "", "20")
									tempList := getPList(tempP)

									if tempList.Data[0].AllowSubmit == true && tempList.Data[0].Status == 0 {
										allowSubmit = true
									}
								}
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
		OJ := r.FormValue("OJ")
		if OJ == "URI" {
			URISubmit(w, r)
			return
		}
		//if other OJ
		//do login first
		body, err := rLogin("ajudgebd", "aj199273", "https://vjudge.net/user/login")
		checkErr(err)
		result := string(body)

		if result == "success" { //if login success
			//submit the code
			body, err := rSubmit(r)
			checkErr(err)

			type result struct { //json reply gives either error or runID
				RunID int64  `json:"runId"`
				Error string `json:"error"`
			}
			var res result
			json.Unmarshal(body, &res)

			if res.Error != "" {
				errorType = res.Error
				popUpCause = "submissionError"
				http.Redirect(w, r, lastPage, http.StatusSeeOther)
			} else if res.RunID == 0 {
				errorType = res.Error
				popUpCause = "submissionErrorCustom"
				http.Redirect(w, r, lastPage, http.StatusSeeOther)
			} else {
				// inserting submission records to DB //
				DB := dbConn()
				defer DB.Close()

				//checking in DB for vSubID already exist or not
				var rows *sql.Rows
				rows, _ = DB.Query("SELECT username FROM submission WHERE vID=?", res.RunID)
				defer rows.Close()

				var uniqueUser bool = true
				var username string
				for rows.Next() {
					err = rows.Scan(&username)
					checkErr(err)

					if username == session.Values["username"] {
						uniqueUser = false
						break
					}
				}

				var id int
				if uniqueUser { //vSubID doesn't exist with this username
					insertQuery, err := DB.Prepare("INSERT INTO submission (username,OJ,pNum,language,submitTime,vID,sourceCode) VALUES (?,?,?,?,?,?,?)")
					checkErr(err)
					insertQuery.Exec(session.Values["username"], r.FormValue("OJ"), r.FormValue("pNum"), r.FormValue("language"), time.Now().Unix(), res.RunID, r.FormValue("source"))
				}
				//taking submission id
				_ = DB.QueryRow("SELECT id FROM submission WHERE username=? AND vID=?", session.Values["username"], res.RunID).Scan(&id)

				Info["Username"] = session.Values["username"]
				Info["Password"] = session.Values["password"]
				Info["IsLogged"] = session.Values["isLogin"]
				Info["PageTitle"] = "Result"
				Info["SubID"] = id //sending submit id to frontend for getting the verdict with ajax call
				Info["OJ"] = r.FormValue("OJ")
				Info["PNum"] = strings.TrimSpace(r.FormValue("pNum"))
				Info["Language"] = r.FormValue("language")
				Info["SourceCode"] = strings.TrimSpace(r.FormValue("source"))

				http.Redirect(w, r, "/result", http.StatusSeeOther)
			}
		}
	}
}
func Result(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	tpl.ExecuteTemplate(w, "result.gohtml", Info)
}
func Verdict(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	path := r.URL.Path
	runes := []rune(path)
	need := "="
	index := strings.Index(path, need)
	subID := string(runes[index+1:])

	DB := dbConn()
	defer DB.Close()

	//checking for submission status from DB
	var verdict, timeExec, memoryExec, vID string
	var submitTime int64
	_ = DB.QueryRow("SELECT verdict,timeExec,memoryExec,submitTime,vID FROM submission WHERE id=?", subID).Scan(&verdict, &timeExec, &memoryExec, &submitTime, &vID)

	if verdict == "Accepted" || verdict == "Wrong Answer" || verdict == "Compilation Error" || verdict == "Time Limit Exceeded" || verdict == "Memory Limit Exceeded" { //already data exist
		mapD := map[string]interface{}{
			"status":     verdict,
			"runtime":    timeExec,
			"memory":     memoryExec,
			"submitTime": time.Unix(submitTime, 0),
		}
		mapB, _ := json.Marshal(mapD)

		b := []byte(mapB)
		w.Write(b)
	} else { //data not exist. go for api call
		apiURL := "https://vjudge.net/solution/data/" + vID
		body, _ := rGET(apiURL)

		type Res struct {
			Status  string `json:"status"`
			Runtime int    `json:"runtime"`
			Memory  int    `json:"memory"`
		}
		var res Res
		json.Unmarshal(body, &res)

		updateQuery, err := DB.Prepare("UPDATE submission SET verdict=?,timeExec=?,memoryExec=? WHERE id=?")
		checkErr(err)
		updateQuery.Exec(res.Status, res.Runtime, res.Memory, subID)

		mapD := map[string]interface{}{
			"status":     res.Status,
			"runtime":    res.Runtime,
			"memory":     res.Memory,
			"submitTime": time.Unix(submitTime, 0),
		}
		mapB, _ := json.Marshal(mapD)

		b := []byte(mapB)
		w.Write(b)
	}
}
