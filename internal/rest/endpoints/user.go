package endpoints

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	"github.com/DenisquaP/yandex_gophermart/internal/models"
	userModels "github.com/DenisquaP/yandex_gophermart/internal/models/users"
)

func (e *Endpoints) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var request userModels.RegisterReq

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Logger.Errorw("error while decoding body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if request.Login == "" || request.Password == "" {
		logger.Logger.Errorw("empty login or password", "login", request.Login, "password", request.Password)
		http.Error(w, "empty login or password", http.StatusBadRequest)

		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			logger.Logger.Errorw("error closing body", "error", err)
		}
	}()

	if err := e.services.RegisterUser(r.Context(), request.Login, request.Password); err != nil {
		var cErr models.CustomError
		if errors.As(err, &cErr) {
			http.Error(w, cErr.Error(), cErr.StatusCode)
		}

		logger.Logger.Errorw("error registering user", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
}
