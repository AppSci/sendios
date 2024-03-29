package sendios

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

const (
	PlatformUnknown = iota
	PlatformDesktop
	PlatformMobile
	PlatformAndroid
	PlatformIos
)

const (
	VipNo = iota
	VipYes
)

const (
	GenderMale   = "m"
	GenderFemale = "f"
)

type UserResponse struct {
	Meta `json:"_meta"`
	Data struct {
		User struct {
			ID           int64       `json:"id"`
			ProjectID    int         `json:"project_id"`
			ProjectTitle string      `json:"project_title"`
			Email        string      `json:"email"`
			Name         string      `json:"name"`
			Gender       string      `json:"gender"`
			Country      interface{} `json:"country"`
			Language     interface{} `json:"language"`
			ErrResponse  int         `json:"err_response"`
			Confirm      int         `json:"confirm"`
			LastOnline   interface{} `json:"last_online"`
			LastReaction interface{} `json:"last_reaction"`
			LastMailed   string      `json:"last_mailed"`
			LastRequest  interface{} `json:"last_request"`
			Activation   interface{} `json:"activation"`
			Meta         interface{} `json:"meta"`
			Clicks       int         `json:"clicks"`
			Sends        int         `json:"sends"`
			CreatedAt    string      `json:"created_at"`
			SentMails    []struct {
				ID          int64       `json:"id"`
				CategoryID  int         `json:"category_id"`
				TypeID      int         `json:"type_id"`
				SubjectID   int         `json:"subject_id"`
				TemplateID  int         `json:"template_id"`
				SplitGroup  int         `json:"split_group"`
				SourceID    interface{} `json:"source_id"`
				ServerID    int         `json:"server_id"`
				CreatedAt   string      `json:"created_at"`
				MailGroupID int         `json:"mail_group_id"`
				PreHeaderID interface{} `json:"pre_header_id"`
				CategorySig string      `json:"category_sig"`
				ServerName  string      `json:"server_name"`
				ServerIp    string      `json:"server_ip"`
				SourceName  string      `json:"source_name"`
				TypeSig     string      `json:"type_sig"`
			} `json:"sent_mails"`
			Unsubscribe      interface{}   `json:"unsubscribe"`
			UnsubscribeTypes []interface{} `json:"unsubscribe_types"`
			UnsubPromo       interface{}   `json:"unsub_promo"`
			Webpush          interface{}   `json:"webpush"`
			LastPayment      interface{}   `json:"last_payment"`
			ChannelID        interface{}   `json:"channel_id"`
			SubchannelID     interface{}   `json:"subchannel_id"`
			CustomFields     interface{}   `json:"custom_fields"`
			Vip              int           `json:"vip"`
		} `json:"user"`
	} `json:"data"`
}

type UserResponseByID struct {
	Meta `json:"_meta"`
	Data struct {
		Result struct {
			User struct {
				// ProjectID    int         `json:"project_id"`
				// ListID       interface{} `json:"list_id"`
				// Email        string      `json:"email"`
				// Name         interface{} `json:"name"` // name can be int
				// LanguageID   interface{} `json:"language_id"`
				// CityID       int         `json:"city_id"`
				// VendorID     int         `json:"vendor_id"`
				// ValidID      interface{} `json:"valid_id"`
				// CountryID    interface{} `json:"country_id"`
				// PlatformID   int         `json:"platform_id"`
				// Gender       string      `json:"gender"`
				Confirm int `json:"confirm"`
				// Vip          int         `json:"vip"`
				// ErrResponse  int         `json:"err_response"`
				// LastOnline   int         `json:"last_online"`
				// LastReaction int         `json:"last_reaction"`
				// LastMailed   int         `json:"last_mailed"`
				// CreatedAt    string      `json:"created_at"`
				// ID           int         `json:"id"`
				// RegisteredAt string      `json:"registered_at"`
				// UpdatedAt    string      `json:"updated_at"`
			} `json:"user"`
			CustomFields interface{} `json:"custom_fields"`
		} `json:"result"`
	} `json:"data"`
}

type (
	UpdateUserFieldsRequest  map[string]interface{}
	UpdateUserFieldsResponse struct {
		Meta `json:"_meta"`
		Data struct {
			Result bool `json:"result"`
		} `json:"data"`
	}
)

func (c *Client) GetUserInfo(email string) (*UserResponse, error) {
	url := fmt.Sprintf("https://api.sendios.io/v1/user/project/%d/email/%s", c.Config.Project, email)

	statusCode, body, err := c.makeRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "get user info")
	}

	if statusCode != http.StatusOK {
		var resp ErrorResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			fmt.Println(string(body))

			return nil, errors.Wrap(err, "map get user error")
		}

		return nil, fmt.Errorf("get user error: %s", resp.Data.Error)
	}

	var resp UserResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Println(string(body))

		return nil, errors.Wrap(err, "map get user response")
	}

	return &resp, nil
}

func (c *Client) GetUserInfoByID(id int64) (*UserResponseByID, error) {
	// https://sendios.readme.io/reference/get-user-custom-fields-by-user
	url := fmt.Sprintf("https://api.sendios.io/v1/userfields/user/%d", id)

	statusCode, body, err := c.makeRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "get user info")
	}

	if statusCode != http.StatusOK {
		var resp ErrorResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			fmt.Println(string(body))

			return nil, errors.Wrap(err, "map get user error")
		}

		return nil, fmt.Errorf("get user error: %s", resp.Data.Error)
	}

	var resp UserResponseByID
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Println(string(body))

		return nil, errors.Wrap(err, "map get user response")
	}

	return &resp, nil
}

func (c *Client) GetUserFields(email string) (*UserResponseByID, error) {
	// https://sendios.readme.io/reference/get-user-custom-fields
	url := fmt.Sprintf("https://api.sendios.io/v1/userfields/project/%d/email/%s", c.Config.Project, email)

	statusCode, body, err := c.makeRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "get user info")
	}

	if statusCode != http.StatusOK {
		var resp ErrorResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			fmt.Println(string(body))

			return nil, errors.Wrap(err, "map get user error")
		}

		return nil, fmt.Errorf("get user error: %s", resp.Data.Error)
	}

	var resp UserResponseByID
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Println(string(body))

		return nil, errors.Wrap(err, "map get user response")
	}

	return &resp, nil
}

func (c *Client) SetUserData(email string, req UpdateUserFieldsRequest) (*UpdateUserFieldsResponse, error) {
	emailHash := base64.StdEncoding.EncodeToString([]byte(email))

	data, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "marshal payload")
	}

	url := fmt.Sprintf("https://api.sendios.io/v1/userfields/project/%d/emailhash/%s", c.Config.Project, emailHash)

	statusCode, body, err := c.makeRequest(http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "save user")
	}

	if statusCode != http.StatusOK {
		var resp ErrorResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			fmt.Println(string(body))

			return nil, errors.Wrap(err, "map save user error")
		}

		return nil, fmt.Errorf("save user error: %s", resp.Data.Error)
	}

	var resp UpdateUserFieldsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Println(string(body))

		return nil, errors.Wrap(err, "map save user response")
	}

	return &resp, nil
}

type UnsubscribeUserResponse struct {
	Meta struct {
		Status string `json:"status"`
		Time   int64  `json:"time"`
		Count  int    `json:"count"`
	} `json:"_meta"`
	Data struct {
		Unsub bool `json:"unsub"`
	} `json:"data"`
}

func (c *Client) UnsubscribeUser(email string) (*UnsubscribeUserResponse, error) {
	emailHash := base64.StdEncoding.EncodeToString([]byte(email))

	url := fmt.Sprintf("https://api.sendios.io/v1/unsub/admin/%d/email/%s", c.Config.Project, emailHash)

	statusCode, body, err := c.makeRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unsubscribe user")
	}

	if statusCode != http.StatusOK {
		var resp ErrorResponse
		if err = json.Unmarshal(body, &resp); err != nil {
			fmt.Println(string(body))

			return nil, errors.Wrap(err, "map unsubscribe user error")
		}

		return nil, fmt.Errorf("unsubscribe user error: %s", resp.Data.Error)
	}

	var resp UnsubscribeUserResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		fmt.Println(string(body))

		return nil, errors.Wrap(err, "map unsubscribe user response")
	}

	return &resp, nil
}

type CreateClientUserResponse struct {
	Meta struct {
		Status string `json:"status"`
		Time   int64  `json:"time"`
		Count  int    `json:"count"`
	} `json:"_meta"`
	Data struct {
		Message string `json:"message"`
		Date    string `json:"date"`
		Status  bool   `json:"status"`
	} `json:"data"`
}

func (c *Client) CreateClientUser(email string) (*CreateClientUserResponse, error) {
	type CreateClientUserRequest struct {
		Email        string `json:"email"`
		ProjectID    int    `json:"project_id"`
		ClientUserID int    `json:"client_user_id"`
	}

	data, err := json.Marshal(CreateClientUserRequest{
		Email:        email,
		ProjectID:    c.Config.Project,
		ClientUserID: c.Config.ClientID,
	})

	statusCode, body, err := c.makeRequest(http.MethodPost, "https://api.sendios.io/v1/clientuser/create", bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "create user")
	}

	if statusCode != http.StatusOK {
		var resp ErrorResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			fmt.Println(string(body))

			return nil, errors.Wrap(err, "map create user error")
		}

		return nil, fmt.Errorf("create user error: %s", resp.Data.Error)
	}

	var resp CreateClientUserResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		fmt.Println(string(body))

		return nil, errors.Wrap(err, "map create user response")
	}

	return &resp, nil
}

type UnsubscribeStatusResponse struct {
	Meta struct {
		Status string `json:"status"`
		Time   int    `json:"time"`
		Count  int    `json:"count"`
	} `json:"_meta"`
	Data struct {
		Result struct {
			ProjectID int         `json:"project_id"`
			UserID    int         `json:"user_id"`
			SourceID  interface{} `json:"source_id"`
			Meta      struct {
				Message string `json:"message"`
			} `json:"meta"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
		} `json:"result"`
	} `json:"data"`
}

type SubscribedStatusResponse struct {
	Meta struct {
		Status string `json:"status"`
		Time   int    `json:"time"`
		Count  int    `json:"count"`
	} `json:"_meta"`
	Data struct {
		Result bool `json:"result"`
	} `json:"data"`
}

func (c *Client) CheckIsUnsubscribedUser(userID int64) (bool, error) {
	statusCode, body, err := c.makeRequest(http.MethodGet, fmt.Sprintf("https://api.sendios.io/v1/unsub/isunsub/%d", userID), nil)
	if err != nil {
		return false, errors.Wrap(err, "check unsubscribed")
	}

	if statusCode != http.StatusOK {
		var resp ErrorResponse
		if err = json.Unmarshal(body, &resp); err != nil {
			fmt.Println(string(body))

			return false, errors.Wrap(err, "map check unsubscribed error")
		}

		return false, fmt.Errorf("check unsubscribed error: %s", resp.Data.Error)
	}

	var resp UnsubscribeStatusResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		var r SubscribedStatusResponse
		if err = json.Unmarshal(body, &r); err != nil {
			return false, errors.Wrap(err, "map check unsubscribed response")
		}

		return false, nil
	}

	return true, nil
}

func (c *Client) CheckIsUnsubscribed(email string) (bool, error) {
	userInfo, err := c.GetUserInfo(email)
	if err != nil {
		return false, errors.Wrap(err, "get user info")
	}

	return c.CheckIsUnsubscribedUser(userInfo.Data.User.ID)
}

type (
	ProvideClientRequest struct {
		Email        string `json:"email"`
		ProjectID    int    `json:"project_id"`
		ClientUserID string `json:"client_user_id"`
	}
	ProvideClientResponse struct {
		Meta `json:"_meta"`
		Data struct {
			Message string `json:"message"`
			Date    string `json:"date"`
			Status  bool   `json:"status"`
		} `json:"data"`
	}
)

func (c *Client) ProvideClientID(req ProvideClientRequest) (*ProvideClientResponse, error) {
	// https://sendios.readme.io/reference/providing-user-id-on-product

	data, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "marshal request")
	}

	statusCode, body, err := c.makeRequest(http.MethodPost, "https://api.sendios.io/v1/clientuser/create", bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "unsubscribe user")
	}

	if statusCode != http.StatusOK {
		var resp ErrorResponse
		if err = json.Unmarshal(body, &resp); err != nil {
			fmt.Println(string(body))

			return nil, errors.Wrap(err, "map unsubscribe user error")
		}

		return nil, fmt.Errorf("unsubscribe user error: %s", resp.Data.Error)
	}

	var resp ProvideClientResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		fmt.Println(string(body))

		return nil, errors.Wrap(err, "map unsubscribe user response")
	}

	return &resp, nil
}
