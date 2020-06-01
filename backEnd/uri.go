package backEnd

import (
	"encoding/json"
	"fmt"
	"html"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var method, csrfToken, tokenFields, tokenUnlocked = "", "", "", ""

func loginURI() string {
	pURL := "https://www.urionlinejudge.com.br/judge/en/login"
	req, err := http.NewRequest("GET", pURL, nil)
	req.Header.Add("Content-Type", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	response, err := client.Do(req)
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	checkErr(err)

	method, _ = document.Find("input[name='_method']").Attr("value")
	csrfToken, _ = document.Find("input[name='_csrfToken']").Attr("value")
	tokenFields, _ = document.Find("input[name='_Token[fields]']").Attr("value")
	tokenUnlocked, _ = document.Find("input[name='_Token[unlocked]']").Attr("value")

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
		max := 83 //there are total 83 page problem in URI
		pageNum := rand.Intn(max-min+1) + min

		pURL = "https://www.urionlinejudge.com.br/judge/en/problems/all?page=" + strconv.Itoa(pageNum)
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
func URIGet(pNum string) string {
	pURL := "https://www.urionlinejudge.com.br/repository/UOJ_" + pNum + "_en.html"
	var URIDes string

	if loginURI() == "success" {
		req, err := http.NewRequest("GET", pURL, nil)
		req.Header.Add("Content-Type", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		response, err := client.Do(req)
		defer response.Body.Close()

		document, err := goquery.NewDocumentFromReader(response.Body)
		checkErr(err)

		pTitle = document.Find("div[class='header']").Find("h1").Text()
		pTimeLimit = document.Find("div[class='header']").Find("strong").Text()
		pTimeLimit = strings.TrimPrefix(pTimeLimit, "Timelimit: ")

		URIDes, _ = document.Find("div[class='problem']").Html()
	}
	return URIDes
}
func URISubmit(w http.ResponseWriter, r *http.Request) {
	//gettimng form data
	pNum := r.FormValue("pNum")
	language := r.FormValue("language")
	source := strings.TrimSpace(r.FormValue("source"))

	//preparing source code for later comparing (ignoring carriage return /r)
	sourceEscape := html.EscapeString(source)
	var sourceMod string

	for i := 0; i < len(sourceEscape); i++ {
		if sourceEscape[i] != 13 { //ignoring (carriage return-/r) for comparing
			sourceMod += string(sourceEscape[i])
		}
	}

	if loginURI() == "success" { //first login to URI - if success
		postData := url.Values{
			"_method":          {method},
			"_csrfToken":       {csrfToken},
			"problem_id":       {pNum},
			"language_id":      {language},
			"template":         {"1"},
			"source_code":      {source},
			"_Token[fields]":   {tokenFields},
			"_Token[unlocked]": {tokenUnlocked},
		}

		//submitting to URI
		pURL := "https://www.urionlinejudge.com.br/judge/en/runs/add"
		req, _ := http.NewRequest("POST", pURL, strings.NewReader(postData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Add("Content-Length", strconv.Itoa(len(postData.Encode())))

		//setting up requset to prevent auto redirect
		client = &http.Client{
			Jar: cookieJar,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		response, _ := client.Do(req)
		defer response.Body.Close()

		//getting submission ID
		pURL = "https://www.urionlinejudge.com.br/judge/en/runs?problem_id=" + pNum + "&language_id=" + language
		req, _ = http.NewRequest("GET", pURL, nil)
		req.Header.Add("Content-Type", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		response, err := client.Do(req)
		defer response.Body.Close()

		// submitTime := response.Header["Date"]
		// runes := []rune(submitTime[0])
		// submitTimeOri := string(runes[17:25]) //00:00:00 from rfc1123TimeFormat
		// fmt.Println("SubmitTimeOri=", submitTimeOri)

		document, err := goquery.NewDocumentFromReader(response.Body)
		checkErr(err)

		//getting all submission ID from latest submission page
		var subID []string
		var ID string
		document.Find("td[class='id']").Each(func(index int, IDSeg *goquery.Selection) {
			ID = IDSeg.Find("a").Text()
			subID = append(subID, ID)
		})

		// var subPNum []string
		// var num string
		// document.Find("td[class='tiny']").Each(func(index int, numSeg *goquery.Selection) {
		// 	num = numSeg.Find("a").Text()
		// 	if len(num) == 4 {
		// 		subPNum = append(subPNum, num)
		// 	}
		// })
		// var subTime []string
		// var time string
		// document.Find("td[class='center']").Each(func(index int, timeSeg *goquery.Selection) {
		// 	time = timeSeg.Text()
		// 	var start int
		// 	var timeTemp string
		// 	if len(time) >= 52 {
		// 		for i := 21; i < len(time); i++ {
		// 			if time[i] == 32 { //if space found
		// 				start = i + 1
		// 				break
		// 			}
		// 		}
		// 		for i := start; i < len(time); i++ {
		// 			timeTemp += string(time[i])
		// 			if string(time[i]) == "M" { //if AM/PM found
		// 				break
		// 			}
		// 		}
		// 		var hour, minute, second, AmPm string
		// 		var next int
		// 		for i := 0; i < len(timeTemp); i++ {
		// 			if string(timeTemp[i]) == ":" {
		// 				next = i + 1
		// 				break
		// 			}
		// 			hour += string(timeTemp[i])
		// 		}
		// 		for i := next; i < len(timeTemp); i++ {
		// 			if string(timeTemp[i]) == ":" {
		// 				next = i + 1
		// 				break
		// 			}
		// 			minute += string(timeTemp[i])
		// 		}
		// 		for i := next; i < len(timeTemp); i++ {
		// 			if string(timeTemp[i]) == " " {
		// 				next = i + 1
		// 				break
		// 			}
		// 			second += string(timeTemp[i])
		// 			// tempSec, _ := strconv.Atoi(second)
		// 			// tempSec += 2 //submited delaying 2 sec
		// 			// if tempSec >= 60 {
		// 			// 	tempSec %= 60
		// 			// 	tempMin, _ := strconv.Atoi(minute)
		// 			// 	tempMin += 1
		// 			// 	if tempMin == 60 {
		// 			// 		tempMin %= 60
		// 			// 		if hour == "11" && AmPm == "A" {
		// 			// 			hour = "12"
		// 			// 			AmPm = "P"
		// 			// 		} else if hour == "12" && AmPm == "A" {
		// 			// 			hour = "1"
		// 			// 		} else if hour == "11" && AmPm == "P" {
		// 			// 			hour = "12"
		// 			// 			AmPm = "A"
		// 			// 		} else if hour == "12" && AmPm == "P" {
		// 			// 			hour = "1"
		// 			// 		}
		// 			// 	}
		// 			// 	minute = strconv.Itoa(tempMin)
		// 			// }
		// 			// second = strconv.Itoa(tempSec)
		// 		}
		// 		AmPm = string(timeTemp[next])
		// 		//brazil GMT-3 time to rfc1123 format
		// 		if hour == "9" && AmPm == "A" { //9AM->12PM=12
		// 			hour = "12"
		// 		} else if hour == "10" && AmPm == "A" { //10AM->1PM=13
		// 			hour = "13"
		// 		} else if hour == "11" && AmPm == "A" { //11AM->2PM=14
		// 			hour = "14"
		// 		} else if hour == "12" && AmPm == "A" { //12AM-3PM=03
		// 			hour = "03"
		// 		} else if AmPm == "A" { //(1-8AM)->(4-11AM)=(04-11)
		// 			temp, _ := strconv.Atoi(hour)
		// 			temp += 3
		// 			if hour != "7" && hour != "8" {
		// 				hour = "0"
		// 			}
		// 			hour += strconv.Itoa(temp)
		// 		} else if hour == "9" && AmPm == "P" { //9PM->12AM=00
		// 			hour = "00"
		// 		} else if hour == "10" && AmPm == "P" { //10PM->1AM=01
		// 			hour = "01"
		// 		} else if hour == "11" && AmPm == "P" { //11PM->2AM=02
		// 			hour = "02"
		// 		} else if hour == "12" && AmPm == "P" { //12PM->3PM=15
		// 			hour = "15"
		// 		} else if AmPm == "P" { //(1-8PM)->(4-11PM)=(16-23)
		// 			temp, _ := strconv.Atoi(hour)
		// 			temp += 3
		// 			temp += 12
		// 			hour = strconv.Itoa(temp)
		// 		}
		// 		timeFinal := hour + ":" + minute + ":" + second
		// 		subTime = append(subTime, timeFinal)
		// 	}
		// })

		//getting submission ID by matching original source code & submitted source code
		var actualSubID string
		for i := 0; i < len(subID); i++ {
			//getting submitted code one by one with collected subID
			req, _ := http.NewRequest("GET", "https://www.urionlinejudge.com.br/judge/en/runs/code/"+subID[i], nil)
			req.Header.Add("Content-Type", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
			response, err := client.Do(req)
			defer response.Body.Close()

			document, err := goquery.NewDocumentFromReader(response.Body)
			checkErr(err)

			subCode, _ := document.Find("pre[id='code']").Html()
			subCode = strings.TrimSpace(subCode)

			if subCode == sourceMod { //if original source code & submitted source code matched
				actualSubID = subID[i] //got actual/this submission ID
				break
			}
		}
		//got actual submit ID. Now getting verdict
		var result string
		for i := 1; i < 100; i++ { //if verdict got '- In queue -' then recheck will be done
			req, _ = http.NewRequest("GET", "https://www.urionlinejudge.com.br/judge/en/runs/code/"+actualSubID, nil)
			req.Header.Add("Content-Type", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
			response, err = client.Do(req)
			defer response.Body.Close()

			document, err = goquery.NewDocumentFromReader(response.Body)
			checkErr(err)

			result = document.Find("dd").Find("span").Text()
			result = strings.TrimSpace(result)

			if result != "- In queue -" && result != "Thinking..." {
				break
			}
		}
		fmt.Println("result=", actualSubID, result, "#")
	}
}
func Scrap(w http.ResponseWriter, r *http.Request) {
	rr, _ := rGET("https://www.urionlinejudge.com.br/repository/UOJ_1001_en.html")
	fmt.Println(string(rr))

	Info["Isrc"] = string(rr)
	tpl.ExecuteTemplate(w, "test2.gohtml", nil)

	fmt.Println("End")
}
