package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Maruqes/Tokenize"
	funchooks "github.com/Maruqes/Tokenize/FuncHooks"
	functions "github.com/Maruqes/Tokenize/Functions"
	login "github.com/Maruqes/Tokenize/Login"
	types "github.com/Maruqes/Tokenize/Types"
	"github.com/Maruqes/Tokenize/UserFuncs"
	"github.com/Maruqes/Tokenize/database"
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

// needs to cheak all code for offline payments and user funcs
func checkIfUserHasLatePayments(id int) (bool, error) {
	endDate, err := UserFuncs.GetEndDateForUser(id)
	if err != nil {
		return false, err
	}

	if endDate == (database.Date{}) {
		return false, nil
	}

	timeDB := time.Date(endDate.Year, time.Month(endDate.Month), endDate.Day, 0, 0, 0, 0, time.UTC).Unix()
	mourosDate, _ := functions.GetMourosStartingDate()

	timeDB = time.Date(time.Now().Year()-2, mourosDate.Month(), mourosDate.Day(), 0, 0, 0, 0, time.UTC).Unix()

	//ultimo mouros date que passou
	lastMourosDate := time.Date(time.Now().Year(), mourosDate.Month(), mourosDate.Day(), 0, 0, 0, 0, time.UTC).Unix()
	yearDifference := time.Unix(timeDB, 0).Year() - time.Unix(lastMourosDate, 0).Year()
	fmt.Println("yearDifference: ", yearDifference)

	if yearDifference > 0 {
		return true, nil
	}
	
	return false, nil
}

func checkIfUserHasLatePaymentsRequest(w http.ResponseWriter, r *http.Request) bool {
	//checkar se o usuario tem pagamentos em atraso se tiver retornar true, bloquear e pedir para pagar
	return false
}

func main() {
	Tokenize.Initialize()

	date, err := functions.GetMourosStartingDate()
	if err != nil {
		fmt.Println("Error getting mouros date")
		return

	}
	fmt.Println("Mouros date: ", date)

	res, err := checkIfUserHasLatePayments(1)
	if err != nil {
		fmt.Println("Error checking if user has late payments")
		return
	}
	fmt.Println("User has late payments: ", res)

	http.HandleFunc("/logado", testLogado)
	http.HandleFunc("/pago", testPago)
	// http.HandleFunc("/hasLatePayments", testPago)

	// UserFuncs.ProhibitUser(0)
	// UserFuncs.UnprohibitUser(0)
	funchooks.SetCheckout_UserFunc(checkIfUserHasLatePaymentsRequest)

	Tokenize.InitListen("4242", "/sucess", "/cancel", types.TypeOfSubscriptionValues.Normal, []types.ExtraPayments{types.ExtraPaymentsValues.Multibanco})
}
