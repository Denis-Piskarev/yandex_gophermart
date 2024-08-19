package endpoints

import (
	"encoding/json"
	"net/http"
)

// GetBalance - gets balance of user
func (e *Endpoints) GetBalance(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	balance, err := e.services.GetBalance(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	balanceByte, err := json.Marshal(balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if _, err := w.Write(balanceByte); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
