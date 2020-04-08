package giveservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func rLogin(username, password, apiURL string)([]byte, error){
	//fmt.Println("In rLogin")

	//apiURL := "https://vjudge.net/user/login"
	data := url.Values{}

	if apiURL == "https://toph.co/login" {
		data.Set("handle", username)
	} else if apiURL == "https://vjudge.net/user/login" {
		data.Set("username", username)
	}
	
	data.Set("password", password)

	return rPOST(apiURL, data)
}

func rSubmit()([]byte, error){
	//fmt.Println("In rSubmit")
	
	apiURL := "https://vjudge.net/problem/submit"

	code := `#include<iostream>
	using namespace std;
	int main()
	{
		int test;
		cin>>test;
		for(int t=1;t<=test;t++)
		{
			int a,o;
			cin>>a>>o;

			cout<<"Case "<<t<<": "<<a+o<<endl;
		}
		return 0;
	}
			`

	data := url.Values{}
	data.Set("language", "C++")
	data.Set("share", "0")
	data.Set("source", code)
	data.Set("captcha", "")
	data.Set("oj", "LightOJ")
	data.Set("probNum", "1000")

	return rPOST(apiURL, data)
}

func verdict(w http.ResponseWriter, r *http.Request, subID string) {
	//fmt.Println("In rVerdict")

	type result struct {
		Memory            int    `json:"memory"`
		StatusType        int    `json:"statusType"`
		Author            string `json:"author"`
		Length            int    `json:"length"`
		Runtime           int    `json:"runtime"`
		Language          string `json:"language"`
		StatusCanonical   string `json:"statusCanonical"`
		AuthorID          int64  `json:"authorId"`
		LanguageCanonical string `json:"languageCanonical"`
		CodeAccessInfo    string `json:"codeAccessInfo"`
		SubmitTime        int64  `json:"submitTime"`
		IsOpen            int    `json:"isOpen"`
		Processing        bool   `json:"processing"`
		RunID             int64  `json:"runId"`
		Oj                string `json:"oj"`
		RemoteRunID       string `json:"remoteRunId"`
		ProbNum           string `json:"probNum"`
		Status            string `json:"status"`
	}
	var res result

	urlStr := "https://vjudge.net/solution/data/" + subID

	for i := 1; i <= 50; i++ {
		body, err := rGET(urlStr)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(body, &res)

		fmt.Println(res.Status)

		if res.Status == "Accepted" {
			break
		}

		time.Sleep(2 * time.Second)
	}

	fmt.Println("===Final Details===")
	fmt.Println("Memory:", res.Memory)
	fmt.Println("StatusType:", res.StatusType)
	fmt.Println("Author:", res.Author)
	fmt.Println("Length:", res.Length)
	fmt.Println("Runtime:", res.Runtime)
	fmt.Println("Language:", res.Language)
	fmt.Println("StatusCanon:", res.StatusCanonical)
	fmt.Println("AuthorID:", res.AuthorID)
	fmt.Println("LanguageCan:", res.Language)
	fmt.Println("CodeAccessI:", res.CodeAccessInfo)
	fmt.Println("SubmitTime:", res.SubmitTime)
	fmt.Println("IsOpen:", res.IsOpen)
	fmt.Println("Processing:", res.Processing)
	fmt.Println("RunID :", res.RunID)
	fmt.Println("Oj:", res.Oj)
	fmt.Println("RemoteRunID:", res.RemoteRunID)
	fmt.Println("ProbNum:", res.ProbNum)
	fmt.Println("Status:", res.Status)

	//tpl.ExecuteTemplate(w, "result.html", res)
}

func Submit(w http.ResponseWriter, r *http.Request) {
	//do login first
	body, err := rLogin("ajudgebd", "aj199273", "https://vjudge.net/user/login")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	result := string(body)
	fmt.Println(result)

	if result == "success" {
		//submit the code
		body, err := rSubmit()
		if err != nil {
			fmt.Println(err)
		}

		type result struct {
			RunID int64 `json:"runId"`
			SubID string
			URL string
		}
		var res result
		json.Unmarshal(body, &res)

		submissionID := strconv.FormatInt(res.RunID, 10)
		//verdict(w, r, submissionID)

		//sending submission id to frontend for getting the verdict with ajax call
		res.URL = "https://vjudge.net/solution/data/"
		res.SubID = submissionID
		fmt.Println(res.URL,res.SubID)

		session, _ := store.Get(r, "mysession")

		if session.Values["isLogin"] == nil {
			session.Values["isLogin"] = false
		}
	
		data := map[string]interface{}{
			"username"	: session.Values["username"],
			"password"	: session.Values["password"],
			"isLogged"	: session.Values["isLogin"],
			"pageTitle"	: "Verdict",
			"Res"		: res,
		}

		tpl.ExecuteTemplate(w, "result.gohtml", data)
	}
}
