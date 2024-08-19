package endpoints

import (
	"encoding/json"
	"errors"
	"github.com/DenisquaP/yandex_gophermart/internal/models/customErrors"
	"net/http"

	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	userModels "github.com/DenisquaP/yandex_gophermart/internal/models/users"
)

// RegisterUser - registers user
func (e *Endpoints) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var request userModels.RegisterReq

	// reading body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Logger.Errorw("error while decoding body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			logger.Logger.Errorw("error closing body", "error", err)
		}
	}()

	if !checkEmpty(request.Login, request.Password) {
		http.Error(w, "empty login or password", http.StatusBadRequest)

		return
	}

	token, err := e.services.RegisterUser(r.Context(), request.Login, request.Password)
	if err != nil {
		// If we can convert error to custom error
		// We paste status code from custom error
		var cErr customErrors.CustomError
		if errors.As(err, &cErr) {
			http.Error(w, cErr.Error(), cErr.StatusCode)

			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "gopherToken",
		Value:  token,
		Path:   "/user",
		Domain: "localhost",
	})
}

// LoginUser - use fro login user
func (e *Endpoints) LoginUser(w http.ResponseWriter, r *http.Request) {
	var request userModels.RegisterReq

	// reading body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Logger.Errorw("error while decoding body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			logger.Logger.Errorw("error closing body", "error", err)
		}
	}()

	if !checkEmpty(request.Login, request.Password) {
		http.Error(w, "empty login or password", http.StatusBadRequest)

		return
	}

	token, err := e.services.LoginUser(r.Context(), request.Login, request.Password)
	if err != nil {
		// If we can convert error to custom error
		// We paste status code from custom error
		var cErr customErrors.CustomError
		if errors.As(err, &cErr) {
			http.Error(w, cErr.Error(), cErr.StatusCode)

			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "gopherToken",
		Value: token,
	})
}

// GetWithdrawals - gets user`s withdrawals
func (e *Endpoints) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	withdrawals, err := e.services.GetWithdrawals(r.Context(), userId)
	if err != nil {
		var cErr customErrors.CustomError
		// if we got custom err then set status code from err
		if errors.As(err, &cErr) {
			http.Error(w, cErr.Error(), cErr.StatusCode)

			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	wByte, err := json.Marshal(withdrawals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(wByte); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

// Checks if login or password is empty
func checkEmpty(login, password string) bool {
	if login == "" || password == "" {
		logger.Logger.Errorw("empty login or password", "login", login, "password", password)

		return false
	}

	return true
}
