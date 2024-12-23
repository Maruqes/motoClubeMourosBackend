package main

import (
	"net/http"
	"time"

	"github.com/Maruqes/Tokenize"
	funchooks "github.com/Maruqes/Tokenize/FuncHooks"
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

	timeNow := time.Now().Unix()
	timeDB := time.Date(endDate.Year, time.Month(endDate.Month), endDate.Day, 0, 0, 0, 0, time.UTC).Unix()

	if timeDB+(365*24*60*60) < timeNow {
		return true, nil
	}

	return false, nil
}

func checkIfUserHasLatePaymentsRequest(w http.ResponseWriter, r *http.Request) bool {
	//checkar se o usuario tem pagamentos em atraso se tiver retornar true, bloquear e pedir para pagar
	return false
}

func main() {

	http.HandleFunc("/logado", testLogado)
	http.HandleFunc("/pago", testPago)
	// http.HandleFunc("/hasLatePayments", testPago)

	// UserFuncs.ProhibitUser(0)
	// UserFuncs.UnprohibitUser(0)
	funchooks.SetCheckout_UserFunc(checkIfUserHasLatePaymentsRequest)

	Tokenize.Init("10951", "/sucess", "/cancel", types.TypeOfSubscriptionValues.Normal, []types.ExtraPayments{types.ExtraPaymentsValues.Multibanco})
}
