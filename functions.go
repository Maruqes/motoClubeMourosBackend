package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func getMourosStartDate() (time.Time, error) {
	startDate := os.Getenv("MOUROS_STARTING_DATE")
	if startDate == "" {
		return time.Time{}, fmt.Errorf("MOUROS_STARTING_DATE is not set")
	}
	startDate = startDate + "/" + strconv.Itoa(time.Now().Year())

	return time.Parse("2/1/2006", startDate)
}

func getMourosEndingDate() (time.Time, error) {
	endingDate := os.Getenv("MOUROS_ENDING_DATE")
	if endingDate == "" {
		return time.Time{}, fmt.Errorf("MOUROS_ENDING_DATE is not set")
	}
	endingDate = endingDate + "/" + strconv.Itoa(time.Now().Year())

	return time.Parse("2/1/2006", endingDate)
}

func checkMourosDate() {
	startingDate, err := getMourosStartDate()
	if err != nil {
		panic(err)
	}

	endingDate, err := getMourosEndingDate()
	if err != nil {
		panic(err)
	}

	nowDate := time.Now()

	if nowDate.After(startingDate) && nowDate.Before(endingDate) {
		fmt.Println("Dentro do prazo")
	} else {
		panic("Fora do prazo")
	}
}
