package order

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
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
func sendRequest(first bool, order int) (modelsOrder.OrderAccrual, int, error) {
	if first {
		if err := registerInSystem(order); err != nil {
			return modelsOrder.OrderAccrual{}, 0, err
		}
	}

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

	//var orderStruct modelsOrder.OrderAccrual
	var res struct {
		Order   string  `json:"order"`
		Status  string  `json:"status"`
		Accrual float32 `json:"accrual"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		logger.Logger.Errorw("error unmarshalling json", "error", err)

		return modelsOrder.OrderAccrual{}, 0, err
	}

	orderInt, err := strconv.Atoi(res.Order)
	if err != nil {
		logger.Logger.Errorw("error unmarshalling json", "error", err)

		return modelsOrder.OrderAccrual{}, 0, err
	}

	return modelsOrder.OrderAccrual{
		Order:   orderInt,
		Status:  res.Status,
		Accrual: res.Accrual,
	}, resp.StatusCode, nil
}

// Use for update order`s status code in database
func (o *Order) updateStatusInDB(ctx context.Context, order int) {
	var lastUpdateStatus string

	// until status != PROCESSED or INVALID
	for lastUpdateStatus != models.PROCESSED {
		// wait for 5 seconds for another try
		t := time.After(5 * time.Second)
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

func registerInSystem(order int) error {
	newName := fmt.Sprintf("%d", rand.Intn(1000000000))

	client := http.Client{Timeout: 5 * time.Second}
	body1 := struct {
		Match      string `json:"match"`
		Reward     int    `json:"reward"`
		RewardType string `json:"reward_type"`
	}{
		Match:      newName,
		Reward:     10,
		RewardType: "%",
	}
	jBody1, err := json.Marshal(body1)
	if err != nil {
		logger.Logger.Errorw("error marshalling json", "error", err)

		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/goods", accuralURL), bytes.NewBuffer(jBody1))
	if err != nil {
		logger.Logger.Errorw("error request to accrual system", "error", err)

		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		logger.Logger.Errorw("error request to accrual system register order", "error", err)

		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Logger.Errorw("not expected status code")

		return fmt.Errorf("not expected status code: %d", resp.StatusCode)
	}

	type goods struct {
		Description string `json:"description"`
		Price       int    `json:"price"`
	}

	body2 := struct {
		Order string  `json:"order"`
		Goods []goods `json:"goods"`
	}{
		Order: fmt.Sprintf("%d", order),
		Goods: []goods{
			{
				Description: newName + "saw",
				Price:       rand.Intn(10000),
			},
		},
	}
	jBody2, err := json.Marshal(body2)
	if err != nil {
		logger.Logger.Errorw("error marshalling json", "error", err)

		return err
	}

	req2, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/orders", accuralURL), bytes.NewBuffer(jBody2))
	if err != nil {
		logger.Logger.Errorw("error request to accrual system", "error", err)

		return err
	}
	req2.Header.Set("Content-Type", "application/json")

	resp2, err := client.Do(req2)
	if err != nil {
		logger.Logger.Errorw("error request to accrual system register order", "error", err)

		return err
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusAccepted {
		logger.Logger.Errorw("not expected status code")

		return fmt.Errorf("not expected status code 2: %d", resp2.StatusCode)
	}

	return nil
}
