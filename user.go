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
			ID int64 `json:"id"`
			//ProjectID    int         `json:"project_id"`
			//ProjectTitle string      `json:"project_title"`
			//Email        string      `json:"email"`
			//Name         string      `json:"name"`
			//Gender       interface{} `json:"gender"`
			//Country      string      `json:"country"`
			//Language     string      `json:"language"`
			//ErrResponse  int         `json:"err_response"`
			//LastOnline   interface{} `json:"last_online"`
			//LastReaction interface{} `json:"last_reaction"`
			//LastMailed   string      `json:"last_mailed"`
			//LastRequest  interface{} `json:"last_request"`
			//Activation   interface{} `json:"activation"`
			//Meta         struct {
			//	Profile struct {
			//		Age       int         `json:"age"`
			//		Ak        interface{} `json:"ak"`
			//		Photo     interface{} `json:"photo"`
			//		PartnerID interface{} `json:"partner_id"`
			//	} `json:"profile"`
			//} `json:"meta"`
			//Clicks    int    `json:"clicks"`
			//Sends     int    `json:"sends"`
			//CreatedAt string `json:"created_at"`
			//SentMails []struct {
			//	ID          int64       `json:"id"`
			//	CategoryID  int         `json:"category_id"`
			//	TypeID      int         `json:"type_id"`
			//	SubjectID   int         `json:"subject_id"`
			//	TemplateID  int         `json:"template_id"`
			//	SplitGroup  int         `json:"split_group"`
			//	SourceID    int         `json:"source_id"`
			//	ServerID    int         `json:"server_id"`
			//	CreatedAt   string      `json:"created_at"`
			//	MailGroupID int         `json:"mail_group_id"`
			//	PreHeaderID interface{} `json:"pre_header_id"`
			//	CategorySig string      `json:"category_sig"`
			//	ServerName  string      `json:"server_name"`
			//	ServerIP    string      `json:"server_ip"`
			//	SourceName  string      `json:"source_name"`
			//	TypeSig     string      `json:"type_sig"`
			//} `json:"sent_mails"`
			//Unsubscribe      []interface{} `json:"unsubscribe"`
			//UnsubscribeTypes []interface{} `json:"unsubscribe_types"`
			//UnsubPromo       []interface{} `json:"unsub_promo"`
			//WebPush          []interface{} `json:"webpush"`
			//LastPayment      struct {
			//	ID           int         `json:"id"`
			//	UserID       int         `json:"user_id"`
			//	ProjectID    int         `json:"project_id"`
			//	StartedAt    int         `json:"started_at"`
			//	PaymentCount int         `json:"payment_count"`
			//	ExpiresAt    int         `json:"expires_at"`
			//	Active       int         `json:"active"`
			//	PaymentType  interface{} `json:"payment_type"`
			//	Amount       interface{} `json:"amount"`
			//} `json:"last_payment"`
			//ChannelID    interface{} `json:"channel_id"`
			//SubChannelID interface{} `json:"subchannel_id"`
		} `json:"user"`
	} `json:"data"`
}

type (
	CreateUserRequest struct {
		Name          string `json:"name,omitempty"`
		Gender        string `json:"gender,omitempty"`
		Age           int    `json:"age,omitempty"`
		Photo         string `json:"photo,omitempty"`
		AK            string `json:"ak,omitempty"`
		VIP           *int   `json:"vip,omitempty"`
		Language      string `json:"language,omitempty"`
		Country       string `json:"country,omitempty"`
		PlatformID    *int   `json:"platform_id,omitempty"`
		ListID        *int   `json:"list_id,omitempty"`
		Status        *int   `json:"status,omitempty"`
		PartnerID     *int   `json:"partner_id,omitempty"`
		Field1        *int   `json:"field1,omitempty"`
		SessionsCount *int   `json:"sessions_count,omitempty"`
		SessionLast   *int   `json:"session_last,omitempty"`
	}
	CreateUserResponse struct {
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

func (c *Client) SaveUser(email string, req CreateUserRequest) (*CreateUserResponse, error) {
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

	var resp CreateUserResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Println(string(body))

		return nil, errors.Wrap(err, "map save user response")
	}

	return &resp, nil
}
