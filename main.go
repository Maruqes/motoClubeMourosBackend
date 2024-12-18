package main

import (
	"net/http"

	"github.com/Maruqes/Tokenize"
	login "github.com/Maruqes/Tokenize/Login"
	types "github.com/Maruqes/Tokenize/Types"
)

func testLogado(w http.ResponseWriter, r *http.Request) {
	if !login.CheckToken(r) {
		w.Write([]byte("Não Logado."))
		return
	}
	w.Write([]byte("Logado"))
}

func testPago(w http.ResponseWriter, r *http.Request) {
	if !login.CheckToken(r) {
		w.Write([]byte("Não Logado."))
		return
	}
	if val, err := login.IsUserActiveRequest(r); err != nil || !val {
		w.Write([]byte("Não Ativo."))
		return
	}
	w.Write([]byte("Pago."))
}

func main() {

	http.HandleFunc("/logado", testLogado)
	http.HandleFunc("/pago", testPago)

	Tokenize.Init("10951", "/sucess", "/cancel", types.TypeOfSubscriptionValues.Normal, []types.ExtraPayments{types.ExtraPaymentsValues.Multibanco})
}
