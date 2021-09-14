package users

import (
	"bytes"
	"context"
	"github.com/bubulearn/bubucore"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const logTag = "[bubucore][users]"

const (
	endpointUserInfo    = "auth/user/"
	endpointUsersGetAll = "auth/users/all"
)

// NewClient creates new Client instance
func NewClient(host string, token string) *Client {
	return &Client{
		host:  host,
		token: token,
	}
}

// Client is a bubulearn users service client
type Client struct {
	host  string
	token string

	redis    *redis.Client
	cacheTTL int

	_client *http.Client
}

// SetRedis sets redis client to cache results with
func (c *Client) SetRedis(client *redis.Client, ttl int) {
	c.redis = client
	c.cacheTTL = ttl
}

// GetAll returns all users
func (c *Client) GetAll() (users []*User, err error) {
	err = c.DoRequest(http.MethodGet, endpointUsersGetAll, nil, &users)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		c.writeToCache(user)
	}

	return users, nil
}

// GetUserInfo fetches user info by user ID
func (c *Client) GetUserInfo(userID string) (user *User, err error) {
	user = c.readFromCache(userID)
	if user != nil {
		return user, nil
	}

	err = c.DoRequest(http.MethodGet, endpointUserInfo+userID, nil, &user)
	if err != nil {
		return nil, err
	}

	c.writeToCache(user)

	return user, nil
}

// Close the Client
func (c *Client) Close() error {
	if c._client != nil {
		c._client.CloseIdleConnections()
	}
	return nil
}

// DoRequest sends request to the users service and decodes response to the respData
func (c *Client) DoRequest(method string, endpoint string, reqData interface{}, respData interface{}) (err error) {
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

// cacheKey returns cache key for the user info
func (c *Client) cacheKey(userID string) string {
	return "bubuusersservice:userinfo:" + userID
}

// readFromCache reads user info from the cache
func (c *Client) readFromCache(userID string) *User {
	if c.redis == nil {
		return nil
	}

	cacheKey := c.cacheKey(userID)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cached, err := c.redis.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil
	}

	var user *User
	err = jsoniter.Unmarshal([]byte(cached), &user)
	if err != nil {
		c.redis.Del(ctx, cacheKey)
		return nil
	}

	return user
}

// writeToCache saves user info to the cache
func (c *Client) writeToCache(user *User) {
	if c.redis == nil {
		return
	}

	cacheKey := c.cacheKey(user.ID)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	data, err := jsoniter.Marshal(user)
	if err != nil {
		return
	}

	c.redis.Set(ctx, cacheKey, data, time.Second*time.Duration(c.cacheTTL))
}
