package main

import (
	"fmt"
	"time"
)

const (
	mutetm=		5*60
	timeout=	60
	stickers=	3
)


type stickent struct {
	n int
	t time.Time
}


func (b * Bot) updates(last_upd int) []Update {
	var upds []Update

	if err := b.GetUpdates(&upds, last_upd, 0, 0, nil); err != nil {
		fmt.Print(err)
	}
	return upds
}

func (b *Bot) send(id, reply int, s string) {
	if s == "" {
		if err := b.SendMessage(id, reply, "-/-/-/-"); err != nil {
			fmt.Print(err)
		}
		return
	}
	for len(s) > 4096 {
		if err := b.SendMessage(id, reply, string(s[:4095])); err != nil {
			fmt.Print(err)
		}
		s = s[4096:]
	}

	if err := b.SendMessage(id, reply, string(s)); err != nil {
		fmt.Print(err)
	}
}

func isstick(m Message) bool {
	return m.Sticker.File_id != ""
}

func sendedbefore(m Message) bool {
	tm := time.Unix(int64(m.Date), 0)
	/*
	 * Lets assume that all messages recived before timeout should be ignored.
	 */
	tm = tm.Add(time.Second*timeout)
	return tm.Before(time.Now())
}

func main() {
	b := Bot{"oops here should be your bot token", "Markdown"}
	last_upd := 0
	stickdb := make(map[int]stickent)
	
	muted := func (sender int) bool {
		v, ok := stickdb[sender]
		return ok && v.n == 0
	}

	for {
		upds := b.updates(last_upd)

		for i := range upds {
			if int(upds[i].Id) <= last_upd {
				continue
			}

			M := upds[i].Mes
			last_upd = int(upds[i].Id)
			chat := int(M.Chat.Id)
			sender := int(M.From.Id)
			mes := int(M.Id)
			firstn := M.From.Firstn

			if muted(sender) {
				b.DeleteMessage(chat, mes)
				continue
			}
			if !isstick(M) {
				continue
			}
			if sendedbefore(M) {
				continue
			}

			if _, ok := stickdb[sender]; !ok {
				var a stickent
				a.n = stickers
				stickdb[sender] = a
			}
			if stickdb[sender].n > 1 {
				var a stickent

				a.n = stickdb[sender].n - 1
				a.t = time.Now().Add(time.Second*mutetm)
				stickdb[sender] = a

				s := fmt.Sprintf("*%s warning!* only %d stickers more allowed\n",
						firstn, stickdb[sender].n)
				b.send(chat, mes, s)
			} else {
				var a stickent

				a.n = 0
				a.t = time.Now().Add(time.Second*mutetm)
				stickdb[sender] = a

				s := fmt.Sprintf("*%s* is now *muted* for *%d* seconds!\n",
					firstn, mutetm)
				b.send(chat, mes, s)
			}
		}

		for k, v := range stickdb {
			if v.t.Before(time.Now()) {
				delete(stickdb, k)
			}
		}
	}
}

