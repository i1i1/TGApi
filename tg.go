package main

import (
  "fmt"
  "net/url"
  "net/http"
  "encoding/json"
)


type Bot struct {
	tok		string
}

type Respond struct {
	Ok		bool			`json:"ok"`
	Desk	string			`json:"description"`
	Res		json.RawMessage	`json:"result"`
}

type User struct {
	Id		float64			`json:"id"`
	Is_bot	bool			`json:"is_bot"`
	Firstn	string			`json:"first_name"`
}

type Chat struct {
	Id		float64			`json:"id"`
	Tp		string			`json:"type"`		
}

type Message struct {
	Id		float64			`json:"message_id"`
	From	User			`json:"from"`
	Date	float64 		`json:"date"`
	Chat	Chat			`json:"chat"`
	Text	string			`json:"text"`
}

type Update struct {
	Id		float64			`json:"update_id"`
	Mes		Message			`json:"message"`
}

const apifmt = "https://api.telegram.org/bot%s/%s"


func (b Bot) Getcmd(cmd string, par url.Values, ret interface {}) (err error, ok bool) {
	var resp Respond

	url := fmt.Sprintf(apifmt, b.tok, cmd)
	res, err := http.PostForm(url, par)
	if err != nil {
		return
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&resp)
	ret = json.Unmarshal(resp.Res, ret)
	ok = resp.Ok

	return
}


