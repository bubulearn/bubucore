package notifications

import (
	"bytes"
	"github.com/bubulearn/bubucore"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// notification-service endpoints
const (
	EndpointAppReport   = "app-report"
	EndpointCustomEmail = "email"
)

// Templates names
const (
	TplResetPassLink   = "reset-pass-link"
	TplResetPassResult = "reset-pass-result"
)

// NewClient creates new notifications service client
func NewClient(host string, token string) *Client {
	return &Client{
		host:  host,
		token: token,
	}
}

// Client is a notifications service client
type Client struct {
	host  string
	token string

	_client *http.Client
}

// Send sends notification request
func (c *Client) Send(endpoint string, data interface{}) error {
	err := c.checkPreconditions()
	if err != nil {
		return err
	}

	endpoint = "/" + strings.TrimLeft(endpoint, "/")

	body, err := jsoniter.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, c.host+endpoint, bytes.NewBuffer(body))
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

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	respErr := &bubucore.Error{}
	err = jsoniter.Unmarshal(body, &respErr)
	if err != nil {
		return bubucore.NewError(http.StatusBadGateway, "failed to send notification and to decode response: "+string(body))
	}

	return respErr
}

// SendPlainText sends plain message notification
func (c *Client) SendPlainText(endpoint string, msg string) error {
	return c.Send(endpoint, &PlainText{Text: msg})
}

// SendEmail send custom email notification
func (c *Client) SendEmail(n Email) error {
	return c.Send(EndpointCustomEmail, n)
}

// SendAppReport sends notification about app report
func (c *Client) SendAppReport(msg string) error {
	return c.SendPlainText(EndpointAppReport, msg)
}

// Close finalizes the Client
func (c *Client) Close() error {
	if c._client != nil {
		c._client.CloseIdleConnections()
	}
	return nil
}

// checkPreconditions validates if Client data is ok
func (c *Client) checkPreconditions() error {
	if c.host == "" {
		return bubucore.NewError(http.StatusInternalServerError, "no notifications host defined")
	}
	if c.token == "" {
		return bubucore.NewError(http.StatusInternalServerError, "no notifications token defined")
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
