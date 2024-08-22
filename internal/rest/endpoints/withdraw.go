package endpoints

import (
	"encoding/json"
	"errors"
	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	"github.com/DenisquaP/yandex_gophermart/internal/models/customerrors"
	modelsUser "github.com/DenisquaP/yandex_gophermart/internal/models/users"
	"net/http"
)

func (e *Endpoints) Withdraw(w http.ResponseWriter, r *http.Request) {
	var request modelsUser.Withdrawal
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Logger.Errorw("error while decoding request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	userID, err := getUserIDFromHeader(r)
	if err != nil {
		return
	}

	// withdrawing
	if err := e.services.Withdraw(r.Context(), userID, request.Sum, request.Order); err != nil {
		var cErr customerrors.CustomError
		// if we got custom err then set status code from err
		if errors.As(err, &cErr) {
			http.Error(w, cErr.Error(), cErr.StatusCode)

			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
