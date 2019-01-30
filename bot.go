package main

import (
  "log"
  "fmt"
  "net/http"
  "net/url"
  "encoding/json"
  "strconv"
)


type tgrespond struct {
	Ok		bool			`json:"ok"`
	Desk	string			`json:"description"`
	Res		json.RawMessage	`json:"result"`
}

type tguser struct {
	Id		float64			`json:"id"`
	Is_bot	bool			`json:"is_bot"`
	Firstn	string			`json:"first_name"`
}

type tgchat struct {
	Id		float64			`json:"id"`
	Tp		string			`json:"type"`		
}

type tgmessage struct {
	Id		float64			`json:"message_id"`
	From	tguser			`json:"from"`
	Date	float64 		`json:"date"`
	Chat	tgchat			`json:"chat"`
	Text	string			`json:"text"`
}

type tgupdate struct {
	Id		float64			`json:"update_id"`
	Mes		tgmessage		`json:"message"`
}


const (
	token = "772888415:AAHrwngGee8dskfmAMSlckm-zGeJyxR4LpY"
	tgurl = "https://api.telegram.org/bot" + token
)


func getcmd(cmd string, par url.Values, ret interface {}) (err error, ok bool) {
	var resp tgrespond

	res, err := http.PostForm(tgurl + "/" + cmd, par)
	if err != nil {
		return
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&resp)
	ret = json.Unmarshal(resp.Res, ret)
	ok = resp.Ok

	return
}

func main() {
	var upds []tgupdate

	last_upd := 0
	
	for {
		v := url.Values{}

		v.Add("offset", strconv.Itoa(last_upd))

		if err, ok := getcmd("getUpdates", v, &upds); err != nil || !ok {
			log.Fatal(err)
		}

		for i := range upds {
			if int(upds[i].Id) <= last_upd {
				continue
			}
			last_upd = int(upds[i].Id)
			fmt.Printf("%s> %s\n",
				upds[i].Mes.From.Firstn,
				upds[i].Mes.Text)
		}
	}
}

