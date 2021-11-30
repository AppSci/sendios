package sendios

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

const (
	CategorySystem  = "1"
	CategoryTrigger = "2"
)

type (
	EmailRequest struct {
		TypeID   string                 `json:"type_id,omitempty"`
		Category string                 `json:"category,omitempty"`
		Data     map[string]interface{} `json:"data"`
		Meta     map[string]interface{} `json:"meta,omitempty"`
	}
)

type EmailResponse struct {
	Meta `json:"_meta"`
	Data struct {
		Queued bool `json:"queued"`
		RTime  int  `json:"r_time"`
	} `json:"data"`
}

func (c *Client) SendEmail(r EmailRequest) (*EmailResponse, error) {
	type req struct {
		EmailRequest
		ClientID  string `json:"client_id"`
		ProjectID int    `json:"project_id"`
	}
	data, err := json.Marshal(req{
		EmailRequest: r,
		ClientID:     strconv.Itoa(c.Config.ClientID),
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
	if err := mapResponse(body, &resp); err != nil {
		fmt.Println(string(body))

		return nil, errors.Wrap(err, "map last payment response")
	}

	return &resp, nil
}
