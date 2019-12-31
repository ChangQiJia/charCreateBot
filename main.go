package main 

import (
	"os"
	"log"
	"fmt"
    tb "gopkg.in/tucnak/telebot.v2"	
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

	b.Start()
	
}