package main

import (
	"fmt"
	"motoClubeMourosBackend/LatePayments"
	stripewebhook "motoClubeMourosBackend/StripeWebhook"
	"net/http"
	"os"
	"time"

	"github.com/Maruqes/Tokenize"
	funchooks "github.com/Maruqes/Tokenize/FuncHooks"
	functions "github.com/Maruqes/Tokenize/Functions"
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

	date, err := functions.GetMourosStartingDate()
	if err != nil {
		fmt.Println("Error getting mouros date")
		return
	}
	fmt.Println("Mouros date: ", date)
	fmt.Println()
	fmt.Println()

	res, numberOfYears, err := LatePayments.CheckIfUserHasLatePayments(2)
	if err != nil {
		fmt.Println("Error checking if user has late payments")
	}
	if res {
		for _, year := range numberOfYears {
			fmt.Println()
			fmt.Println(time.Unix(year, 0).Format("01/12/2006"))
		}
	}
	fmt.Println()
	fmt.Println()

	http.HandleFunc("/logado", testLogado)
	http.HandleFunc("/pago", testPago)
	// http.HandleFunc("/hasLatePayments", testPago)

	// UserFuncs.ProhibitUser(0)
	// UserFuncs.UnprohibitUser(0)
	funchooks.SetCheckout_UserFunc(LatePayments.CheckIfUserHasLatePaymentsRequest)
	funchooks.SetStripeWebhook_UserFunc(stripewebhook.HandleEvents)

	if os.Getenv("DEV") == "True" {
		http.FileServer(http.Dir("public"))
	}
	Tokenize.InitListen("10951", "/sucess", "/cancel", types.TypeOfSubscriptionValues.MourosSubscription, []types.ExtraPayments{types.ExtraPaymentsValues.Multibanco})
}
