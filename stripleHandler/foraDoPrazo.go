package StripleHandler

import (
	"net/http"

	"github.com/Maruqes/Tokenize/database"
)

func PagamentoForaDoPrazo(w http.ResponseWriter, r *http.Request, user database.User) {
	w.Write([]byte("Pagamento fora do prazo"))
}
