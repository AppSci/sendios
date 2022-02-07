package sendios

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

type Client struct {
	Config        Config
	HTTPClient    *http.Client
	DebugRequests bool
}

type Config struct {
	Project     int
	ClientID    int
	ClientToken string
}

func New(projectID, clientID int, clientToken string) *Client {
	return &Client{
		Config: Config{
			Project:     projectID,
			ClientID:    clientID,
			ClientToken: clientToken,
		},
		HTTPClient: &http.Client{},
	}
}

func NewFromEnv() (*Client, error) {
	env := os.Getenv("SENDIOS_CONFIG")
	if env == "" {
		return nil, nil
	}

	values, err := url.ParseQuery(env)
	if err != nil {
		return nil, errors.Wrap(err, "parse config")
	}

	c := Client{HTTPClient: &http.Client{}}

	if val, ok := values["client_id"]; ok {
		c.Config.ClientID, err = strconv.Atoi(val[0])
		if err != nil {
			return nil, errors.Wrap(err, "parse client_id from env")
		}
	}

	if val, ok := values["client_token"]; ok {
		c.Config.ClientToken = val[0]
	}

	if val, ok := values["project"]; ok {
		c.Config.Project, err = strconv.Atoi(val[0])
		if err != nil {
			return nil, errors.Wrap(err, "parse project from env")
		}
	}

	if c.Config.ClientID < 0 || c.Config.ClientToken == "" || c.Config.Project < 0 {
		return nil, nil
	}

	return &c, nil
}

func (c *Client) makeRequest(method, url string, reader io.Reader) (int, []byte, error) {
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return 0, nil, err
	}

	req.SetBasicAuth(strconv.Itoa(c.Config.ClientID), c.Config.ClientToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	if c.DebugRequests {
		data, _ := httputil.DumpRequest(req, true)
		fmt.Printf("Sendios request:\n%s\n", data)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, nil, errors.Wrap(err, "send request")
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	if c.DebugRequests {
		data, _ := httputil.DumpResponse(resp, true)
		fmt.Printf("Sendios response:\n%s\n", data)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, errors.Wrap(err, "read body")
	}

	return resp.StatusCode, data, nil
}

type Meta struct {
	Status string `json:"status"`
	Time   int    `json:"time"`
	Count  int    `json:"count"`
}

type ErrorResponse struct {
	Meta `json:"_meta"`
	Data struct {
		Error string `json:"error"`
	} `json:"data"`
}
