package backEnd

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func loginURI() string {
	pURL := "https://www.urionlinejudge.com.br/judge/en/login"
	req, err := http.NewRequest("GET", pURL, nil)
	req.Header.Add("Content-Type", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	response, err := client.Do(req)
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	checkErr(err)

	method, _ := document.Find("input[name='_method']").Attr("value")
	csrfToken, _ := document.Find("input[name='_csrfToken']").Attr("value")
	tokenFields, _ := document.Find("input[name='_Token[fields]']").Attr("value")
	tokenUnlocked, _ := document.Find("input[name='_Token[unlocked]']").Attr("value")

	//fmt.Println(method, csrfToken, tokenFields, tokenUnlocked)

	postData := url.Values{
		"_method":          {method},
		"_csrfToken":       {csrfToken},
		"email":            {"ajudge.bd@gmail.com"},
		"password":         {"aj199273"},
		"remember_me":      {"0"},
		"_Token[fields]":   {tokenFields},
		"_Token[unlocked]": {tokenUnlocked},
	}

	pURL = "https://www.urionlinejudge.com.br/judge/en/login"
	req, err = http.NewRequest("POST", pURL, strings.NewReader(postData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	response, err = client.Do(req)
	defer response.Body.Close()

	document, err = goquery.NewDocumentFromReader(response.Body)
	checkErr(err)

	resp := document.Find("div[class='h-user']").Find("i").Text()

	if resp == "ajudge.bd@gmail.com" {
		return "success"
	} else {
		return "failed"
	}
}
func URISearch(sQuery string) []byte {
	var res []byte
	var pURL string

	if sQuery == "URI Only" {
		rand.Seed(time.Now().UnixNano())
		min := 1
		max := 83
		pageNum := rand.Intn(max-min+1) + min

		pURL = "https://www.urionlinejudge.com.br/judge/en/problems/all?page="+strconv.Itoa(pageNum)
	} else {
		pURL = "https://www.urionlinejudge.com.br/judge/en/search?q=" + sQuery + "&for=problems"
	}

	if loginURI() == "success" {
		req, err := http.NewRequest("GET", pURL, nil)
		req.Header.Add("Content-Type", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		response, err := client.Do(req)
		defer response.Body.Close()

		document, err := goquery.NewDocumentFromReader(response.Body)
		checkErr(err)

		var uNumList, uNameList []string
		var uNum, uName string

		document.Find("td[class='id']").Each(func(index int, linkStr *goquery.Selection) {
			uNum = linkStr.Find("a").Text()
			uNumList = append(uNumList, uNum)
		})
		document.Find("td[class='id ']").Each(func(index int, linkStr *goquery.Selection) {
			uNum = linkStr.Find("a").Text()
			uNumList = append(uNumList, uNum)
		})
		document.Find("td[class='id tour-step-1001']").Each(func(index int, linkStr *goquery.Selection) {
			uNum = linkStr.Find("a").Text()
			uNumList = append(uNumList, uNum)
		})
		document.Find("td[class='large']").Each(func(index int, linkStr *goquery.Selection) {
			uName = linkStr.Find("a").Text()
			uNameList = append(uNameList, uName)
		})

		type list struct {
			Num, Name string
		}
		var uList []list
		for i := 0; i < len(uNumList); i++ {
			temp := list{uNumList[i], uNameList[i]}

			uList = append(uList, temp)
		}

		res, _ = json.Marshal(uList)
	}

	return res
}

func Scrap(w http.ResponseWriter, r *http.Request) {
	rr, _ := rGET("https://www.urionlinejudge.com.br/repository/UOJ_1001_en.html")
	fmt.Println(string(rr))

	Info["Isrc"] = string(rr)
	tpl.ExecuteTemplate(w, "test2.gohtml", nil)

	fmt.Println("End")
}
