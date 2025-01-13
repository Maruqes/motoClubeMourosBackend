package LatePayments

import (
	"fmt"
	"log"
	mourosFunctions "motoClubeMourosBackend/MourosFunctions"
	"net/http"
	"time"

	functions "github.com/Maruqes/Tokenize/Functions"
	"github.com/Maruqes/Tokenize/UserFuncs"
	"github.com/Maruqes/Tokenize/database"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
)

// needs to cheak all code for offline payments and user funcs
func CheckIfUserHasLatePayments(id int) (bool, []int64, error) {
	endDate, err := UserFuncs.GetEndDateForUser(id)
	if err != nil {
		return false, []int64{}, err
	}

	if endDate == (database.Date{}) {
		return false, []int64{}, nil
	}

	timeDB := time.Date(endDate.Year, time.Month(endDate.Month), endDate.Day, 0, 0, 0, 0, time.UTC).Unix()

	mourosDate, _ := functions.GetMourosStartingDate()
	lastMourosDate := time.Date(time.Now().Year(), mourosDate.Month(), mourosDate.Day(), 0, 0, 0, 0, time.UTC).Unix()

	yearDifference := time.Unix(timeDB, 0).Year() - time.Unix(lastMourosDate, 0).Year()
	fmt.Println("yearDifference: ", yearDifference)
	fmt.Println("timeDB: ", time.Unix(timeDB, 0))
	fmt.Println("lastMourosDate: ", time.Unix(lastMourosDate, 0))

	if yearDifference < 0 {
		yearsLeft := []int64{}
		for i := 0; i < -yearDifference; i++ {
			yearsLeft = append(yearsLeft, time.Date(int(time.Unix(timeDB, 0).Year())+i, 1, 1, 0, 0, 0, 0, time.UTC).Unix())
		}

		return true, yearsLeft, nil
	}

	return false, []int64{}, nil
}

func sendToCheckOutPaymentStripe(w http.ResponseWriter, r *http.Request, numberOfPayments int) {
	priceOne := mourosFunctions.GetSubscriptionPrice() - mourosFunctions.GetCouponDiscount()

	checkoutParams := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card", "multibanco"}),
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("eur"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Late Payments/Anos em Falta"),
					},
					UnitAmount: stripe.Int64(priceOne * int64(numberOfPayments)),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Metadata: map[string]string{
			"LatePayments":     "true",
			"numberOfPayments": fmt.Sprintf("%d", numberOfPayments),
		},
		SuccessURL: stripe.String("http://localhost:3000"),
		CancelURL:  stripe.String("http://localhost:3000"),
	}

	session, err := session.New(checkoutParams)
	if err != nil {
		http.Error(w, "Failed to create checkout session", http.StatusInternalServerError)
		log.Printf("Failed to create checkout session: %v", err)
		return
	}

	// Redirecionar para a página de pagamento
	http.Redirect(w, r, session.URL, http.StatusSeeOther)
	log.Printf("Redirecionado para URL: %s", session.URL)
}

func CheckIfUserHasLatePaymentsRequest(w http.ResponseWriter, r *http.Request) bool {
	res, numberOfYears, err := CheckIfUserHasLatePayments(1)
	if err != nil {
		fmt.Println("Error checking if user has late payments")
		return false
	}
	if res {
		sendToCheckOutPaymentStripe(w, r, len(numberOfYears))
		return true
	}
	return false
}
