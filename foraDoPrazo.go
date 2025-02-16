package main

import "net/http"

func pagamentoForaDoPrazo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pagamento fora do prazo"))
}
