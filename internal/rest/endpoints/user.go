package endpoints

import (
	"encoding/json"
	"net/http"

	userModels "github.com/DenisquaP/yandex_gophermart/internal/models/users"
)

func (e *Endpoints) registerUser(w http.ResponseWriter, r *http.Request) {
	var request userModels.RegisterReq

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		e.logger.Errorw("error while decoding body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if request.Login == "" || request.Password == "" {
		e.logger.Errorw("empty login or password", "login", request.Login, "password", request.Password)
		http.Error(w, "empty login or password", http.StatusBadRequest)

		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			e.logger.Errorw("error closing body", "error", err)
		}
	}()

	if err := e.services.RegisterUser(request.Login, request.Password); err != nil {
		e.logger.Errorw("error registering user", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
}
