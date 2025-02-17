package StripleHandler

import "os"

var SUB_PRICE_ID string

func GetPriceId() string {
	priceId := os.Getenv("SUBSCRIPTION_PRICE_ID")
	if priceId == "" {
		panic("SUBSCRIPTION_PRICE_ID is not set")
	}
	return priceId
}
