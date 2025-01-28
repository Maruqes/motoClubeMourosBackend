package main

import (
	"motoClubeMourosBackend/LatePayments"
	stripewebhook "motoClubeMourosBackend/StripeWebhook"
	"net/http"
	"os"

	"github.com/Maruqes/Tokenize"
	funchooks "github.com/Maruqes/Tokenize/FuncHooks"
	login "github.com/Maruqes/Tokenize/Login"
	types "github.com/Maruqes/Tokenize/Types"
)

func testLogado(w http.ResponseWriter, r *http.Request) {
	if !login.CheckToken(r) {
		w.Write([]byte("Não Logadoo"))
		return
	}
	w.Write([]byte("Logado"))
}

func testPago(w http.ResponseWriter, r *http.Request) {
	if !login.CheckToken(r) {
		w.Write([]byte("Não Logado"))
		return
	}
	if val, err := login.IsUserActiveRequest(r); err != nil || !val {
		w.Write([]byte("Não Ativo"))
		return
	}
	w.Write([]byte("Pago"))
}

func main() {
	Tokenize.Initialize()

	http.HandleFunc("/logado", testLogado)
	http.HandleFunc("/pago", testPago)
	// http.HandleFunc("/hasLatePayments", testPago)

	// UserFuncs.ProhibitUser(0)
	// UserFuncs.UnprohibitUser(0)

	//retornam true para cancelar o evento
	funchooks.SetCheckout_UserFunc(LatePayments.CheckIfUserHasLatePaymentsRequest)
	funchooks.SetStripeWebhook_UserFunc(stripewebhook.HandleEvents)

	if os.Getenv("DEV") == "True" {
		http.FileServer(http.Dir("public"))
	}
	Tokenize.InitListen("10951", "/sucess", "/cancel", types.TypeOfSubscriptionValues.MourosSubscription, []types.ExtraPayments{types.ExtraPaymentsValues.Multibanco})
}
