// Package b9s is a Backendless API client
package b9s

import (
	"github.com/bubulearn/bubucore"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
	host = "https://api.backendless.com"
)

// NewClient creates a Client instance
func NewClient(opt *ClientOpt) *Client {
	return &Client{
		opt: opt,
	}
}

// ClientOpt is a B9s Client options
type ClientOpt struct {
	// ProjectID is a B9s project UUID
	ProjectID string

	// APIKey is a B9s REST API key
	APIKey string
}

// Client is a B9s client
type Client struct {
	opt     *ClientOpt
	_client *http.Client
}

// NewRequest creates new b9s GET request
func (c *Client) NewRequest(table string) *Request {
	return &Request{
		Table:     table,
		ProjectID: c.opt.ProjectID,
		APIKey:    c.opt.APIKey,
	}
}

// Do executes the Request and puts parsed response JSON to the target
func (c *Client) Do(req *Request, target interface{}) error {
	r, err := http.NewRequest(http.MethodGet, req.String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.client().Do(r)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		bb, _ := ioutil.ReadAll(resp.Body)
		return bubucore.NewError(
			http.StatusBadGateway,
			"b9s: got code "+strconv.FormatInt(int64(resp.StatusCode), 10)+"; resp: "+string(bb),
		)
	}

	err = jsoniter.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return err
	}

	return nil
}

// client returns http.Client
func (c *Client) client() *http.Client {
	if c._client == nil {
		tr := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   60 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   15 * time.Second,
			ExpectContinueTimeout: 5 * time.Second,
		}
		c._client = &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		}
	}
	return c._client
}
