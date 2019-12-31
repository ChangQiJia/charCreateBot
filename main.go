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

	b.Handle("/help", func(m *tb.Message) {
		help(b, m)
	})

	b.Handle("/ordered", func(m *tb.Message) {
		orderedRoll(b, m)
	})

	b.Handle("/unordered", func(m *tb.Message) {
		unorderedRoll(b, m)
	})

	b.Start()
	
}

func help(b *tb.Bot, m *tb.Message){
	outputStr := ""
	outputStr += "For ordered rolls please enter /ordered \n For unordered rolls please enter /unordered"
	b.Send(m.Sender, outputStr)
}

func orderedRoll(b *tb.Bot, m *tb.Message){
	rand.Seed(time.Now().UTC().UnixNano())
	outputStr := ""

	for i:= 0; i<6 ;i++{
		eachScore := 0
		min := 10
		outputStr += "("
		
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

		outputStr += "\n"
	}

	b.Send(m.Sender, outputStr)
}

func unorderedRoll(b *tb.Bot, m *tb.Message){
	rand.Seed(time.Now().UTC().UnixNano())
	
	totalScore := 0
	outputStr := ""

	for ok := true; ok; ok = (totalScore < 70 || totalScore > 75) {
	
		outputStr := ""
		totalScore := 0
		
		for i:= 0; i<6 ;i++{
			eachScore := 0
			min := 10
			outputStr += "("
			
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
			
			totalScore += eachScore
			outputStr += "\n"
		}
	}
	outputStr += "Total Score: "
	outputStr += strconv.Itoa(totalScore)

	b.Send(m.Sender, outputStr)
}