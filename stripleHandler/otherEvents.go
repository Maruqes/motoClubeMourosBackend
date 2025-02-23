package StripleHandler

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Maruqes/Tokenize/Logs"
	"github.com/Maruqes/Tokenize/database"
	"github.com/stripe/stripe-go/v81"
)

func Custumer_subscription_deleted(event stripe.Event) {
	fmt.Println("custumer_subscription_deleted")
	// active -> 0
	var subscription stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &subscription)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
		return
	}

	user, err := database.GetUserByStripeID(subscription.Customer.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !user.IsActive {
		fmt.Println("user is already inactive")
		return
	}

	err = database.DeactivateUser(user.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func HandleOtherEvents(event stripe.Event) {
	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case "customer.subscription.deleted":
		Custumer_subscription_deleted(event)
	case "invoice.payment_succeeded":
		PagamentoDentroDoPrazoCallBack(event)
	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
		Logs.LogMessage("Unhandled event type: " + string(event.Type))
	}
}
