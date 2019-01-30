package main

import (
	"log"
	"fmt"
)


func main() {
	var upds []Update

	b := Bot{"772888415:AAHrwngGee8dskfmAMSlckm-zGeJyxR4LpY", "Markdown"}
	last_upd := 0
	
	for {
		if err := b.GetUpdates(&upds, last_upd, 0, 0, nil); err != nil {
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

			if upds[i].Mes.Text == "Hello" {
				id := int(upds[i].Mes.Chat.Id)
				fmt.Printf("%s> %s\n", "Bot", "Hi!")
				if err := b.SendMessage(id, "Hi!"); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

