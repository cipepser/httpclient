package main

import (
	"fmt"
	"log"

	"./bf"
)

func main() {
	c, _ := bf.NewClient(bf.URL, nil)

	mkt := "FX_BTC_JPY"
	before := "31777915"
	after := "31775063"

	es, err := c.GetExecutions(mkt, "", before, after)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range es {
		fmt.Println(e)
	}

}
