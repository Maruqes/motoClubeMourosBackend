package main

import (
	"motoClubeMourosBackend/LatePayments"
	stripewebhook "motoClubeMourosBackend/StripeWebhook"
	"motoClubeMourosBackend/member"
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
	db := Tokenize.Initialize()

	member.CreateSociosTable(db)
	member.ChangeMember(member.Member{ // Exemplo de inserção de um membro
		ID:                "1",
		NumeroSocio:       1,
		Junior:            false,
		SocioResponsavel:  "",
		DataNascimento:    "2000-01-01",
		DataAdesao:        "2021-01-01",
		MembroResponsavel: "1",
		Nome:              "João",
		Email:             "testeemail@teste.com",
		Telefone:          "123456789",
		TipoSangue:        "B+",
		Rua:               "Rua",
		Numero:            "1",
		Concelho:          "Concelho",
		Distrito:          "Distrito",
		CodPostal:         "1234-567",
		Tipo:              "Sócio",
		GrupoWA:           false,
		DataInscricao:     "2021-01-02",
	})

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
