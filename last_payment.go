package sendios

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

type (
	LastPaymentRequest struct {
		UserID int64 `json:"user_id"` // User id in sendios
		Payment
	}
	Payment struct {
		StartDate   string `json:"start_date"`             // Payment date or subscription start date
		ExpireDate  string `json:"expire_date,omitempty"`  // Subscription end date (optional, default false)
		TotalCount  string `json:"total_count,omitempty"`  // Pay count (optional, default false)
		PaymentType string `json:"payment_type,omitempty"` // Pay type (optional, default false)
		Amount      string `json:"amount,omitempty"`       // Pay amount (optional, default false)
	}
	LastPaymentResponse struct {
		Meta `json:"_meta"`
		Data struct {
			Message string `json:"message"`
			Date    string `json:"date"`
			Status  bool   `json:"status"`
		} `json:"data"`
	}
)

func (c *Client) SendLastPayment(req LastPaymentRequest) (*LastPaymentResponse, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "marshal payload")
	}

	statusCode, body, err := c.makeRequest(http.MethodPost, "https://api.sendios.io/v1/lastpayment", bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "send last payment")
	}

	if statusCode == http.StatusNotFound {
		var resp ErrorResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			fmt.Println(string(body))

			return nil, errors.Wrap(err, "map last payment error")
		}

		return nil, fmt.Errorf("last payment error: %s", resp.Data.Error)
	}

	var resp LastPaymentResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Println(string(body))

		return nil, errors.Wrap(err, "map last payment response")
	}

	return &resp, nil
}

func (c *Client) SendLastPaymentByEmail(email string, p Payment) (*LastPaymentResponse, error) {
	user, err := c.GetUserInfo(email)
	if err != nil {
		return nil, errors.Wrap(err, "get user")
	}

	payment, err := c.SendLastPayment(LastPaymentRequest{UserID: user.Data.User.ID, Payment: p})
	if err != nil {
		return nil, errors.Wrap(err, "send last payment")
	}

	return payment, nil
}
