package backEnd

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	//"github.com/gocolly/colly"
)

var client = &http.Client{
	Jar: cookieJar,
}

func getTimeLimit(body string) (a, b string) {
	//fmt.Println(body)

	need := "ms"

	index := strings.Index(body, need)
	fmt.Println("got Time")
	fmt.Println(index)

	runes := []rune(body)
	sb := "-"
	if index != -1 {
		sb = string(runes[0 : index+2])
	}
	fmt.Println("Time Limit:", sb)

	// if oj == "CodeChef" {
	// 	need = "B"
	// }
	need = "B"
	index2 := strings.Index(body, need)
	fmt.Println("got Memory")
	fmt.Println(index2)

	//runes = []rune(body)
	sbm := "-"
	if index2 != -1 {
		sbm = string(runes[index+2 : index2+1])
	}
	fmt.Println("Memory Limit:", sbm)

	return sb, sbm
}

func getTitle(oj, pNum string) (a, b, c, d string) {
	fmt.Println("In GetTitle")
	client := &http.Client{
		Jar: cookieJar,
	}
	pURL := "https://vjudge.net/problem/" + oj + "-" + pNum

	response, err := client.Get(pURL)
	if err != nil {
		log.Fatalln("Error fetching response. ", err)
	}
	defer response.Body.Close()
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	timeLimitText := document.Find("dl[class='card row']").Find("dd[class='col-sm-7']").Text()
	timeLimit, memoryLimit := getTimeLimit(timeLimitText)

	title := document.Find("div[id='prob-title']").Find("h2").Text()
	src, iv := document.Find("iframe").Attr("src")

	fmt.Println("src=", src, "iv", iv)
	return title, timeLimit, memoryLimit, src
}

func Scrap(w http.ResponseWriter, r *http.Request) {

	fmt.Println("IN scrap")

	title, timeLimit, memoryLimit, dNum := getTitle("CodeForces", "831A")

	fmt.Println("After getting title ", title)
	fmt.Println("Finished Scrap")

	session, _ := store.Get(r, "mysession")
	data := map[string]interface{}{
		"username":    session.Values["username"],
		"password":    session.Values["password"],
		"isLogged":    session.Values["isLogin"],
		"Oj":          "CF",
		"PNum":        "1111",
		"PName":       title,
		"TimeLimit":   timeLimit,
		"MemoryLimit": memoryLimit,
		"pageTitle":   title,
	}

	// getting problem description
	url := "https://vjudge.net" + dNum
	res, _ := rGET(url)

	fmt.Printf(string(res))

	tpl.ExecuteTemplate(w, "problemView.gohtml", data)

	// c := colly.NewCollector(
	// //colly.AllowedDomains("localhost"),
	// )

	// // Find and visit all links
	// c.OnHTML("iframe[src]", func(e *colly.HTMLElement) {
	// 	desID := e.Attr("src")
	// 	fmt.Println("Description ID:",desID)

	// 	// matched, err := regexp.MatchString("problem", res)
	// 	// if err != nil {
	// 	// 	panic(err)
	// 	// }
	// 	// //fmt.Println(matched)
	// 	// if matched == true {
	// 	// 	fmt.Println(res)
	// 	// }
	// })

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL)
	// })

	// c.OnScraped(func(r *colly.Response) {
	// 	fmt.Println("Finished", r.Request.URL)
	// })

	// c.Visit("https://vjudge.net/problem/CodeForces-341C")
}
