package main 

import (
	"os"
	"log"
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"	
	"math/rand"
	"time"
	"strconv"
	"database/sql"
	"strings"
	_ "github.com/lib/pq"
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

	b.Handle("/crit", func(m *tb.Message) {
		critSuccess(b, m)
	})

	b.Start()
	
}

func help(b *tb.Bot, m *tb.Message){
	outputStr := ""
	outputStr += "For ordered rolls please enter /ordered \nFor unordered rolls please enter /unordered"
	b.Send(m.Chat, outputStr)
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

	b.Send(m.Chat, outputStr)
}

func unorderedRoll(b *tb.Bot, m *tb.Message){
	rand.Seed(time.Now().UTC().UnixNano())
	
	var reroll = true; 
	var valid = false; 

	for reroll {
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
			
			if (eachScore >= 15){
				valid = true
			}

			totalScore += eachScore
			outputStr += "\n"
		}

		if (totalScore >= 70 && totalScore <= 75 && valid){

			reroll=false

			outputStr += "Total Score: "
			outputStr += strconv.Itoa(totalScore)

			b.Send(m.Chat, outputStr)
		}
		fmt.Println (outputStr)
		fmt.Println ("~~ Current total score: ")
		fmt.Println (strconv.Itoa(totalScore))
	}
	
}

func critSuccess(b *tb.Bot, m *tb.Message){

	s := strings.Split(m.Payload, " ")
	fmt.Println (s)
	critType, critValue := s[0], s[1]

	fmt.Println (critType)
	fmt.Println (critValue)

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	
	if err != nil {
		log.Fatal(err)
		fmt.Println (err)
	}

	var (
		effect string
		description string
	)

	rows, err := db.Query("select Effect, Description from Crit where Crit_Type = 'spells' and Min <= 5 and Max >= 5" , critType, critValue, critValue)
	
	if err != nil {
		log.Fatal(err)
		fmt.Println (err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&effect, &description)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(effect, description)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		fmt.Println (err)
	}

	outputStr := "Description: "
	outputStr += description
	outputStr += "\n"
	outputStr += "Effect: "
	outputStr += effect

	b.Send(m.Chat, outputStr)

}