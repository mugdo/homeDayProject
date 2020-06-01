package backEnd

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"net/url"
	"strconv"
	"strings"
)

func checkErr(err error) {
	if err != nil {
		//fmt.Println(err)
		panic(err)
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
func pSearch(OJ, pNum, pName, length string) ([]byte, error) {
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
		"length":   {length},
		"OJId":     {OJ},
		"probNum":  {pNum},
		"title":    {pName},
		"source":   {""},
		"category": {"all"},
	}

	return rPOST(apiURL, data)
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
func getLanguage(w http.ResponseWriter, r *http.Request) {
	var languagePack []LanguagePack

	path := r.URL.Path
	runes := []rune(path)
	need := "="
	index := strings.Index(path, need)
	OJ := string(runes[index+1:])

	if OJ == "AtCoder" {
		languagePack = []LanguagePack{
			LanguagePack{LangValue: "3506", LangName: "Awk (mawk 1.3.3)"},
			LanguagePack{LangValue: "3001", LangName: "Bash (GNU bash v4.3.11)"},
			LanguagePack{LangValue: "3507", LangName: "Brainfuck (bf 20041219)"},
			LanguagePack{LangValue: "3004", LangName: "C (Clang 3.8.0)"},
			LanguagePack{LangValue: "3002", LangName: "C (GCC 5.3.0)"},
			LanguagePack{LangValue: "3006", LangName: "C# (Mono 4.2.2.30)"},
			LanguagePack{LangValue: "3005", LangName: "C++14 (Clang 3.8.0)"},
			LanguagePack{LangValue: "3003", LangName: "C++14 (GCC 5.3.0)"},
			LanguagePack{LangValue: "3517", LangName: "Ceylon (1.2.1)"},
			LanguagePack{LangValue: "3007", LangName: "Clojure (1.8.0)"},
			LanguagePack{LangValue: "3008", LangName: "Common Lisp (SBCL 1.1.14)"},
			LanguagePack{LangValue: "3511", LangName: "Crystal (0.12.0)"},
			LanguagePack{LangValue: "3009", LangName: "D (DMD64 v2.070.1)"},
			LanguagePack{LangValue: "3011", LangName: "D (GDC 4.9.3)"},
			LanguagePack{LangValue: "3010", LangName: "D (LDC 0.17.0)"},
			LanguagePack{LangValue: "3512", LangName: "F# (Mono 4.2.2.30)"},
			LanguagePack{LangValue: "3012", LangName: "Fortran (gfortran v4.8.4)"},
			LanguagePack{LangValue: "3013", LangName: "Go (1.6)"},
			LanguagePack{LangValue: "3014", LangName: "Haskell (GHC 7.10)"},
			LanguagePack{LangValue: "3015", LangName: "Java7 (OpenJDK 1.7.0)"},
			LanguagePack{LangValue: "3016", LangName: "Java8 (OpenJDK 1.8.0)"},
			LanguagePack{LangValue: "3017", LangName: "JavaScript (node.js v5.7)"},
			LanguagePack{LangValue: "3518", LangName: "Julia (0.4.2)"},
			LanguagePack{LangValue: "3523", LangName: "Kotlin (1.0.0)"},
			LanguagePack{LangValue: "3514", LangName: "Lua (5.3.2)"},
			LanguagePack{LangValue: "3515", LangName: "LuaJIT (2.0.2)"},
			LanguagePack{LangValue: "3516", LangName: "MoonScript (0.4.0)"},
			LanguagePack{LangValue: "3520", LangName: "Nim (0.13.0)"},
			LanguagePack{LangValue: "3502", LangName: "Objective-C (Clang3.7.1)"},
			LanguagePack{LangValue: "3501", LangName: "Objective-C (GCC 5.3.0)"},
			LanguagePack{LangValue: "3018", LangName: "OCaml (4.02.3)"},
			LanguagePack{LangValue: "3519", LangName: "Octave (4.0.0)"},
			LanguagePack{LangValue: "3019", LangName: "Pascal (FPC 2.6.2)"},
			LanguagePack{LangValue: "3020", LangName: "Perl (v5.18.2)"},
			LanguagePack{LangValue: "3522", LangName: "Perl6 (rakudo-star 2016.01)"},
			LanguagePack{LangValue: "3021", LangName: "PHP (5.6.18)"},
			LanguagePack{LangValue: "3524", LangName: "PHP7 (7.0.4)"},
			LanguagePack{LangValue: "3509", LangName: "PyPy2 (4.0.1)"},
			LanguagePack{LangValue: "3510", LangName: "PyPy3 (2.4.0)"},
			LanguagePack{LangValue: "3022", LangName: "Python2 (2.7.6)"},
			LanguagePack{LangValue: "3023", LangName: "Python3 (3.4.3)"},
			LanguagePack{LangValue: "3024", LangName: "Ruby (2.3.0)"},
			LanguagePack{LangValue: "3504", LangName: "Rust (1.7.0)"},
			LanguagePack{LangValue: "3025", LangName: "Scala (2.11.7)"},
			LanguagePack{LangValue: "3026", LangName: "Scheme (Gauche 0.9.3.3)"},
			LanguagePack{LangValue: "3505", LangName: "Sed (GNU sed 4.2.2)"},
			LanguagePack{LangValue: "3508", LangName: "Standard ML (MLton 20100608)"},
			LanguagePack{LangValue: "3503", LangName: "Swift (swift-2.2-RELEASE)"},
			LanguagePack{LangValue: "3027", LangName: "Text (cat)"},
			LanguagePack{LangValue: "3521", LangName: "TypeScript (1.8.2)"},
			LanguagePack{LangValue: "3513", LangName: "Unlambda (0.1.3)"},
			LanguagePack{LangValue: "3028", LangName: "Visual Basic (Mono 4.2.2.30)"},
		}
	} else if OJ == "CodeChef" {
		languagePack = []LanguagePack{
			LanguagePack{LangValue: "7", LangName: "ADA 95(gnat 6.3)"},
			LanguagePack{LangValue: "13", LangName: "Assembler(nasm 2.12.01)"},
			LanguagePack{LangValue: "28", LangName: "Bash(bash 4.4.5)"},
			LanguagePack{LangValue: "12", LangName: "Brainf**k(bff 1.0.6)"},
			LanguagePack{LangValue: "11", LangName: "C(gcc 6.3)"},
			LanguagePack{LangValue: "27", LangName: "C#(gmcs 4.6.2)"},
			LanguagePack{LangValue: "44", LangName: "C++14(gcc 6.3)"},
			LanguagePack{LangValue: "63", LangName: "C++17(gcc 9.1)"},
			LanguagePack{LangValue: "14", LangName: "Clips(clips 6.24)"},
			LanguagePack{LangValue: "111", LangName: "Clojure(clojure 1.8.0)"},
			LanguagePack{LangValue: "118", LangName: "Cobol(open-cobol 1.1.0)"},
			LanguagePack{LangValue: "32", LangName: "Common Lisp(clisp 2.49)"},
			LanguagePack{LangValue: "31", LangName: "Common Lisp(sbcl 1.3.13)"},
			LanguagePack{LangValue: "20", LangName: "D(gdc 6.3)"},
			LanguagePack{LangValue: "36", LangName: "Erlang(erl 19)"},
			LanguagePack{LangValue: "124", LangName: "F#(mono 4.0.0)"},
			LanguagePack{LangValue: "5", LangName: "Fortran(gfortran 6.3)"},
			LanguagePack{LangValue: "114", LangName: "Go(go 1.7.4)"},
			LanguagePack{LangValue: "21", LangName: "Haskell(ghc 8.0.1)"},
			LanguagePack{LangValue: "16", LangName: "Icon(iconc 9.5.1)"},
			LanguagePack{LangValue: "9", LangName: "Intercal(ick 0.3)"},
			LanguagePack{LangValue: "10", LangName: "Java(HotSpot 8u112)"},
			LanguagePack{LangValue: "56", LangName: "JavaScript(node 7.4.0)"},
			LanguagePack{LangValue: "35", LangName: "JavaScript(rhino 1.7.7)"},
			LanguagePack{LangValue: "47", LangName: "Kotlin(kotlin 1.2.50)"},
			LanguagePack{LangValue: "26", LangName: "Lua(luac 5.3.3)"},
			LanguagePack{LangValue: "30", LangName: "Nemerle(ncc 1.2.0)"},
			LanguagePack{LangValue: "25", LangName: "Nice(nicec 0.9.13)"},
			LanguagePack{LangValue: "8", LangName: "Ocaml(ocamlopt 4.01)"},
			LanguagePack{LangValue: "22", LangName: "Pascal(fpc 3.0.0)"},
			LanguagePack{LangValue: "2", LangName: "Pascal(gpc 20070904)"},
			LanguagePack{LangValue: "3", LangName: "Perl(perl 5.24.1)"},
			LanguagePack{LangValue: "54", LangName: "Perl6(perl 6)"},
			LanguagePack{LangValue: "29", LangName: "PHP(php 7.1.0)"},
			LanguagePack{LangValue: "19", LangName: "Pike(pike 8.0)"},
			LanguagePack{LangValue: "15", LangName: "Prolog(swi 7.2.3)"},
			LanguagePack{LangValue: "109", LangName: "PyPy 3(PyPy&nbsp;3.5)"},
			LanguagePack{LangValue: "99", LangName: "PyPy(PyPy 2.6.0)"},
			LanguagePack{LangValue: "4", LangName: "Python(cpython 2.7.13)"},
			LanguagePack{LangValue: "116", LangName: "Python3(python 3.6)"},
			LanguagePack{LangValue: "117", LangName: "R(3.3.2)"},
			LanguagePack{LangValue: "17", LangName: "Ruby(ruby 2.3.3)"},
			LanguagePack{LangValue: "93", LangName: "Rust(rust 1.14.0)"},
			LanguagePack{LangValue: "39", LangName: "Scala(scala 2.12.1)"},
			LanguagePack{LangValue: "97", LangName: "Scheme(chicken 4.11.0)"},
			LanguagePack{LangValue: "33", LangName: "Scheme(guile 2.0.13)"},
			LanguagePack{LangValue: "18", LangName: "Scheme(stalin 0.3)"},
			LanguagePack{LangValue: "23", LangName: "Smalltalk(gst 3.2.5)"},
			LanguagePack{LangValue: "40", LangName: "SQL(sqlite 3.27.2)"},
			LanguagePack{LangValue: "85", LangName: "Swift(swift 3.0.2)"},
			LanguagePack{LangValue: "38", LangName: "Tcl(tcl 8.6)"},
			LanguagePack{LangValue: "62", LangName: "Text(pure text)"},
			LanguagePack{LangValue: "6", LangName: "Whitespace(wspace 0.3)"},
		}
	} else if OJ == "CodeForces" {
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
	} else if OJ == "Gym" {
		languagePack = []LanguagePack{
			LanguagePack{LangValue: "9", LangName: "C# Mono 5.18"},
			LanguagePack{LangValue: "52", LangName: "Clang++17 Diagnostics"},
			LanguagePack{LangValue: "28", LangName: "D DMD32 v2.091.0"},
			LanguagePack{LangValue: "3", LangName: "Delphi 7"},
			LanguagePack{LangValue: "4", LangName: "Free Pascal 3.0.2"},
			LanguagePack{LangValue: "42", LangName: "GNU G++11 5.1.0"},
			LanguagePack{LangValue: "50", LangName: "GNU G++14 6.4.0"},
			LanguagePack{LangValue: "54", LangName: "GNU G++17 7.3.0"},
			LanguagePack{LangValue: "61", LangName: "GNU G++17 9.2.0 (64 bit, msys 2)"},
			LanguagePack{LangValue: "43", LangName: "GNU GCC C11 5.1.0"},
			LanguagePack{LangValue: "32", LangName: "Go 1.14"},
			LanguagePack{LangValue: "12", LangName: "Haskell GHC 8.6.3"},
			LanguagePack{LangValue: "36", LangName: "Java 1.8.0_162"},
			LanguagePack{LangValue: "60", LangName: "Java 11.0.5"},
			LanguagePack{LangValue: "34", LangName: "JavaScript V8 4.8.0"},
			LanguagePack{LangValue: "48", LangName: "Kotlin 1.3.70"},
			LanguagePack{LangValue: "2", LangName: "Microsoft Visual C++ 2010"},
			LanguagePack{LangValue: "59", LangName: "Microsoft Visual C++ 2017"},
			LanguagePack{LangValue: "55", LangName: "Node.js 9.4.0"},
			LanguagePack{LangValue: "19", LangName: "OCaml 4.02.1"},
			LanguagePack{LangValue: "51", LangName: "PascalABC.NET 3.4.2"},
			LanguagePack{LangValue: "13", LangName: "Perl 5.20.1"},
			LanguagePack{LangValue: "6", LangName: "PHP 7.2.13"},
			LanguagePack{LangValue: "40", LangName: "PyPy 2.7 (7.2.0)"},
			LanguagePack{LangValue: "41", LangName: "PyPy 3.6 (7.2.0)"},
			LanguagePack{LangValue: "7", LangName: "Python 2.7.15"},
			LanguagePack{LangValue: "31", LangName: "Python 3.7.2"},
			LanguagePack{LangValue: "8", LangName: "Ruby 2.0.0p645"},
			LanguagePack{LangValue: "49", LangName: "Rust 1.42.0"},
			LanguagePack{LangValue: "20", LangName: "Scala 2.12.8"},
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
	} else if OJ == "TopCoder" {
		languagePack = []LanguagePack{
			LanguagePack{LangValue: "4", LangName: "C#"},
			LanguagePack{LangValue: "3", LangName: "C++"},
			LanguagePack{LangValue: "1", LangName: "Java"},
			LanguagePack{LangValue: "6", LangName: "Python"},
			LanguagePack{LangValue: "5", LangName: "VB"},
		}
	} else if OJ == "URI" {
		languagePack = []LanguagePack{
			LanguagePack{LangValue: "1", LangName: "C (gcc  4.8.5)"},
			LanguagePack{LangValue: "14", LangName: "C99 (gcc  4.8.5)"},
			LanguagePack{LangValue: "7", LangName: "C# (mono 5.10.1.20)"},
			LanguagePack{LangValue: "2", LangName: "C++ (g++ 4.8.5)"},
			LanguagePack{LangValue: "16", LangName: "C++17 (g++ 7.3.0)"},
			LanguagePack{LangValue: "12", LangName: "Go (go 1.8.1)"},
			LanguagePack{LangValue: "17", LangName: "Haskell (ghc 7.6.3)"},
			LanguagePack{LangValue: "3", LangName: "Java 7 (OpenJDK 1.7.0)"},
			LanguagePack{LangValue: "11", LangName: "Java 8 (OpenJDK 1.8.0)"},
			LanguagePack{LangValue: "10", LangName: "JavaScript (nodejs 8.4.0)"},
			LanguagePack{LangValue: "15", LangName: "Kotlin (1.2.10)"},
			LanguagePack{LangValue: "9", LangName: "Lua (lua 5.2.3)"},
			LanguagePack{LangValue: "18", LangName: "OCaml (ocamlc 4.01.0)"},
			LanguagePack{LangValue: "19", LangName: "Pascal (fpc 2.6.2)"},
			LanguagePack{LangValue: "13", LangName: "PostgreSQL (psql 9.4.19)"},
			LanguagePack{LangValue: "4", LangName: "Python 2 (Python 2.7.6)"},
			LanguagePack{LangValue: "5", LangName: "Python 3 (Python 3.4.3)"},
			LanguagePack{LangValue: "6", LangName: "Ruby (ruby 2.3.0)"},
			LanguagePack{LangValue: "8", LangName: "Scala (scalac 2.11.8)"},
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

	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(languagePack)
	w.Write(b)
}
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func sendMail(email, link, resetType string) {
	// Choose auth method and set it up
	auth := smtp.PlainAuth("", "ajudge.bd", "aj199273", "smtp.gmail.com")

	// Here we do it all: connect to our server, set up a message and send it
	to := []string{email}

	var msg []byte
	if resetType == "password" {
		msg = []byte("From: ajudge Team\r\n" +
			"To: " + email + "\r\n" +
			"Subject: Ajudge Password Reset Link\r\n" +
			"\r\n" +
			"Here’s the link of password reset. Click on the link below:\r\n" +
			link)
	} else if resetType == "email" {
		msg = []byte("From: ajudge Team\r\n" +
			"To: " + email + "\r\n" +
			"Subject: Ajudge Email verification Link\r\n" +
			"\r\n" +
			"Here’s the link of account activation. Click on the link below:\r\n" +
			link)
	}

	err := smtp.SendMail("smtp.gmail.com:587", auth, "", to, msg)
	checkErr(err)
}
func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)

	hasher := md5.New()
	hasher.Write(b)
	return hex.EncodeToString(hasher.Sum(nil))
}
func getOriginLink(apiURL string) string {
	req, _ := http.NewRequest("GET", apiURL, nil)

	//setting up requset to prevent auto redirect
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	header, _ := resp.Location()

	return header.String()
}
