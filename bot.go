package main

import (
  "log"
  "strconv"
  "fmt"
  "net/url"
)


func main() {
	var upds []Update

	b := Bot{"772888415:AAHrwngGee8dskfmAMSlckm-zGeJyxR4LpY"}
	last_upd := 0
	
	for {
		v := url.Values{}

		v.Add("offset", strconv.Itoa(last_upd))

		if err, ok := b.Getcmd("getUpdates", v, &upds); err != nil || !ok {
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

