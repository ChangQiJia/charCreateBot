package main 

import (
	"os"
	"log"
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"	
	"math/rand"
	"time"
	"strconv"
)

func main(){
	
	fmt.Println ("~~ Starting App")

	var (
        port      = os.Getenv("PORT")       // sets automatically
        publicURL = os.Getenv("PUBLIC_URL") // you must add it to your config vars
        token     = os.Getenv("TOKEN")      // you must add it to your config vars
	)
	
	webhook := &tb.Webhook{
        Listen:   ":" + port,
        Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
	}
	
	pref := tb.Settings{
        Token:  token,
        Poller: webhook,
	}
	
	fmt.Println ("~~ Creating bot")

	b, err := tb.NewBot(pref)
	
    if err != nil {
		fmt.Println ("~~ Oh No!")
        log.Fatal(err)
	}

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hi!")
	})

	b.Handle("/order", func(m *tb.Message) {
		orderedRoll(b, m)
	})

	b.Start()
	
}

func orderedRoll(b *tb.Bot, m *tb.Message){
	rand.Seed(time.Now().UTC().UnixNano())

	for i:= 0; i<6 ;i++{
		eachScore := 0
		min := 10
		outputStr := "("

		for roll:=0 ; roll < 4; roll++{
			oneDsix := rand.Intn(6)+1
			if (oneDsix < min){
				min = oneDsix
			}
			eachScore += oneDsix

			outputStr += strconv.Itoa(oneDsix)

			if (roll < 3){
				outputStr += " + "
			}else{
				outputStr += ") = "
				eachScore -= min
				outputStr += strconv.Itoa(eachScore)
			}
		}

		b.Send(m.Sender, outputStr)
	}
}