package order

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DenisquaP/yandex_gophermart/internal/models"
	modelsOrder "github.com/DenisquaP/yandex_gophermart/internal/models/orders"
	"net/http"
	"net/url"
	"time"

	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	"github.com/DenisquaP/yandex_gophermart/internal/models/customErrors"
	"github.com/DenisquaP/yandex_gophermart/internal/service/validation"
)

func (o *Order) UploadOrder(ctx context.Context, userId int, order string) (int, error) {
	userIdOrder, err := o.db.GetOrder(ctx, order)
	if err != nil {
		return 0, err
	}

	if userIdOrder != 0 {
		// check if user uploaded by other user
		if userIdOrder != userId {
			logger.Logger.Errorw("project does not belong to user", "userId", userId, "order", order)
			cErr := customErrors.NewCustomError("project does not belong to user", http.StatusConflict)

			return 0, cErr
		}

		// check if user already uploaded
		if userIdOrder == userId {
			return http.StatusOK, nil
		}
	}

	// check order for valid
	if !validation.IsValidLuhnNumber(order) {
		logger.Logger.Errorw("invalid order number", "userId", userId, "order", order)
		cErr := customErrors.NewCustomError("invalid order number", http.StatusUnprocessableEntity)

		return 0, cErr
	}

	_, statusCode, err := sendRequest(order)
	if err != nil {
		if statusCode != http.StatusTooManyRequests {
			return 0, err
		}

		for err != nil {
			// if to many requests > trying to send request every second
			_, statusCode, err = sendRequest(order)
			if statusCode != http.StatusTooManyRequests {
				return 0, err
			}

			t := time.After(time.Second)

			<-t
		}
	}

	go o.updateStatusInDB(context.Background(), order)

	return http.StatusAccepted, nil
}

// Sends request to accural system
func sendRequest(order string) (modelsOrder.Order, int, error) {
	client := http.Client{Timeout: 5 * time.Second}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Path: fmt.Sprintf(`http://%s/api/orders/%s`, accuralUrl, order)},
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Logger.Errorw("error sending request to accural system", "error", err)

		return modelsOrder.Order{}, 0, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Logger.Errorw("error sending request to accural system", "error", err)
		}
	}()

	var orderStruct modelsOrder.Order
	if err := json.NewDecoder(resp.Body).Decode(&orderStruct); err != nil {
		logger.Logger.Errorw("error unmarshalling json", "error", err)

		return modelsOrder.Order{}, 0, err
	}

	return orderStruct, resp.StatusCode, nil
}

// Use for update order`s status code in database
func (o *Order) updateStatusInDB(ctx context.Context, order string) {
	var lastUpdateStatus string

	// until status != PROCESSED or INVALID
	for lastUpdateStatus != models.PROCESSED {
		// wait for 5 seconds for another try
		t := time.After(5 * time.Second)
		<-t

		orderStruct, statusCode, err := sendRequest(order)
		// if status code not StatusTooManyRequests returning error
		if statusCode != http.StatusTooManyRequests {
			return
		}

		// if to many requests > trying to send request every second
		for err != nil {
			_, statusCode, err = sendRequest(order)
			if statusCode != http.StatusTooManyRequests {
				return
			}
			t := time.After(time.Second)

			<-t
		}

		// updating status of order
		if err := o.db.UpdateStatus(ctx, order, orderStruct.Status); err != nil {
			return
		}

		if orderStruct.Status == models.INVALID {
			return
		}
		lastUpdateStatus = orderStruct.Status
	}
}
