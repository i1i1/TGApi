package main

import (
	"fmt"
	"time"
)

const (
	mutetm=		60
	stickers=	5
)


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
	stickdb := make(map[int]int)
	mutedb := make(map[int]time.Time)
	
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
			if _, ok := mutedb[sender]; ok {
				b.DeleteMessage(chat, mes)
				continue
			}
			/* If not a sticker */
			if stick.File_id == "" {
				continue
			}

			if _, ok := stickdb[sender]; !ok {
				stickdb[sender] = stickers
			}
			if stickdb[sender] > 0 {
				stickdb[sender]--
				s := fmt.Sprintf("*%s warning!* only %d stickers more allowed\n",
						firstn, stickdb[sender])
				b.send(chat, mes, s)
			} else {
				mutedb[sender] = time.Now().Add(time.Second*mutetm)
				s := fmt.Sprintf("*%s* is now *muted* for *%d* seconds!\n",
					firstn, mutetm)
				b.send(chat, mes, s)
			}
		}

		for k, v := range mutedb {
			if v.Before(time.Now()) {
				stickdb[k] = stickers
				delete(mutedb, k)
			}
		}
	}
}

