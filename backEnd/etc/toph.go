package backEnd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"github.com/PuerkitoBio/goquery"
)

type App struct {
	Client *http.Client
}

type AuthenticityToken struct {
	Token string
}

var tokenId string

func findToken(body string) (string) {
	//fmt.Println(body)

	need := "tokenId"

	index := strings.Index(body,need)
	fmt.Println("got")
	fmt.Println(index)

	runes := []rune(body)
	sb := string(runes[index+9:index+41])
	fmt.Println(sb)
	fmt.Println(len(sb))

	return sb;
}

func (app *App) getToken() AuthenticityToken {
	loginURL := "http://toph.co"
	client := app.Client

	response, err := client.Get(loginURL)

	if err != nil {
		log.Fatalln("Error fetching response. ", err)
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	//resp, err := ioutil.ReadAll(response.Body)
	//body := string(resp)

	token := document.Find("script").Text()

	token = findToken(token)

	fmt.Println("before token")
	authenticityToken := AuthenticityToken{
		Token: token,
	}

	return authenticityToken
}
func (app *App) login() {
	authenticityToken := app.getToken()
	tokenId = authenticityToken.Token
}
func pSubmit() ([]byte, error) {
	//fmt.Println("In rSubmit")

	apiURL := "https://toph.co/api/problems/5be6911043ca9d0001941344/submissions"

	code := `#include <iostream>	
	using namespace std;
	
	int main() {
		int n;
		cin>>n;
		
		cout<<n<<endl;
		
		return 0;
	}`

	data := url.Values{}
	data.Set("languageId", "5d828f1e9d55050001e97ee4")
	data.Set("source", code)

	u, _ := url.ParseRequestURI(apiURL)
	urlStr := u.String()

	client := &http.Client{
		Jar: cookieJar,
	}

	tk := "Token "+tokenId
	fmt.Println(tk)
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	req.Header.Add("Authorization", tk)
	req.Header.Add("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundaryCvE4v3gnlpsLdeEg")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

func Toph(w http.ResponseWriter, r *http.Request) {
	//login to toph
	_, err := rLogin("ajudge.bd", "aj199273", "https://toph.co/login")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	//result := string(body)
	//fmt.Println(result)
	fmt.Println("login done")

	//getting token
	app := App{
		Client: &http.Client{Jar: cookieJar},
	}
	app.login()
	//fmt.Println(tokenId)

	//submitting
	fmt.Println("Now submitting")
	body, err := pSubmit()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(body))
	fmt.Println("ENDDDDDDDDDDD")
}
