package backEnd

import (
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"net/http/cookiejar"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("frontEnd/*/*"))
}

var cookieJar, _ = cookiejar.New(nil)
var client = &http.Client{
	Jar: cookieJar,
}

var store = sessions.NewCookieStore([]byte("mysession"))

var OJSet = map[string]bool{
	"AtCoder":    true,
	"CodeChef":   true,
	"CodeForces": true,
	"Gym":        true,
	"HackerRank": true,
	"LightOJ":    true,
	"TopCoder":   true,
	"URI":        true,
	"UVA":        true,
}

var Info = map[string]interface{}{}

var lastPage, popUpCause, errorType = "/", "", ""
var pTitle, pTimeLimit, pMemoryLimit, pDesSrc, pOrigin = "", "", "", "", ""

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

type LanguagePack struct {
	LangValue string
	LangName  string
}
