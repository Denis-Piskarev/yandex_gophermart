package order

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/DenisquaP/yandex_gophermart/internal/logger"
	"github.com/DenisquaP/yandex_gophermart/internal/models"
	"github.com/DenisquaP/yandex_gophermart/internal/models/customerrors"
	modelsOrder "github.com/DenisquaP/yandex_gophermart/internal/models/orders"
	"github.com/DenisquaP/yandex_gophermart/internal/service/validation"
)

func (o *Order) UploadOrder(ctx context.Context, userID int, order string) (int, error) {
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
	if !validation.IsValidLuhnNumber(order) {
		logger.Logger.Errorw("invalid order number", "userID", userID, "order", order)
		cErr := customerrors.NewCustomError("invalid order number", http.StatusUnprocessableEntity)

		return 0, cErr
	}

	orderSt, statusCode, err := sendRequest(true, order)
	if err != nil {
		if statusCode != http.StatusTooManyRequests {
			return 0, err
		}

		for err != nil {
			t := time.After(time.Second)
			<-t

			// if to many requests > trying to send request every second
			_, statusCode, err = sendRequest(false, order)
			if statusCode != http.StatusTooManyRequests {
				return 0, err
			}
		}
	}

	if err := o.db.UploadOrder(ctx, userID, &orderSt); err != nil {
		return 0, err
	}

	// if order is not processed starting goroutine to check and update it
	if orderSt.Status != models.PROCESSED {
		go o.updateStatusInDB(context.Background(), order)
	}

	return http.StatusAccepted, nil
}

// Sends request to accrual system
func sendRequest(first bool, order string) (modelsOrder.OrderAccrual, int, error) {
	// register order in system
	if first {
		var err error

		codeRegister, err := registerInSystem()
		if !errors.Is(err, fmt.Errorf("not expected status code: %d", http.StatusTooManyRequests)) {
			return modelsOrder.OrderAccrual{}, codeRegister, err
		}

		for codeRegister == http.StatusTooManyRequests {
			t := time.After(30 * time.Millisecond)
			<-t

			codeRegister, err = registerInSystem()
			if err != nil {
				if codeRegister == http.StatusTooManyRequests {
					err = nil
				}
			}
		}

		codeUpload, err := uploadOrderToSystem(order)
		if !errors.Is(err, fmt.Errorf("not expected status code uploading order to system: %d", http.StatusTooManyRequests)) {
			return modelsOrder.OrderAccrual{}, codeRegister, err
		}

		for codeUpload == http.StatusTooManyRequests {
			t := time.After(30 * time.Millisecond)
			<-t
			codeUpload, err = uploadOrderToSystem(order)
			if err != nil {
				if codeUpload == http.StatusTooManyRequests {
					err = nil
				}
			}
		}

	}

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/orders/%s", accuralURL, order), nil)
	if err != nil {
		return modelsOrder.OrderAccrual{}, http.StatusInternalServerError, err
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Logger.Errorw("error sending request to accrual system", "error", err)

		return modelsOrder.OrderAccrual{}, 0, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Logger.Errorw("error sending request to accrual system", "error", err)
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
func (o *Order) updateStatusInDB(ctx context.Context, order string) {
	var lastUpdateStatus string

	// until status != PROCESSED or INVALID
	for lastUpdateStatus != models.PROCESSED {
		// wait for 5 seconds for another try
		t := time.After(2 * time.Second)
		<-t

		orderStruct, statusCode, err := sendRequest(false, order)
		// if status code not StatusTooManyRequests returning error
		if statusCode != http.StatusTooManyRequests {
			return
		}

		// if to many requests > trying to send request every second
		for err != nil {
			_, statusCode, err = sendRequest(false, order)
			if statusCode != http.StatusTooManyRequests {
				return
			}
			t := time.After(30 * time.Millisecond)

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

func registerInSystem() (int, error) {
	newGoods := fmt.Sprintf("%d", rand.Intn(1000000000))

	client := http.Client{Timeout: 5 * time.Second}
	body1 := struct {
		Match      string `json:"match"`
		Reward     int    `json:"reward"`
		RewardType string `json:"reward_type"`
	}{
		Match:      newGoods,
		Reward:     10,
		RewardType: "%",
	}
	jBody1, err := json.Marshal(body1)
	if err != nil {
		logger.Logger.Errorw("error marshalling json", "error", err)

		return 0, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/goods", accuralURL), bytes.NewBuffer(jBody1))
	if err != nil {
		logger.Logger.Errorw("error request to accrual system", "error", err)

		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		logger.Logger.Errorw("error request to accrual system register order", "error", err)

		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Logger.Errorw("not expected status code")

		return resp.StatusCode, fmt.Errorf("not expected status code: %d", resp.StatusCode)
	}

	return resp.StatusCode, nil
}

func uploadOrderToSystem(order string) (int, error) {
	newName := fmt.Sprintf("%d", rand.Intn(1000000000))

	client := http.Client{Timeout: 5 * time.Second}

	type goods struct {
		Description string `json:"description"`
		Price       int    `json:"price"`
	}

	body := struct {
		Order string  `json:"order"`
		Goods []goods `json:"goods"`
	}{
		Order: order,
		Goods: []goods{
			{
				Description: newName + "saw",
				Price:       rand.Intn(10000),
			},
		},
	}
	jBody, err := json.Marshal(body)
	if err != nil {
		logger.Logger.Errorw("error marshalling json", "error", err)

		return 0, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/orders", accuralURL), bytes.NewBuffer(jBody))
	if err != nil {
		logger.Logger.Errorw("error request to accrual system", "error", err)

		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		logger.Logger.Errorw("error request to accrual system register order", "error", err)

		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		logger.Logger.Errorw("not expected status code")

		return resp.StatusCode, fmt.Errorf("not expected status code uploading order to system: %d", resp.StatusCode)
	}

	return resp.StatusCode, nil
}
