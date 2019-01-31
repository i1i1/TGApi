package main

import (
	"log"
	"os/exec"
	"fmt"
	"strings"
)


func execute(cmd, in string) string {
	args := strings.Fields(cmd)
	c := exec.Command(args[0], args[1:]...)

	stdin, err := c.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(stdin, "%s", in)
	stdin.Close()

	out, err := c.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	return string(out)
}

func (b * Bot) updates(last_upd int) []Update {
	var upds []Update

	if err := b.GetUpdates(&upds, last_upd, 0, 0, nil); err != nil {
		log.Fatal(err)
	}
	return upds
}

func (b *Bot) send(id int, s string) {
	if s == "" {
		if err := b.SendMessage(id, ""); err != nil {
			log.Fatal(err)
		}
		return
	}
	for len(s) > 4096 {
		if err := b.SendMessage(id, string(s[:4095])); err != nil {
			log.Fatal(err)
		}
		s = s[4096:]
	}

	if err := b.SendMessage(id, string(s)); err != nil {
		log.Fatal(err)
	}
}

func main() {
	b := Bot{"772888415:AAHrwngGee8dskfmAMSlckm-zGeJyxR4LpY", "Markdown"}
	last_upd := 0
	
	for {
		upds := b.updates(last_upd)

		for i := range upds {
			if int(upds[i].Id) <= last_upd {
				continue
			}
			last_upd = int(upds[i].Id)
			id := int(upds[i].Mes.Chat.Id)
			mes := upds[i].Mes.Text

			if mes[0] == '/' {
				continue
			}

			b.send(id, execute("bc -l", mes+"\n"))
		}
	}
}

