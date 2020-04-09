package backEnd

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"github.com/PuerkitoBio/goquery"
)

func Des(w http.ResponseWriter, r *http.Request){
	 url := "https://vjudge.net/problem/description/39028"
	// body, _ := rGET(url)
	// fmt.Printf(string(body))

	client := &http.Client{
		Jar: cookieJar,
	}
	pURL := url
	response, err := client.Get(pURL)
	if err != nil {
		log.Fatalln("Error fetching response. ", err)
	}
	defer response.Body.Close()
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}
	textArea := document.Find("textarea").Text()
	//fmt.Println("Textarea:",textArea)

	b := []byte(textArea)
	//fmt.Println("Byte:",b)

	type Inner2 struct{
		Format	string	`json:"format"`
		Content	string	`json:"content"`
	}
	type Inner struct{
		Title	string `json:"title"`
		Value	Inner2 `json:"value"`
	}
	type Res struct{
		Trustable	bool	`json:"trustable"`
		Sections	[]Inner	`json:"sections"`
	}
	var res Res
	json.Unmarshal(b, &res)


	//Eliminating Default CSS on Example-Input-Output
	need := "<script"

	index := strings.Index(res.Sections[0].Value.Content,need)
	fmt.Println("got Script")
	fmt.Println(index)

	runes := []rune(res.Sections[0].Value.Content)
	res.Sections[0].Value.Content = string(runes[index:])

	//fmt.Println("Trustable:",res.Trustable)
	// for i:=0;i<len(res.Sections);i++ {
	// 	fmt.Println("i:",i,"Title:",res.Sections[i].Title)
		
	// 	fmt.Println("Format:",res.Sections[i].Value.Format)
	// 	fmt.Println("Content:",res.Sections[i].Value.Content)
	// }
	session, _ := store.Get(r, "mysession")

	data := map[string]interface{}{
		"username"	: session.Values["username"],
		"password"	: session.Values["password"],
		"isLogged"	: session.Values["isLogin"],
		"pageTitle"	: "",

		"Oj"			: "CF",
		"PNum"			: "1111",
		"PName"			: "title",
		"TimeLimit"		: "timeLimit",
		"MemoryLimit"	: "memoryLimit",
		
		"PDes"		: template.HTML(res.Sections[0].Value.Content),
		"PInputH"	: template.HTML(res.Sections[1].Title),
		"PInputD"	: template.HTML(res.Sections[1].Value.Content),
		"POutputH"	: template.HTML(res.Sections[2].Title),
		"POutputD"	: template.HTML(res.Sections[2].Value.Content),
		"PExampleH"	: template.HTML(res.Sections[3].Title),
		"PExampleD"	: template.HTML(res.Sections[3].Value.Content),
		"PNoteH"	: template.HTML(res.Sections[4].Title),
		"PNoteD"	: template.HTML(res.Sections[4].Value.Content),
	}

	tpl.ExecuteTemplate(w, "problemView.gohtml", data)
}