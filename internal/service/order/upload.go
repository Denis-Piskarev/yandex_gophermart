package order

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	"github.com/DenisquaP/yandex_gophermart/internal/models"
	"github.com/DenisquaP/yandex_gophermart/internal/models/customerrors"
	modelsOrder "github.com/DenisquaP/yandex_gophermart/internal/models/orders"
	"github.com/DenisquaP/yandex_gophermart/internal/service/validation"
)

func (o *Order) UploadOrder(ctx context.Context, userID int, order int) (int, error) {
	userIDOrder, err := o.db.GetOrder(ctx, order)
	if err != nil {
		return 0, err
	}

	if userIDOrder != 0 {
		// check if user uploaded by other user
		if userIDOrder != userID {
			logger.Logger.Errorw("project does not belong to user", "userID", userID, "order", order)
			cErr := customerrors.NewCustomError("project does not belong to user", http.StatusConflict)

			return 0, cErr
		}

		// check if user already uploaded
		if userIDOrder == userID {
			return http.StatusOK, nil
		}
	}

	// check order for valid
	if !validation.IsValidLuhnNumber(strconv.Itoa(order)) {
		logger.Logger.Errorw("invalid order number", "userID", userID, "order", order)
		cErr := customerrors.NewCustomError("invalid order number", http.StatusUnprocessableEntity)

		return 0, cErr
	}

	orderSt, statusCode, err := sendRequest(order)
	if err != nil {
		if statusCode != http.StatusTooManyRequests {
			return 0, err
		}

		for err != nil {
			t := time.After(time.Second)
			<-t

			// if to many requests > trying to send request every second
			_, statusCode, err = sendRequest(order)
			if statusCode != http.StatusTooManyRequests {
				return 0, err
			}
		}
	}

	if err := o.db.UploadOrder(ctx, userID, &orderSt); err != nil {
		return 0, err
	}

	go o.updateStatusInDB(context.Background(), order)

	return http.StatusAccepted, nil
}

// Sends request to accural system
func sendRequest(order int) (modelsOrder.OrderAccrual, int, error) {
	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, fmt.Sprintf("%s/api/orders/%d", accuralURL, order), nil)
	if err != nil {
		return modelsOrder.OrderAccrual{}, http.StatusInternalServerError, err
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Logger.Errorw("error sending request to accural system", "error", err)

		return modelsOrder.OrderAccrual{}, 0, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Logger.Errorw("error sending request to accural system", "error", err)
		}
	}()

	if resp.StatusCode == http.StatusNoContent {
		cErr := customerrors.NewCustomError("no content", http.StatusNoContent)
		return modelsOrder.OrderAccrual{}, http.StatusNoContent, cErr
	}

	var orderStruct modelsOrder.OrderAccrual
	if err := json.NewDecoder(resp.Body).Decode(&orderStruct); err != nil {
		logger.Logger.Errorw("error unmarshalling json", "error", err)

		return modelsOrder.OrderAccrual{}, 0, err
	}

	return orderStruct, resp.StatusCode, nil
}

// Use for update order`s status code in database
func (o *Order) updateStatusInDB(ctx context.Context, order int) {
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
