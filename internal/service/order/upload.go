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
	"github.com/DenisquaP/yandex_gophermart/internal/validation"
)

func (o *Order) UploadOrder(ctx context.Context, userID int, order string) (int, error) {
	sc, err := o.validateOrder(ctx, userID, order)
	if err != nil {
		return 0, err
	}
	if sc != 0 {
		return sc, nil
	}

	client := http.Client{}
	orderSt, statusCode, err := processOrderRequest(&client, order)
	if err != nil {
		if statusCode != http.StatusTooManyRequests {
			return 0, err
		}

		for err != nil {
			t := time.After(time.Second)
			<-t

			// if to many requests > trying to send request every second
			_, statusCode, err = processOrderRequest(&client, order)
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
		go o.updateStatusInDB(context.Background(), &client, order)
	}

	return http.StatusAccepted, nil
}

func (o *Order) validateOrder(ctx context.Context, userID int, order string) (int, error) {
	userIDOrder, err := o.db.GetUserIdByOrder(ctx, order)
	if err != nil {
		return 0, err
	}

	if userIDOrder != 0 {
		// check if user uploaded by other user
		if userIDOrder != userID {
			logger.Logger.Errorw("project does not belong to user", "userID", userID, "owner", userIDOrder, "order", order)
			cErr := customerrors.NewCustomError("project does not belong to user", http.StatusConflict)

			return 0, cErr
		}

		// check if user already uploaded
		if userIDOrder == userID {
			return http.StatusOK, nil
		}
	}

	// check order for valid
	if !validation.ValidateLuhn(order) {
		logger.Logger.Errorw("invalid order number", "userID", userID, "order", order)
		cErr := customerrors.NewCustomError("invalid order number", http.StatusUnprocessableEntity)

		return 0, cErr
	}

	return http.StatusAccepted, nil
}

// Sends request to accrual system
func processOrderRequest(client *http.Client, order string) (modelsOrder.OrderAccrual, int, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/orders/%s", accrualURL, order), nil)
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
		if err := addNew(order); err != nil {
			return modelsOrder.OrderAccrual{}, resp.StatusCode, err
		}

		orderS, sc, err := processOrderRequest(client, order)
		orderS.Status = models.PROCESSING
		return orderS, sc, err
	}

	var orderStruct modelsOrder.OrderAccrual
	if err := json.NewDecoder(resp.Body).Decode(&orderStruct); err != nil {
		logger.Logger.Errorw("error unmarshalling json", "error", err)

		return modelsOrder.OrderAccrual{}, 0, err
	}

	return orderStruct, resp.StatusCode, nil
}

// Use for update order`s status code in database
func (o *Order) updateStatusInDB(ctx context.Context, client *http.Client, order string) {
	var lastUpdateStatus string

	// until status != PROCESSED or INVALID
	for lastUpdateStatus != models.PROCESSED {
		// wait for 5 seconds for another try
		t := time.After(2 * time.Second)
		<-t

		orderStruct, statusCode, err := processOrderRequest(client, order)
		// if status code not StatusTooManyRequests returning error
		if statusCode != http.StatusTooManyRequests {
			return
		}

		// if to many requests > trying to send request every second
		for err != nil {
			_, statusCode, err = processOrderRequest(client, order)
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

func addNew(order string) error {
	codeRegister, err := registerInSystem()
	if err != nil && !errors.Is(err, fmt.Errorf("not expected status code: %d", http.StatusTooManyRequests)) {
		cErr := customerrors.NewCustomError(err.Error(), codeRegister)

		return cErr
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
	if err != nil && !errors.Is(err, fmt.Errorf("not expected status code uploading order to system: %d", http.StatusTooManyRequests)) {
		cErr := customerrors.NewCustomError(err.Error(), codeRegister)

		return cErr
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

	return nil
}

func registerInSystem() (int, error) {
	newGoods := fmt.Sprintf("%d", rand.Intn(1000000000))

	client := http.Client{}
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

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/goods", accrualURL), bytes.NewBuffer(jBody1))
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
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Logger.Errorw("error closing response body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		logger.Logger.Errorw("not expected status code")

		return resp.StatusCode, fmt.Errorf("not expected status code: %d", resp.StatusCode)
	}

	return resp.StatusCode, nil
}

func uploadOrderToSystem(order string) (int, error) {
	newName := fmt.Sprintf("%d", rand.Intn(1000000000))

	client := http.Client{}

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

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/orders", accrualURL), bytes.NewBuffer(jBody))
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
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Logger.Errorw("error closing response body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusAccepted {
		logger.Logger.Errorw("not expected status code")

		return resp.StatusCode, fmt.Errorf("not expected status code uploading order to system: %d", resp.StatusCode)
	}

	return resp.StatusCode, nil
}
