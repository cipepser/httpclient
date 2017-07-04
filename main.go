package main

import (
	"log"

	"./bf"
)

func main() {
	c, _ := bf.NewClient(bf.URL, "user", "passwd", nil)

	mkt := "FX_BTC_JPY"
	before := "31777915"
	after := "31775063"

	es, err := c.GetExecutions(mkt, "", before, after)
	if err != nil {
		log.Fatal(err)
	}

}
