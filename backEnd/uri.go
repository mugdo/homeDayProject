package backEnd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

var c = colly.NewCollector(
//colly.AllowedDomains("localhost"),
)

func getCsrfToken(apiURL string) string {
	var csrfToken string
	c.OnHTML("input[name='csrf_token']", func(e *colly.HTMLElement) {
		csrfToken = e.Attr("value")
	})
	c.Visit(apiURL)

	return csrfToken
}
func getFtaa(apiURL string, w http.ResponseWriter) {
	//var ftaa string
	client := &http.Client{
		Jar: cookieJar,
	}
	req, err := http.NewRequest("GET", "http://codeforces.com/enter", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	client = &http.Client{
		Jar: cookieJar,
	}
	req, err = http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err = client.Do(req)
	body, err = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	fmt.Fprintln(w, string(body))

	//return csrfToken
}

func Scrap(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println(resp)

		pURL := "https://www.urionlinejudge.com.br/judge/en/search?q=hard&for=problems"
		req, err := http.NewRequest("GET", pURL, nil)
		req.Header.Add("Content-Type", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		response, err := client.Do(req)
		defer response.Body.Close()

		document, err := goquery.NewDocumentFromReader(response.Body)
		checkErr(err)

		var uNumList,uNameList []string
		var uNum,uName string

		document.Find("td[class='id ']").Each(func(index int, linkStr *goquery.Selection) {
			uNum = linkStr.Find("a").Text()
			uNumList = append(uNumList, uNum)
		})
		document.Find("td[class='large']").Each(func(index int, linkStr *goquery.Selection) {
			uName = linkStr.Find("a").Text()
			uNameList = append(uNameList, uName)
		})

		type list struct{
			Num,Name string
		}
		var uList []list
		for i:=0;i<len(uNumList);i++{
			temp:=list{uNumList[i],uNameList[i]}

			uList = append(uList,temp)
		}
		b,err:=json.Marshal(uList)
		checkErr(err)

		for i := 0; i < len(uList); i++ {
			fmt.Println(i+1,uList[i].Num,uList[i].Name)
		}
		fmt.Println(string(b))
		fmt.Fprintln(w,string(b))

		fmt.Println("End")
	}
}
