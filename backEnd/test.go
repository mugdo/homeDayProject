package backEnd

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"

	"golang.org/x/net/html"
)
 
func Body(doc *html.Node) (*html.Node, error) {
	var body *html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "div class=\"problem-statement\"" {
			body = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if body != nil {
		return body, nil
	}
	return nil, errors.New("Missing <body> in the node tree")
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

const htm = `<!DOCTYPE html>
<html>
<head>
    <title></title>
</head>
<body>
    body content
    <p>more content</p>
</body>
</html>`

func Test1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// getting problem description///////////////
	//pURL := "https://vjudge.net/problem/description/12222"
	pURL := "http://codeforces.com/problemset/problem/67/D"
	req, _ := http.NewRequest("GET", pURL, nil)
	//req.Header.Add("Content-Type", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	response, _ := client.Do(req)
	defer response.Body.Close()
	resp, _ := ioutil.ReadAll(response.Body)
	sresp := string(resp)

	re := regexp.MustCompile(`<div class="problem-statement".*?>(.*)</div>`)
	var des string
	submatchall := re.FindAllStringSubmatch(sresp, -1)
	for _, element := range submatchall {
		des = element[1]
		//fmt.Fprintln(w,element[1])
	}
	fmt.Println(des)
	data := map[string]interface{}{
		"Content": template.HTML(des),
	}
	tpl.ExecuteTemplate(w, "test1.gohtml", data)
}
func Test2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//fmt.Fprintf(w, "Hello Res")

	// test := r.FormValue("test")
	// fmt.Println(test)

	// data := map[string]interface{}{
	// 	"Test": test,
	// }

	tpl.ExecuteTemplate(w, "test2.gohtml", nil)
}
func TestSub(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	//fmt.Fprintf(w, "Hello")
	//tpl.ExecuteTemplate(w, "test.gohtml", nil)
	http.Redirect(w, r, "/testPage", http.StatusSeeOther)
}
