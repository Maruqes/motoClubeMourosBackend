package StripleHandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Maruqes/Tokenize/StripeFunctions"
	"github.com/Maruqes/Tokenize/database"
	"github.com/stripe/stripe-go/v81"
)

func activateWithInvoice(invoice stripe.Invoice) error {
	user, err := database.GetUserByStripeID(invoice.Customer.ID)
	if err != nil {
		return err
	}

	if user.IsActive {
		return nil
	}

	err = database.ActivateUser(user.ID)
	if err != nil {
		return err
	}

	return nil
}

func PagamentoDentroDoPrazoCallBack(event stripe.Event) {
	if event.Type != "invoice.payment_succeeded" {
		log.Println("event is not invoice.payment_succeeded")
		return
	}

	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		log.Printf("Error parsing webhook JSON: %v\n", err)
		return
	}

	if invoice.Subscription == nil {
		log.Println("invoice has no subscription")
		return
	}

	err := activateWithInvoice(invoice)
	if err != nil {
		fmt.Println(err)
	}
}

func PagamentoDentroDoPrazo(w http.ResponseWriter, r *http.Request, user database.User) {

	checkout, err := StripeFunctions.CreateSubscriptionPage(user.ID, SUB_PRICE_ID, map[string]string{}, "http://localhost:4242/sucess", "http://localhost:4242/cancel")
	if err != nil {
		http.Error(w, "Erro ao criar checkout", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, checkout.URL, http.StatusSeeOther)
}
