package backEnd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"golang.org/x/net/html"
)

func Test1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "Hello")
	tpl.ExecuteTemplate(w, "test1.gohtml", nil)
}
func Test2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "Hello Res")

	test := r.FormValue("test")
	fmt.Println(test)

	data := map[string]interface{}{
		"Test": test,
	}

	tpl.ExecuteTemplate(w, "testPage.gohtml", data)
}
func TestSub(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	//fmt.Fprintf(w, "Hello")
	//tpl.ExecuteTemplate(w, "test.gohtml", nil)
	http.Redirect(w,r,"/testPage",http.StatusSeeOther)
}

func Body(doc *html.Node) (*html.Node, error) {
    var body *html.Node
    var crawler func(*html.Node)
    crawler = func(node *html.Node) {
        if node.Type == html.ElementNode && node.Data == "body" {
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

func main() {
    doc, _ := html.Parse(strings.NewReader(htm))
    bn, err := Body(doc)
    if err != nil {
        return
    }
    body := renderNode(bn)
    fmt.Println(body)
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