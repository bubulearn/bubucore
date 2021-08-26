package users

import (
	"bytes"
	"github.com/bubulearn/bubucore"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const logTag = "[bubucore][users]"

const endpointUserInfo = "auth/user/"

// NewClient creates new Client instance
func NewClient(host string, token string) *Client {
	return &Client{
		host:  host,
		token: token,
	}
}

// Client is a bubulearn users service client
type Client struct {
	host    string
	token   string
	_client *http.Client
}

// GetUserInfo fetches user info by user ID
func (c *Client) GetUserInfo(userID string) (user *User, err error) {
	err = c.doRequest(http.MethodGet, endpointUserInfo+userID, nil, &user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Close the Client
func (c *Client) Close() error {
	if c._client != nil {
		c._client.CloseIdleConnections()
	}
	return nil
}

// doRequest sends request to the users service and decodes response to the respData
func (c *Client) doRequest(method string, endpoint string, reqData interface{}, respData interface{}) (err error) {
	err = c.checkPreconditions()
	if err != nil {
		return err
	}

	endpoint = "/" + strings.TrimLeft(endpoint, "/")

	body := []byte("")
	if reqData != nil {
		body, err = jsoniter.Marshal(reqData)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, c.host+endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-type", "application/json")

	resp, err := c.client().Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return bubucore.NewError(resp.StatusCode, logTag+" non-OK response: "+string(body))
	}

	err = jsoniter.Unmarshal(body, respData)
	if err != nil {
		return bubucore.NewError(http.StatusBadGateway, logTag+" failed to decode response: "+string(body))
	}

	return nil
}

// checkPreconditions validates if Client data is ok
func (c *Client) checkPreconditions() error {
	if c.host == "" {
		return bubucore.NewError(http.StatusInternalServerError, logTag+" no host defined")
	}
	if c.token == "" {
		return bubucore.NewError(http.StatusInternalServerError, logTag+" no token defined")
	}
	return nil
}

// client returns http.Client instance
func (c *Client) client() *http.Client {
	if c._client == nil {
		c._client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}
	return c._client
}
