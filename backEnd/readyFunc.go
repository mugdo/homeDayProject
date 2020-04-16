package backEnd

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

var cookieJar, _ = cookiejar.New(nil)
var client = &http.Client{
	Jar: cookieJar,
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func rGET(urlStr string) ([]byte, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}
func rPOST(urlStr string, data url.Values) ([]byte, error) {
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

func pSearch(OJ, pNum, pName string) ([]byte, error) {
	apiURL := "https://vjudge.net/problem/data"

	data := url.Values{
		"draw":                      {"1"},
		"columns[0][data]":          {"0"},
		"columns[1][data]":          {"1"},
		"columns[2][data]":          {"2"},
		"columns[3][data]":          {"3"},
		"columns[4][data]":          {"4"},
		"columns[5][data]":          {"5"},
		"columns[6][data]":          {"6"},
		"columns[0][name]":          {""},
		"columns[1][name]":          {""},
		"columns[2][name]":          {""},
		"columns[3][name]":          {""},
		"columns[4][name]":          {""},
		"columns[5][name]":          {""},
		"columns[6][name]":          {""},
		"columns[0][searchable]":    {"true"},
		"columns[1][searchable]":    {"true"},
		"columns[2][searchable]":    {"false"},
		"columns[3][searchable]":    {"true"},
		"columns[4][searchable]":    {"true"},
		"columns[5][searchable]":    {"true"},
		"columns[6][searchable]":    {"true"},
		"columns[0][orderable]":     {"false"},
		"columns[1][orderable]":     {"true"},
		"columns[2][orderable]":     {"false"},
		"columns[3][orderable]":     {"true"},
		"columns[4][orderable]":     {"true"},
		"columns[5][orderable]":     {"true"},
		"columns[6][orderable]":     {"false"},
		"columns[0][search][value]": {""},
		"columns[1][search][value]": {""},
		"columns[2][search][value]": {""},
		"columns[3][search][value]": {""},
		"columns[4][search][value]": {""},
		"columns[5][search][value]": {""},
		"columns[6][search][value]": {""},
		"columns[0][search][regex]": {"false"},
		"columns[1][search][regex]": {"false"},
		"columns[2][search][regex]": {"false"},
		"columns[3][search][regex]": {"false"},
		"columns[4][search][regex]": {"false"},
		"columns[5][search][regex]": {"false"},
		"columns[6][search][regex]": {"false"},
		"order[0][column]":          {"5"},
		"order[0][dir]":             {"desc"},
		"search[value]":             {""},
		"search[regex]":             {"false"},

		"start":    {"0"},
		"length":   {"20"},
		"OJId":     {OJ},
		"probNum":  {pNum},
		"title":    {pName},
		"source":   {""},
		"category": {"all"},
	}

	return rPOST(apiURL, data)
}

type Inner struct {
	OriginOJ    string `json:"originOJ"`
	OriginProb  string `json:"originProb"`
	AllowSubmit bool   `json:"allowSubmit"`
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	TriggerTime int64  `json:"triggerTime"`
	IsFav       int    `json:"isFav"`
	Status      int    `json:"status"`
}
type searchResult struct {
	Data            []Inner `json:"data"`
	RecordsTotal    int     `json:"recordsTotal"`
	RecordsFiltered int     `json:"recordsFiltered"`
	Draw            int     `json:"draw"`
}

func getPList(body []byte) searchResult {
	var res searchResult
	json.Unmarshal(body, &res)

	return res
}

func getTimeMemory(body string) (string, string) {
	//getting time limit
	need := "ms"
	index1 := strings.Index(body, need)
	runes := []rune(body)

	time := "-"
	if index1 != -1 {
		time = string(runes[0 : index1+2])
	}

	//getting memory limit
	need = "B"
	index2 := strings.Index(body, need)

	memory := "-"
	if index2 != -1 {
		memory = string(runes[index1+2 : index2+1])
	}

	return time, memory
}
func findPResource(OJ, pNum string) {
	pURL := "https://vjudge.net/problem/" + OJ + "-" + pNum

	req, err := http.NewRequest("GET", pURL, nil)
	req.Header.Add("Content-Type", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	response, err := client.Do(req)
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	checkErr(err)

	pTitle = document.Find("div[id='prob-title']").Find("h2").Text()
	pDesSrc, _ = document.Find("iframe").Attr("src")

	time_Memory := document.Find("dl[class='card row']").Find("dd[class='col-sm-7']").Text()
	pTimeLimit, pMemoryLimit = getTimeMemory(time_Memory)
}

func removeStyle(styleBody string) string {
	need1 := "<style"
	index1 := strings.Index(styleBody, need1)
	need2 := "</style>"
	index2 := strings.Index(styleBody, need2)
	runes := []rune(styleBody)

	var part1, part2 string
	if index1 != -1 {
		part1 = string(runes[0:index1])
		part2 = string(runes[index2+8:])

		styleBody = part1 + part2
	}

	return styleBody
}

type LanguagePack struct {
	LangValue string
	LangName  string
}

func getLanguage(OJ string) []LanguagePack {
	var languagePack []LanguagePack

	if OJ == "CodeForces" {
		languagePack = []LanguagePack{
			LanguagePack{LangValue: "14", LangName: "ActiveTcl 8.5"},
			LanguagePack{LangValue: "33", LangName: "Ada GNAT 4"},
			LanguagePack{LangValue: "18", LangName: "Befunge"},
			LanguagePack{LangValue: "52", LangName: "Clang++17 Diagnostics"},
			LanguagePack{LangValue: "9", LangName: "C# Mono 5.18"},
			LanguagePack{LangValue: "28", LangName: "D DMD32 v2.091.0"},
			LanguagePack{LangValue: "3", LangName: "Delphi 7"},
			LanguagePack{LangValue: "25", LangName: "Factor"},
			LanguagePack{LangValue: "39", LangName: "FALSE"},
			LanguagePack{LangValue: "4", LangName: "Free Pascal 3.0.2"},
			LanguagePack{LangValue: "43", LangName: "GNU GCC C11 5.1.0"},
			LanguagePack{LangValue: "45", LangName: "GNU C++11 5 ZIP"},
			LanguagePack{LangValue: "42", LangName: "GNU G++11 5.1.0"},
			LanguagePack{LangValue: "50", LangName: "GNU G++14 6.4.0"},
			LanguagePack{LangValue: "54", LangName: "GNU G++17 7.3.0"},
			LanguagePack{LangValue: "61", LangName: "GNU G++17 9.2.0 (64 bit, msys 2)"},
			LanguagePack{LangValue: "32", LangName: "Go 1.14"},
			LanguagePack{LangValue: "12", LangName: "Haskell GHC 8.6.3"},
			LanguagePack{LangValue: "15", LangName: "Io-2008-01-07 (Win32)"},
			LanguagePack{LangValue: "47", LangName: "J"},
			LanguagePack{LangValue: "36", LangName: "Java 1.8.0_162"},
			LanguagePack{LangValue: "60", LangName: "Java 11.0.5"},
			LanguagePack{LangValue: "46", LangName: "Java 8 ZIP"},
			LanguagePack{LangValue: "34", LangName: "JavaScript V8 4.8.0"},
			LanguagePack{LangValue: "48", LangName: "Kotlin 1.3.70"},
			LanguagePack{LangValue: "56", LangName: "Microsoft Q#"},
			LanguagePack{LangValue: "2", LangName: "Microsoft Visual C++ 2010"},
			LanguagePack{LangValue: "59", LangName: "Microsoft Visual C++ 2017"},
			LanguagePack{LangValue: "38", LangName: "Mysterious Language"},
			LanguagePack{LangValue: "55", LangName: "Node.js 9.4.0"},
			LanguagePack{LangValue: "19", LangName: "OCaml 4.02.1"},
			LanguagePack{LangValue: "22", LangName: "OpenCobol 1.0"},
			LanguagePack{LangValue: "51", LangName: "PascalABC.NET 3.4.2"},
			LanguagePack{LangValue: "13", LangName: "Perl 5.20.1"},
			LanguagePack{LangValue: "6", LangName: "PHP 7.2.13"},
			LanguagePack{LangValue: "44", LangName: "Picat 0.9"},
			LanguagePack{LangValue: "17", LangName: "Pike 7.8"},
			LanguagePack{LangValue: "40", LangName: "PyPy 2.7 (7.2.0)"},
			LanguagePack{LangValue: "41", LangName: "PyPy 3.6 (7.2.0)"},
			LanguagePack{LangValue: "7", LangName: "Python 2.7.15"},
			LanguagePack{LangValue: "31", LangName: "Python 3.7.2"},
			LanguagePack{LangValue: "27", LangName: "Roco"},
			LanguagePack{LangValue: "8", LangName: "Ruby 2.0.0p645"},
			LanguagePack{LangValue: "49", LangName: "Rust 1.42.0"},
			LanguagePack{LangValue: "20", LangName: "Scala 2.12.8"},
			LanguagePack{LangValue: "26", LangName: "Secret_171"},
			LanguagePack{LangValue: "57", LangName: "Text"},
		}
	} else if OJ == "HackerRank" {
		languagePack = []LanguagePack{
			LanguagePack{LangValue: "ada", LangName: "ada"},
			LanguagePack{LangValue: "bash", LangName: "bash"},
			LanguagePack{LangValue: "c", LangName: "c"},
			LanguagePack{LangValue: "clojure", LangName: "clojure"},
			LanguagePack{LangValue: "coffeescript", LangName: "coffeescript"},
			LanguagePack{LangValue: "cpp", LangName: "cpp"},
			LanguagePack{LangValue: "cpp14", LangName: "cpp14"},
			LanguagePack{LangValue: "csharp", LangName: "csharp"},
			LanguagePack{LangValue: "d", LangName: "d"},
			LanguagePack{LangValue: "elixir", LangName: "elixir"},
			LanguagePack{LangValue: "erlang", LangName: "erlang"},
			LanguagePack{LangValue: "fortran", LangName: "fortran"},
			LanguagePack{LangValue: "fsharp", LangName: "fsharp"},
			LanguagePack{LangValue: "go", LangName: "go"},
			LanguagePack{LangValue: "groovy", LangName: "groovy"},
			LanguagePack{LangValue: "haskell", LangName: "haskell"},
			LanguagePack{LangValue: "java", LangName: "java"},
			LanguagePack{LangValue: "java8", LangName: "java8"},
			LanguagePack{LangValue: "javascript", LangName: "javascript"},
			LanguagePack{LangValue: "julia", LangName: "julia"},
			LanguagePack{LangValue: "lolcode", LangName: "lolcode"},
			LanguagePack{LangValue: "lua", LangName: "lua"},
			LanguagePack{LangValue: "objectivec", LangName: "objectivec"},
			LanguagePack{LangValue: "ocaml", LangName: "ocaml"},
			LanguagePack{LangValue: "octave", LangName: "octave"},
			LanguagePack{LangValue: "pascal", LangName: "pascal"},
			LanguagePack{LangValue: "perl", LangName: "perl"},
			LanguagePack{LangValue: "php", LangName: "php"},
			LanguagePack{LangValue: "pypy", LangName: "pypy"},
			LanguagePack{LangValue: "pypy3", LangName: "pypy3"},
			LanguagePack{LangValue: "python", LangName: "python"},
			LanguagePack{LangValue: "python3", LangName: "python3"},
			LanguagePack{LangValue: "r", LangName: "r"},
			LanguagePack{LangValue: "racket", LangName: "racket"},
			LanguagePack{LangValue: "ruby", LangName: "ruby"},
			LanguagePack{LangValue: "rust", LangName: "rust"},
			LanguagePack{LangValue: "sbcl", LangName: "sbcl"},
			LanguagePack{LangValue: "scala", LangName: "scala"},
			LanguagePack{LangValue: "smalltalk", LangName: "smalltalk"},
			LanguagePack{LangValue: "swift", LangName: "swift"},
			LanguagePack{LangValue: "tcl", LangName: "tcl"},
			LanguagePack{LangValue: "visualbasic", LangName: "visualbasic"},
			LanguagePack{LangValue: "whitespace", LangName: "whitespace"},
		}
	} else if OJ == "LightOJ" {
		languagePack = []LanguagePack{
			LanguagePack{LangValue: "C", LangName: "C"},
			LanguagePack{LangValue: "C++", LangName: "C++"},
			LanguagePack{LangValue: "JAVA", LangName: "JAVA"},
			LanguagePack{LangValue: "PASCAL", LangName: "PASCAL"},
		}
	} else if OJ == "UVA" {
		languagePack = []LanguagePack{
			LanguagePack{LangValue: "1", LangName: "ANSI C 5.3.0"},
			LanguagePack{LangValue: "3", LangName: "C++ 5.3.0"},
			LanguagePack{LangValue: "5", LangName: "C++11 5.3.0"},
			LanguagePack{LangValue: "2", LangName: "JAVA 1.8.0"},
			LanguagePack{LangValue: "4", LangName: "PASCAL 3.0.0"},
			LanguagePack{LangValue: "6", LangName: "PYTH3 3.5.1"},
		}
	}

	return languagePack
}
