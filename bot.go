package main

import (
	"fmt"
	"time"
)

const (
	mutetm=		5*60
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

func main() {
	b := Bot{"oops here should be your bot token", "Markdown"}
	last_upd := 0
	stickdb := make(map[int]stickent)
	
	for {
		upds := b.updates(last_upd)

		for i := range upds {
			if int(upds[i].Id) <= last_upd {
				continue
			}

			last_upd = int(upds[i].Id)
			chat := int(upds[i].Mes.Chat.Id)
			stick := upds[i].Mes.Sticker
			sender := int(upds[i].Mes.From.Id)
			mes := int(upds[i].Mes.Id)
			firstn := upds[i].Mes.From.Firstn

			/* If muted */
			if v, ok := stickdb[sender]; ok && v.n == 0 {
				b.DeleteMessage(chat, mes)
				continue
			}
			/* If not a sticker */
			if stick.File_id == "" {
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

