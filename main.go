package main

import (
	"encoding/json"
	"motoClubeMourosBackend/member"
	"net/http"
	"os"

	"github.com/Maruqes/Tokenize"
	login "github.com/Maruqes/Tokenize/Login"
	types "github.com/Maruqes/Tokenize/Types"
)

//contato sos opcional
//ver coisas obrigatorias ou nao

//criar tipo de eventos nome do tipo
//criar perguntas dinamicamente (id, secao, pergunta, tipo de resposta, opcional ou nao)
//resposta (id_membro, id_pergunta, resposta)

/*
//informacos do evento
eventos-> cartaz img
titulo
data do evento
numero incricoes
preco/preço socio
data limite
programa do evento -> horas, descricao
infos extra para o evento
*/

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

/*
um socio pode inserir apenas a sua propria informação :D
*/
func insertMemberInfo(w http.ResponseWriter, r *http.Request) {
	if !login.CheckToken(r) {
		w.Write([]byte("Não Logado"))
		return
	}

	if val, err := login.IsUserActiveRequest(r); err != nil || !val {
		w.Write([]byte("Não Ativo"))
		return
	}
	var m member.Member

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("Member info received"))

	cookie, err := r.Cookie("id")
	if err != nil {
		w.Write([]byte("ID do cookie não encontrado"))
		return
	}
	if cookie.Value != m.ID {
		w.Write([]byte("ID do cookie diferente do ID do membro"))
		return
	}

	member.InsertMember(m)
}

func getMemberInfo(w http.ResponseWriter, r *http.Request) {
	if !login.CheckToken(r) {
		w.Write([]byte("Não Logado"))
		return
	}
	if val, err := login.IsUserActiveRequest(r); err != nil || !val {
		w.Write([]byte("Não Ativo"))
		return
	}

	cookie, err := r.Cookie("id")
	if err != nil {
		w.Write([]byte("ID do cookie não encontrado"))
		return
	}

	m, err := member.GetMemberData(cookie.Value)
	if err != nil {
		w.Write([]byte("Membro não encontrado"))
		return
	}

	json.NewEncoder(w).Encode(m)
}

func main() {
	db := Tokenize.Initialize()

	member.CreateSociosTable(db)

	http.HandleFunc("/logado", testLogado)
	http.HandleFunc("/pago", testPago)
	http.HandleFunc("/insertMemberInfo", insertMemberInfo)
	http.HandleFunc("/getMemberInfo", getMemberInfo)
	// http.HandleFunc("/hasLatePayments", testPago)

	// UserFuncs.ProhibitUser(0)
	// UserFuncs.UnprohibitUser(0)

	if os.Getenv("DEV") == "True" {
		http.FileServer(http.Dir("public"))
	}

	Tokenize.InitListen("10951", "/sucess", "/cancel", types.TypeOfSubscriptionValues.MourosSubscription, []types.ExtraPayments{types.ExtraPaymentsValues.Multibanco})
}
