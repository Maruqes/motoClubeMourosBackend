package main

import "net/http"

func pagamentoDentroDoPrazo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pagamento dentro do prazo"))
}
