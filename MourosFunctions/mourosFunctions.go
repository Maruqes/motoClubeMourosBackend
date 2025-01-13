package mourosFunctions

import (
	"log"
	"os"

	"github.com/stripe/stripe-go/v81/coupon"
	"github.com/stripe/stripe-go/v81/price"
)

func GetSubscriptionPrice() int64 {

	priceID := os.Getenv("SUBSCRIPTION_PRICE_ID")
	if priceID == "" {
		log.Fatalf("ID do preço da assinatura não encontrado")
	}

	p, err := price.Get(priceID, nil)
	if err != nil {
		log.Fatalf("Erro ao obter o preço: %v", err)
	}

	return int64(p.UnitAmount)
}

func GetCouponDiscount() int64 {

	couponID := os.Getenv("COUPON_ID")
	if couponID == "" {
		log.Fatalf("ID do cupom não encontrado")
	}

	c, err := coupon.Get(couponID, nil)
	if err != nil {
		log.Fatalf("Erro ao obter o cupom: %v", err)
	}

	if c.AmountOff == 0 {
		return GetSubscriptionPrice() * int64(c.PercentOff) / 100
	}

	return int64(c.AmountOff)
}
