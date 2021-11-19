package sendios

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"os"
)

const (
	CategorySystem  = "1"
	CategoryTrigger = "2"
)

type (
	EmailRequest struct {
		TypeID   string                 `json:"type_id,omitempty"`
		Category string                 `json:"category,omitempty"`
		Data     Data                   `json:"data"`
		Meta     map[string]interface{} `json:"meta,omitempty"`
	}
	Data struct {
		User User `json:"user,omitempty"`
	}
	User struct {
		Email      string `json:"email,omitempty"`
		Name       string `json:"name,omitempty"`
		Age        string `json:"age,omitempty"`
		Gender     string `json:"gender,omitempty"`
		Language   string `json:"language,omitempty"`
		Country    string `json:"country,omitempty"`
		PlatformID string `json:"platform_id,omitempty"`
		Vip        string `json:"vip,omitempty"`
		Photo      string `json:"photo,omitempty"`
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
		ClientID:     os.Getenv("SENDIOS_CLIENT_ID"),
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
