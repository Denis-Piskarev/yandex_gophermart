package endpoints

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	"github.com/DenisquaP/yandex_gophermart/internal/models/customerrors"
)

// UploadOrder - uploads order to system
func (e *Endpoints) UploadOrder(w http.ResponseWriter, r *http.Request) {
	userID, err := getuserIDFromHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	var order string
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		logger.Logger.Errorw("error unmarshalling request", "error", err)

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	statusCode, err := e.services.UploadOrder(r.Context(), userID, order)
	if err != nil {
		var cErr customerrors.CustomError
		// if we got custom err then set status code from err
		if errors.As(err, &cErr) {
			http.Error(w, cErr.Error(), cErr.StatusCode)

			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
}

// GetOrders - gets order info
func (e *Endpoints) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := getuserIDFromHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	orders, err := e.services.GetOrders(r.Context(), userID)
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
