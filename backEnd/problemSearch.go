package backEnd

import (
	"encoding/json"
	"net/url"
)

func pSearch(oj,pNum,pName string) ([]byte, error) {
	apiURL := "https://vjudge.net/problem/data"

	data := url.Values{}
	data.Set("draw", "1")
	data.Set("columns[0][data]", "0")
	data.Set("columns[0][name]", "")
	data.Set("columns[0][searchable]", "true")
	data.Set("columns[0][orderable]", "false")
	data.Set("columns[0][search][value]", "")
	data.Set("columns[0][search][regex]", "false")
	data.Set("columns[1][data]", "1")
	data.Set("columns[1][name]", "")
	data.Set("columns[1][searchable]", "true")
	data.Set("columns[1][orderable]", "true")
	data.Set("columns[1][search][value]", "")
	data.Set("columns[1][search][regex]", "false")
	data.Set("columns[2][data]", "2")
	data.Set("columns[2][name]", "")
	data.Set("columns[2][searchable]", "false")
	data.Set("columns[2][orderable]", "false")
	data.Set("columns[2][search][value]", "")
	data.Set("columns[2][search][regex]", "false")
	data.Set("columns[3][data]", "3")
	data.Set("columns[3][name]", "")
	data.Set("columns[3][searchable]", "true")
	data.Set("columns[3][orderable]", "true")
	data.Set("columns[3][search][value]", "")
	data.Set("columns[3][search][regex]", "false")
	data.Set("columns[4][data]", "4")
	data.Set("columns[4][name]", "")
	data.Set("columns[4][searchable]", "true")
	data.Set("columns[4][orderable]", "true")
	data.Set("columns[4][search][value]", "")
	data.Set("columns[4][search][regex]", "false")
	data.Set("columns[5][data]", "5")
	data.Set("columns[5][name]", "")
	data.Set("columns[5][searchable]", "true")
	data.Set("columns[5][orderable]", "true")
	data.Set("columns[5][search][value]", "")
	data.Set("columns[5][search][regex]", "false")
	data.Set("columns[6][data]", "6")
	data.Set("columns[6][name]", "")
	data.Set("columns[6][searchable]", "true")
	data.Set("columns[6][orderable]", "false")
	data.Set("columns[6][search][value]", "")
	data.Set("columns[6][search][regex]", "false")
	data.Set("order[0][column]", "5")
	data.Set("order[0][dir]", "desc")
	data.Set("start", "0")
	data.Set("length", "20")
	data.Set("search[value]", "")
	data.Set("search[regex]", "false")
	data.Set("OJId", oj)
	data.Set("probNum", pNum)
	data.Set("title", pName)
	data.Set("source", "")
	data.Set("category", "all")

	return rPOST(apiURL, data)
}

type Inner struct {
	OriginOJ		string		`json:"originOJ"`
	OriginProb		string		`json:"originProb"`
	AllowSubmit		bool		`json:"allowSubmit"`
	ID				int64		`json:"id"`
	Title			string		`json:"title"`
	TriggerTime		int64		`json:"triggerTime"`
	IsFav			int			`json:"isFav"`
	Status			int			`json:"status"`
}
type searchResult struct {
	Data				[]Inner	`json:"data"`
	RecordsTotal		int		`json:"recordsTotal"`
	RecordsFiltered		int		`json:"recordsFiltered"`
	Draw				int		`json:"draw"`
}

func getPList(body []byte) (searchResult){
	var res searchResult

	json.Unmarshal(body, &res)

	// fmt.Println("Problem Search Done")
	// fmt.Println(res.RecordsFiltered)

	return res
}
