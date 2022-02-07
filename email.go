package sendios

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

const (
	CategorySystem  = 1
	CategoryTrigger = 2
)

type (
	EmailRequest struct {
		TypeID   int                    `json:"type_id,omitempty"`
		Category int                    `json:"category,omitempty"`
		Data     map[string]interface{} `json:"data"`
		Meta     map[string]interface{} `json:"meta,omitempty"`
	}
)

type EmailResponse struct {
	Meta `json:"_meta"`
	Data struct {
		Error  string `json:"error"`
		Queued bool   `json:"queued"`
		RTime  int    `json:"r_time"`
	} `json:"data"`
}

// https://sendios.readme.io/reference/send-system-email

func (c *Client) SendEmail(r EmailRequest) (*EmailResponse, error) {
	type req struct {
		EmailRequest
		ClientID  int `json:"client_id"`
		ProjectID int `json:"project_id"`
	}
	data, err := json.Marshal(req{
		EmailRequest: r,
		ClientID:     c.Config.ClientID,
		ProjectID:    c.Config.Project,
	})
	if err != nil {
		return nil, errors.Wrap(err, "marshal payload")
	}

	_, body, err := c.makeRequest(http.MethodPost, "https://api.sendios.io/v1/push/system", bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "send email")
	}

	var resp EmailResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Println(string(body))

		return nil, errors.Wrap(err, "map eemail response")
	}

	if resp.Data.Error != "" {
		return &resp, errors.New(resp.Data.Error)
	}

	return &resp, nil
}
