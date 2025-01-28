package stripewebhook

import (
	"github.com/stripe/stripe-go/v81"
)

// retorna true para cancelar o evento
func HandleEvents(event stripe.Event) bool {

	switch event.Type {
	case "customer.subscription.deleted":
	case "customer.subscription.created":
	case "customer.created":
	case "invoice.payment_succeeded":
	case "charge.succeeded":
	case "invoice.created":
	default:
		return false
	}
	return false
}
