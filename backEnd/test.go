package backEnd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
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
	fmt.Fprintf(w, "Hello Res")

	if loginURI() == "success" {
		req, _ := http.NewRequest("GET", "https://www.urionlinejudge.com.br/judge/en/runs/code/18359922", nil)
		req.Header.Add("Content-Type", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		response, _ := client.Do(req)
		defer response.Body.Close()
		document, _ := goquery.NewDocumentFromReader(response.Body)

		var result string
		document.Find("dd").Find("span").Each(func(index int, resSeg *goquery.Selection) {
			result = resSeg.Text()
		})

		result=strings.TrimSpace(result)
		fmt.Println("result=", result, "#")
	} else {
		fmt.Println("Nope")
	}

	//tpl.ExecuteTemplate(w, "test2.gohtml", Info)
}
func TestSub(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	//fmt.Fprintf(w, "Hello")
	//tpl.ExecuteTemplate(w, "test.gohtml", nil)

	http.Redirect(w, r, "/testPage", http.StatusSeeOther)
}

func goGet(w http.ResponseWriter) {
	var headings, row []string
	var rows [][]string

	data := `<html><body>
	<table>
	<thead>
		<tr>
			<th><a href="/judge/en/search?q=hard&amp;for=problems&amp;direction=asc&amp;sort=Problems.id">#</a></th>
			<th></th>
			<th class="left"><a href="/judge/en/search?q=hard&amp;for=problems&amp;direction=asc&amp;sort=Problems.name_en">Name</a></th>
			<th class="left">
							</th>
			<th><a href="/judge/en/search?q=hard&amp;for=problems&amp;direction=asc&amp;sort=Problems.solved">Solved</a></th>
			<th><a href="/judge/en/search?q=hard&amp;for=problems&amp;direction=asc&amp;sort=Problems.level">Level</a></th>
		</tr>
	</thead>
	<tbody>
					<tr class="impar">
								<td class="id ">
					<a href="/judge/en/problems/view/1728">1728</a>				</td>
				<td class="tiny">
					 
				</td>
				<td class="large">
					<a href="/judge/en/problems/view/1728">Hard to Believe, But True!</a>				</td>
				<td class="wide">
										</td>
				<td class="small">
					1,019				</td>
				<td class="tiny">3</td>
			</tr>

					<tr class="par">
								<td class="id ">
					<a href="/judge/en/problems/view/1260">1260</a>				</td>
				<td class="tiny">
					 
				</td>
				<td class="large">
					<a href="/judge/en/problems/view/1260">Hardwood Species</a>				</td>
				<td class="wide">
										</td>
				<td class="small">
					2,295				</td>
				<td class="tiny">5</td>
			</tr>

					<tr class="impar">
								<td class="id ">
					<a href="/judge/en/problems/view/1831">1831</a>				</td>
				<td class="tiny">
					 
				</td>
				<td class="large">
					<a href="/judge/en/problems/view/1831">Hard Day At Work</a>				</td>
				<td class="wide">
										</td>
				<td class="small">
					108				</td>
				<td class="tiny">8</td>
			</tr>

					<tr class="par">
								<td class="id ">
					<a href="/judge/en/problems/view/2986">2986</a>				</td>
				<td class="tiny">
					 
				</td>
				<td class="large">
					<a href="/judge/en/problems/view/2986">Not Everything is Strike Hard Version</a>				</td>
				<td class="wide">
										</td>
				<td class="small">
					67				</td>
				<td class="tiny">5</td>
			</tr>

					<tr class="impar">
								<td class="id ">
					<a href="/judge/en/problems/view/2702">2702</a>				</td>
				<td class="tiny">
					 
				</td>
				<td class="large">
					<a href="/judge/en/problems/view/2702">Hard Choice</a>				</td>
				<td class="wide">
										</td>
				<td class="small">
					3,180				</td>
				<td class="tiny">1</td>
			</tr>

					<tr class="par">
								<td class="id ">
					<a href="/judge/en/problems/view/2742">2742</a>				</td>
				<td class="tiny">
					 
				</td>
				<td class="large">
					<a href="/judge/en/problems/view/2742">Richards Multiverse</a>				</td>
				<td class="wide">
										</td>
				<td class="small">
					1,302				</td>
				<td class="tiny">5</td>
			</tr>

				
						<tr class="impar">
					<td colspan="7"></td>
				</tr>
								<tr class="par">
					<td colspan="7"></td>
				</tr>
								<tr class="impar">
					<td colspan="7"></td>
				</tr>
								<tr class="par">
					<td colspan="7"></td>
				</tr>
								<tr class="impar">
					<td colspan="7"></td>
				</tr>
								<tr class="par">
					<td colspan="7"></td>
				</tr>
								<tr class="impar">
					<td colspan="7"></td>
				</tr>
								<tr class="par">
					<td colspan="7"></td>
				</tr>
								<tr class="impar">
					<td colspan="7"></td>
				</tr>
								<tr class="par">
					<td colspan="7"></td>
				</tr>
								<tr class="impar">
					<td colspan="7"></td>
				</tr>
								<tr class="par">
					<td colspan="7"></td>
				</tr>
								<tr class="impar">
					<td colspan="7"></td>
				</tr>
								<tr class="par">
					<td colspan="7"></td>
				</tr>
					</tbody>
</table>
	</body>
	</html>
	`
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		fmt.Println("No url found")
		log.Fatal(err)
	}

	// Find each table
	doc.Find("table").Each(func(index int, tablehtml *goquery.Selection) {
		tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
			rowhtml.Find("th").Each(func(indexth int, tableheading *goquery.Selection) {
				headings = append(headings, tableheading.Text())
			})
			rowhtml.Find("td[class='id ']").Find("a").Each(func(indexth int, tablecell *goquery.Selection) {
				temp, _ := tablecell.Attr("href")
				row = append(row, temp)
			})
			rows = append(rows, row)
			row = nil
		})
	})
	fmt.Println("####### headings = ", len(headings), headings)
	fmt.Println("####### rows = ", len(rows), rows)

	for i := 0; i < len(rows); i++ {
		fmt.Println(rows[i], rows[i])
		fmt.Fprintln(w, rows[i], rows[i])
	}
}
