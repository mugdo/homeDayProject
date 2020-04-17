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