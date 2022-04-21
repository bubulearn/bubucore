package staticservice

import (
	"bytes"
	"github.com/bubulearn/bubucore"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const logTag = "[bubucore][staticservice]"

const (
	endpointUpload  = "/uploader/upload"
	endpointUploads = "/uploader/uploads"
)

const (
	inputKeyFile  = "file"
	inputKeyTTL   = "ttl"
	inputKeyTitle = "title"
)

// NewClient creates new Client instance
func NewClient(host string, sign string) *Client {
	return &Client{
		host: host,
		sign: sign,
	}
}

// Client is a bubulearn static service client
type Client struct {
	host string
	sign string

	_client *http.Client
}

// GetAll returns all uploads
func (c *Client) GetAll() (uploads []*Upload, err error) {
	err = c.DoJSONRequest(http.MethodGet, endpointUploads, nil, &uploads)
	if err != nil {
		return nil, err
	}

	return uploads, nil
}

// GetUploadInfo fetches upload info by upload ID
func (c *Client) GetUploadInfo(uploadID string) (upload *Upload, err error) {
	err = c.DoJSONRequest(http.MethodGet, endpointUpload+"/"+uploadID, nil, &upload)
	if err != nil {
		return nil, err
	}

	return upload, nil
}

// Upload sends file upload request
func (c *Client) Upload(title string, data []byte, ttl uint64) (upload *Upload, err error) {
	err = c.checkPreconditions()
	if err != nil {
		return nil, err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(inputKeyFile, title)
	if err != nil {
		return nil, err
	}

	_, err = part.Write(data)
	if err != nil {
		return nil, err
	}

	if err := writer.WriteField(inputKeyTitle, title); err != nil {
		return nil, err
	}

	if err := writer.WriteField(inputKeyTTL, strconv.FormatUint(ttl, 10)); err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.host+endpointUpload, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	var upl *Upload

	err = c.doRequest(req, &upl)
	if err != nil {
		return nil, err
	}

	return upl, nil
}

// Close the Client
func (c *Client) Close() error {
	if c._client != nil {
		c._client.CloseIdleConnections()
	}
	return nil
}

// DoJSONRequest sends request to the users service and decodes response to the respData
func (c *Client) DoJSONRequest(method string, endpoint string, reqData interface{}, respData interface{}) (err error) {
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

	req.Header.Set("Content-type", "application/json")

	return c.doRequest(req, respData)
}

func (c *Client) doRequest(req *http.Request, respData interface{}) (err error) {
	req.Header.Set("Authorization", "Bearer "+c.sign)

	resp, err := c.client().Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
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
	if c.sign == "" {
		return bubucore.NewError(http.StatusInternalServerError, logTag+" no sign defined")
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
