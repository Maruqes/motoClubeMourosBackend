package main

import (
	"motoClubeMourosBackend/member"
	"net/http"
	"os"
	"time"

	"github.com/Maruqes/Tokenize"
	login "github.com/Maruqes/Tokenize/Login"
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

func pagarSubscricao(w http.ResponseWriter, r *http.Request) {
	if !login.CheckToken(r) {
		w.Write([]byte("Não Logadoo"))
		return
	}
	w.Write([]byte("Pagamento"))
	startingDate, err := getMourosStartDate()
	if err != nil {
		panic(err)
	}
	endingDate, err := getMourosEndingDate()
	if err != nil {
		panic(err)
	}

	nowDate := time.Now()

	if nowDate.After(startingDate) && nowDate.Before(endingDate) {
		pagamentoDentroDoPrazo(w, r)
		return
	} else {
		pagamentoForaDoPrazo(w, r)
		return
	}

}

func main() {
	db := Tokenize.Initialize()

	member.CreateSociosTable(db)

	http.HandleFunc("/logado", testLogado)
	http.HandleFunc("/pagarSubcricao", pagarSubscricao)

	if os.Getenv("DEV") == "True" {
		http.FileServer(http.Dir("public"))
	}

	Tokenize.InitListen("4242", "/sucesso", "/cancelado")
}
