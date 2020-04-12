package backEnd

import (
	"fmt"
	"net/http"
	"github.com/gocolly/colly"
)

var client = &http.Client{
	Jar: cookieJar,
}

func Scrap(w http.ResponseWriter, r *http.Request) {
	c := colly.NewCollector(
	//colly.AllowedDomains("localhost"),
	)

	// Find and visit all links
	c.OnHTML("iframe[src]", func(e *colly.HTMLElement) {
		desID := e.Attr("src")
		fmt.Println("Description ID:",desID)

		// matched, err := regexp.MatchString("problem", res)
		// if err != nil {
		// 	panic(err)
		// }
		// //fmt.Println(matched)
		// if matched == true {
		// 	fmt.Println(res)
		// }
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit("https://vjudge.net/problem/CodeForces-341C")
}
