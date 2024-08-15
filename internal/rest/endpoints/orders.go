package endpoints

import (
	"encoding/json"
	"net/http"
)

func (e *Endpoints) GetOrders(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	orders, err := e.services.GetOrders(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if len(orders) == 0 {
		http.Error(w, "No orders found", http.StatusNoContent)

		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if _, err := w.Write(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
