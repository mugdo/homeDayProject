package backEnd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	fmt.Println("Textarea:",textArea)

	b := []byte(textArea)
	fmt.Println("Byte:",b)

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

	fmt.Println("Trustable:",res.Trustable)

	for i:=1;i<len(res.Sections);i++ {
		fmt.Println("i:",i,"Title:",res.Sections[i].Title)
		
		fmt.Println("Format:",res.Sections[i].Value.Format)
		fmt.Println("Content:",res.Sections[i].Value.Content)
	}
}