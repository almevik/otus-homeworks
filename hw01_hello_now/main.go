package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	// Place your code here
	t := time.Now()

	tntp, err := ntp.Time("localhost")

	if err == nil {
		log.Fatalf(err.Error())
		return
	}

	text := "current time: " + t.Truncate(time.Second).String()
	text += "\nexact time: " + tntp.Round(time.Second).String()

	fmt.Println(text)
}
