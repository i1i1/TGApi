package main

import (
	"fmt"
	"net/url"
	"net/http"
	"encoding/json"
	"strconv"
)


type Bot struct {
	tok		string
	parm	string
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

type Sticker struct {
	File_id	string			`json:"file_id"`
	Width	float64			`json:"width"`
	Height	float64			`json:"height"`
}

type Message struct {
	Id		float64			`json:"message_id"`
	From	User			`json:"from"`
	Date	float64 		`json:"date"`
	Chat	Chat			`json:"chat"`
	Text	string			`json:"text"`
	Sticker	Sticker			`json:"sticker"`
	Newmem	[]User			`json:"new_chat_members"`
}

type Update struct {
	Id		float64			`json:"update_id"`
	Mes		Message			`json:"message"`
}

type Error struct {
	s		string
}


const apifmt = "https://api.telegram.org/bot%s/%s"


func (e Error) Error() string {
	return fmt.Sprintf("Telegram error: %s\n", e.s)
}

func (b *Bot) request(cmd string, par url.Values, ret interface {}) error {
	var resp Respond

	url := fmt.Sprintf(apifmt, b.tok, cmd)
	res, err := http.PostForm(url, par)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return err
	}
	if !resp.Ok {
		return Error{resp.Desk}
	}

	json.Unmarshal(resp.Res, ret)
	return nil
}

func (b *Bot) GetUpdates(offset, limit, timeout int) (upds []Update, e error) {
	pars := url.Values{}
	if offset > 0 {
		pars.Add("offset", strconv.Itoa(offset))
	}
	if limit > 0 {
		pars.Add("limit", strconv.Itoa(limit))
	}
	if timeout > 0 {
		pars.Add("timeout", strconv.Itoa(timeout))
	}
	e = b.request("getUpdates", pars, &upds)
	return upds, e
}

/*
 * TODO: Some more arguments to add
 */
func (b *Bot) SendMessage(chat_id, mes int, text string) error {
	pars := url.Values{}

	pars.Add("chat_id", strconv.Itoa(chat_id))
	if mes > 0 {
		pars.Add("reply_to_message_id", strconv.Itoa(mes))
	}
	pars.Add("parse_mode", b.parm)
	pars.Add("text", text)

	return b.request("sendMessage", pars, nil)
}

func (b *Bot) KickChatMember(chat, user, until int) error {
	pars := url.Values{}

	pars.Add("chat_id", strconv.Itoa(chat))
	pars.Add("user_id", strconv.Itoa(user))
	pars.Add("until_date", strconv.Itoa(until))

	return b.request("kickChatMember", pars, nil)
}

func (b *Bot) DeleteMessage(chat, mes int) error {
	pars := url.Values{}

	pars.Add("chat_id", strconv.Itoa(chat))
	pars.Add("message_id", strconv.Itoa(mes))

	return b.request("deleteMessage", pars, nil)
}

func (b *Bot) EditMessageText(chat, mes int, text string) error {
	pars := url.Values{}

	pars.Add("chat_id", strconv.Itoa(chat))
	pars.Add("message_id", strconv.Itoa(mes))
	pars.Add("text", text)

	return b.request("editMessageText", pars, nil)
}

func (b *Bot) GetMe() (u User, e error) {
	e = b.request("getMe", url.Values{}, &u)
	return
}

