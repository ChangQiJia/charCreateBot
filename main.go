package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {

	fmt.Println("~~ Starting App")

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

	fmt.Println("~~ Creating bot")

	b, err := tb.NewBot(pref)

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}

	if err != nil {
		fmt.Println("~~ Oh No!")
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
		critSuccess(b, m, db)
	})

	b.Handle("/fail", func(m *tb.Message) {
		critFail(b, m, db)
	})

	b.Handle("/health", func(m *tb.Message) {
		rollHealth(b, m)
	})

	b.Handle("/classname", func(m *tb.Message) {
		getClassNames(b, m)
	})

	b.Handle("/roll", func(m *tb.Message) {
		roll(b, m)
	})

	b.Start()

}

func help(b *tb.Bot, m *tb.Message) {
	outputStr := ""
	outputStr += "For ordered rolls please enter /ordered \nFor unordered rolls please enter /unordered"
	outputStr += "\nFor crit success please enter /crit followed by Spell or Attack follow by a space and a number example /crit Spell 20\n"
	outputStr += "\nFor crit fails please enter /fail followed by Spell or Attack follow by a space and a number example /fail Attack 90\n"
	outputStr += "\nFor rolling health please enter /health follow by con mod and how you level your characther , example /health 2 rogue:3 pala:5, please note you can only enter 9 different classes max\n"
	outputStr += "\nFor class abbreviations please enter /classname\n"
	outputStr += "\nFor rolls please enter /roll followed by dice:number of rolls. example /roll 6:2 10:2 will roll d6 2 times followed by a d10 2 times"
	b.Send(m.Chat, outputStr)
}

func orderedRoll(b *tb.Bot, m *tb.Message) {
	rand.Seed(time.Now().UTC().UnixNano())
	outputStr := ""
	largestAmount := 0
	secondAmount := 0
	largestIndex := -1
	secondIndex := -1

	for i := 0; i < 6; i++ {
		eachScore := 0
		min := 10
		outputStr += "("

		for roll := 0; roll < 4; roll++ {
			oneDsix := rand.Intn(6) + 1
			if oneDsix < min {
				min = oneDsix
			}
			eachScore += oneDsix

			outputStr += strconv.Itoa(oneDsix)

			if roll < 3 {
				outputStr += " + "
			} else {
				outputStr += ") = "
				eachScore -= min
				outputStr += strconv.Itoa(eachScore)
			}
		}

		if eachScore > largestAmount {
			secondAmount = largestAmount
			secondIndex = largestIndex

			largestAmount = eachScore
			largestIndex = i
		} else if eachScore > secondAmount {
			secondAmount = eachScore
			secondIndex = i
		}

		outputStr += "\n"
	}

	suggestion := getSuggestion(largestIndex, secondIndex)

	outputStr += "\n"
	outputStr += suggestion

	b.Send(m.Chat, outputStr)
}

func getSuggestion(first int, second int) string {
	output := "Suggestion: "

	if first == 0 {
		if second == 1 {
			output += "Barbarian"
		} else if second == 2 {
			output += "Fighter"
		} else if second == 3 {
			output += "Eldritch Knight"
		} else if second == 4 {
			output += "Ranger"
		} else if second == 5 {
			output += "Bard-barian"
		}
	} else if first == 1 {
		if second == 0 {
			output += "Rogue-barian"
		} else if second == 2 {
			output += "Monk"
		} else if second == 3 {
			output += "Arcane Trickster"
		} else if second == 4 {
			output += "DruidMonk"
		} else if second == 5 {
			output += "Roga-din"
		}
	} else if first == 2 {
		if second == 0 {
			output += "Barbarian"
		} else if second == 1 {
			output += "Ranged Fighter"
		} else if second == 3 {
			output += "Wizard"
		} else if second == 4 {
			output += "Cleric"
		} else if second == 5 {
			output += "Sorc"
		}
	} else if first == 3 {
		if second == 0 {
			output += "Artificer"
		} else if second == 1 {
			output += "BladeSinger"
		} else if second == 2 {
			output += "Mystic"
		} else if second == 4 {
			output += "WizardCleric"
		} else if second == 5 {
			output += "Wiz-Sorc"
		}
	} else if first == 4 {
		if second == 0 {
			output += "War Cleric"
		} else if second == 1 {
			output += "Ranger"
		} else if second == 2 {
			output += "Druid"
		} else if second == 3 {
			output += "Druid"
		} else if second == 5 {
			output += "Cleric Warlock (Child of divorced patrons)"
		}
	} else if first == 5 {
		if second == 0 {
			output += "Hexadin"
		} else if second == 1 {
			output += "Bard"
		} else if second == 2 {
			output += "Sorc-lock"
		} else if second == 3 {
			output += "Wiz-lock"
		} else if second == 4 {
			output += "Bard-ric"
		}
	}

	output += ", Ranger"

	return output

}

func unorderedRoll(b *tb.Bot, m *tb.Message) {
	rand.Seed(time.Now().UTC().UnixNano())

	var reroll = true
	var valid = false

	for reroll {
		outputStr := ""
		totalScore := 0

		for i := 0; i < 6; i++ {
			eachScore := 0
			min := 10
			outputStr += "("

			for roll := 0; roll < 4; roll++ {
				oneDsix := rand.Intn(6) + 1
				if oneDsix < min {
					min = oneDsix
				}
				eachScore += oneDsix

				outputStr += strconv.Itoa(oneDsix)

				if roll < 3 {
					outputStr += " + "
				} else {
					outputStr += ") = "
					eachScore -= min
					outputStr += strconv.Itoa(eachScore)
				}
			}

			if eachScore >= 15 {
				valid = true
			}

			totalScore += eachScore
			outputStr += "\n"
		}

		if totalScore >= 70 && totalScore <= 75 && valid {

			reroll = false

			outputStr += "Total Score: "
			outputStr += strconv.Itoa(totalScore)

			b.Send(m.Chat, outputStr)
		}
		fmt.Println(outputStr)
		fmt.Println("~~ Current total score: ")
		fmt.Println(strconv.Itoa(totalScore))
	}

}

func critSuccess(b *tb.Bot, m *tb.Message, db *sql.DB) {

	s := strings.Split(m.Payload, " ")
	fmt.Println(s)
	critType, critValue := s[0], s[1]
	critType = strings.Title(strings.ToLower(critType))

	fmt.Println(critType)
	fmt.Println(critValue)
	i1, _ := strconv.Atoi(critValue)

	var (
		effect      string
		description string
	)

	rows, err := db.Query("select \"Effect\", \"Description\" from public.\"Crit\" where \"Crit_Type\" = $1 and \"Min\" <= $2 and \"Max\" >= $2", critType, i1)

	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
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
		fmt.Println(err)
	}

	outputStr := "Description: "
	outputStr += description
	outputStr += "\n\n"
	outputStr += "Effect: "
	outputStr += effect

	b.Send(m.Chat, outputStr)

}

func critFail(b *tb.Bot, m *tb.Message, db *sql.DB) {

	s := strings.Split(m.Payload, " ")
	fmt.Println(s)
	critType, critValue := s[0], s[1]
	critType = strings.Title(strings.ToLower(critType))

	fmt.Println(critType)
	fmt.Println(critValue)
	i1, _ := strconv.Atoi(critValue)

	var (
		effect      string
		description string
	)

	rows, err := db.Query("select \"Effect\", \"Description\" from public.\"CritFail\" where \"Crit_Type\" = $1 and \"Min\" <= $2 and \"Max\" >= $2", critType, i1)

	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
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
		fmt.Println(err)
	}

	outputStr := "Description: "
	outputStr += description
	outputStr += "\n\n"
	outputStr += "Effect: "
	outputStr += effect

	b.Send(m.Chat, outputStr)

}

func rollHealth(b *tb.Bot, m *tb.Message) {

	s := strings.Split(m.Payload, " ")
	conMod, _ := strconv.Atoi(s[0])
	totalHealth := 0

	output := ""

	for index, value := range s {

		if index != 0 {

			classHpString := strings.Split(value, ":")

			hpDice := getClassHP(classHpString[0])
			numberOfRolls, _ := strconv.Atoi(classHpString[1])

			if index == 1 {
				output += "Your result = ("
				output += strconv.Itoa(hpDice + conMod)
				totalHealth += (hpDice + conMod)
				numberOfRolls--
			}

			for roll := 0; roll < numberOfRolls; roll++ {
				hp := rand.Intn(hpDice) + 1

				fmt.Print("Hp Roll : ")
				fmt.Println(hp)

				hp += conMod

				if hp <= 0 {
					hp = 1
				}

				totalHealth += hp

				output += " + "
				output += strconv.Itoa(hp)
			}
		}
	}

	output += " ) = "
	output += strconv.Itoa(totalHealth)
	b.Send(m.Chat, output)
}

func getClassHP(name string) int {

	className := strings.ToLower(name)

	if className == "wiz" || className == "sorc" {
		fmt.Println(className + " : 6")
		return 6
	} else if className == "bard" || className == "cleric" || className == "druid" || className == "monk" || className == "rogue" || className == "war" || className == "arti" {
		fmt.Println(className + " : 8")
		return 8
	} else if className == "fighter" || className == "pala" || className == "rang" {
		fmt.Println(className + " : 10")
		return 10
	} else if className == "barb" {
		fmt.Println(className + " : 12")
		return 12
	} else {
		return 0
	}
}

func getClassNames(b *tb.Bot, m *tb.Message) {

	outputStr := "These are the only short forms: \n    Artificer = arti \n    Barbarian = barb \n    Paladin = pala \n    Ranger = rang \n    Sorcerer = sorc \n    Warlock = war \n    Wizard = wiz"

	b.Send(m.Chat, outputStr)
}

func roll(b *tb.Bot, m *tb.Message) {
	s := strings.Split(m.Payload, " ")

	output := ""
	totalRoll := 0

	for index, value := range s {
		diceAndRolls := strings.Split(value, ":")

		if index == 0 {
			output += "Your rolls : "
		}

		dice, _ := strconv.Atoi(diceAndRolls[0])

		numberOfRolls := 0
		if len(diceAndRolls) < 2 {
			numberOfRolls = 1
		} else {
			tempNum, err := strconv.Atoi(diceAndRolls[1])

			if err != nil {
				tempNum = 1
			}
			numberOfRolls = tempNum
		}

		for roll := 0; roll < numberOfRolls; roll++ {
			if roll == 0 {
				output += "("
			}
			eachRoll := rand.Intn(dice) + 1
			totalRoll += eachRoll
			output += strconv.Itoa(eachRoll)

			if roll < numberOfRolls-1 {
				output += ", "
			}

		}

		output += ") "

	}
	outputStr := "Your total rolls: "
	outputStr += strconv.Itoa(totalRoll)
	outputStr += "\n"
	outputStr += output

	b.Send(m.Chat, outputStr)

}
