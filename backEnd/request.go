package giveservice

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)
var cookieJar, _ = cookiejar.New(nil)

func rGET(urlStr string) ([]byte, error) {
	client := &http.Client{
		Jar: cookieJar,
	}

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

func rPOST(urlStr string, data url.Values)([]byte, error){
	u, _ := url.ParseRequestURI(urlStr)
	urlStr = u.String()

	client := &http.Client{
		Jar: cookieJar,
	}

	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	//req.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}
