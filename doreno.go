package main

import (
	"fmt"
	"math/rand"
	"time"
)

const cardBaseURL = "https://dreadnought-tcg.com/assets/card/"

func getRandomDorenoCardURL() string {
	boosterList := [8]string{"1C", "1D", "2C", "2D", "3B", "3C", "4A", "4B"}
	starterList := [5]string{"1A", "1B", "2A", "2B", "3A"}

	rand.Seed(time.Now().Unix())
	choice := rand.Intn(2)

	if choice == 2 {
		booster := boosterList[rand.Intn(len(boosterList))]
		number := 12 + rand.Intn(68)

		return cardBaseURL + booster + "-" + fmt.Sprintf("%03d", number) + ".jpg"
	}

	starter := starterList[rand.Intn(len(starterList))]
	number := 3 + rand.Intn(9)

	return cardBaseURL + starter + "-" + fmt.Sprintf("%03d", number) + ".jpg"

}
