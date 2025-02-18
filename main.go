package main

import (
	"motoClubeMourosBackend/joia"
	"motoClubeMourosBackend/member"
	StripleHandler "motoClubeMourosBackend/stripleHandler"
	"net/http"
	"os"
	"time"

	"github.com/Maruqes/Tokenize"
	login "github.com/Maruqes/Tokenize/Login"
	"github.com/Maruqes/Tokenize/StripeFunctions"
	"github.com/Maruqes/Tokenize/UserFuncs"
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

func pagarSubscricao(w http.ResponseWriter, r *http.Request) {
	if !login.CheckToken(r) {
		w.Write([]byte("Não Logadoo"))
		return
	}
	w.Write([]byte("Pagamento"))

	userId, err := login.GetIdWithRequest(r)
	if err != nil {
		http.Error(w, "Erro ao obter id do utilizador", http.StatusInternalServerError)
		return
	}

	user, err := UserFuncs.GetUserByID(userId)
	if err != nil {
		http.Error(w, "Erro ao obter utilizador", http.StatusInternalServerError)
		return
	}

	if user.IsProhibited || user.IsActive {
		http.Error(w, "Utilizador proibido ou já ativo", http.StatusForbidden)
		return
	}

	startingDate, err := getMourosStartDate()
	if err != nil {
		http.Error(w, "Erro ao obter data de inicio", http.StatusInternalServerError)
		return
	}

	endingDate, err := getMourosEndingDate()
	if err != nil {
		http.Error(w, "Erro ao obter data de fim", http.StatusInternalServerError)
		return
	}

	nowDate := time.Now()

	if nowDate.After(startingDate) && nowDate.Before(endingDate) {
		StripleHandler.PagamentoDentroDoPrazo(w, r, user)
		return
	} else {
		StripleHandler.PagamentoForaDoPrazo(w, r, user)
		return
	}

}

func main() {
	db := Tokenize.Initialize()

	member.CreateSociosTable(db)
	joia.CreateJoiaTable(db)

	StripleHandler.SUB_PRICE_ID = StripleHandler.GetPriceId()

	StripeFunctions.SetCreateSubscriptionPageCallback(StripleHandler.PagamentoDentroDoPrazoCallBack)
	StripeFunctions.SetOtherEventCallback(StripleHandler.HandleOtherEvents)

	http.HandleFunc("/pagarSubcricao", pagarSubscricao)

	if os.Getenv("DEV") == "True" {
		http.FileServer(http.Dir("public"))
	}

	Tokenize.InitListen("4242")
}
