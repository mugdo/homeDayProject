package backEnd

import (
	"fmt"
	"net/http"
)

var client = &http.Client{
	Jar: cookieJar,
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
