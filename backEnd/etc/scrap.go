package backEnd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/gocolly/colly"
	//"golang.org/x/net/html"
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

	// req, _ := http.NewRequest("GET", "http://codeforces.com/problemset/problem/231/A", nil)
	// resp, err := client.Do(req)
	// defer resp.Body.Close()
	// doc, err := html.Parse(resp.Body)
	// if err != nil {
	// 	panic(err)
	// }
	// var f func(*html.Node)
	// f = func(n *html.Node) {
	// 	if n.Type == html.ElementNode && n.Data == "a" {
	// 		// Do something with n...
	// 		fmt.Fprintln(w, n.Attr)
	// 	}
	// 	for c := n.FirstChild; c != nil; c = c.NextSibling {
	// 		f(c)
	// 	}
	// }
	// f(doc)

	//csrfToken := getCsrfToken("http://codeforces.com/enter")
	getFtaa("http://codeforces.com/enter",w)
	// Find and visit all links

	// var c = colly.NewCollector(
	// //colly.AllowedDomains("localhost"),
	// )

	// c.OnHTML("meta[property='og:image']", func(e *colly.HTMLElement) {
	// 	fmt.Println("Inside")
	// 	//name := e.Text
	// 	csrfToken := e.Attr("content")
	// 	fmt.Println("ID:", csrfToken)
	// 	fmt.Println(csrfToken)
	// 	matched, err := regexp.MatchString("ftaa", csrfToken)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	//fmt.Println(matched)
	// 	if matched == true {
	// 		fmt.Println(csrfToken)
	// 	}
	// })

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL)
	// })

	// c.OnScraped(func(r *colly.Response) {
	// 	fmt.Println("Finished", r.Request.URL)
	// })

	// c.Visit("http://sta.codeforces.com/s/58652/js/ftaa.js")
}
