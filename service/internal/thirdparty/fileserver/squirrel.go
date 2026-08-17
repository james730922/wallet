package fileserver

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/conf"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

var (
	client *SquirrelClient
)

const (
	loginURI  string = "/v1/user/login"
	fileURI   string = "/v1/fs"
	healthURI string = "/health"

	fileServerLoginTimeout     = 5 * time.Second
	fileServerUploadTimeout    = 30 * time.Second
	fileServerOperationTimeout = 10 * time.Second
	fileServerHealthTimeout    = 3 * time.Second
)

type loginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type response struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type FileInfo struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	Type    string    `json:"type"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mtime"`
}

type SquirrelClient struct {
	mx         sync.Mutex
	Host       string
	user       string
	passwd     string
	token      string
	httpClient *http.Client
}

func New(configs ...Config) (FileServer, error) {
	config := Config{
		Host:     conf.FileServer().GetHost(),
		Username: conf.FileServer().GetUser(),
		Password: conf.FileServer().GetPassword(),
	}
	if len(configs) > 0 {
		config = configs[0]
	}

	client = &SquirrelClient{
		mx:         sync.Mutex{},
		Host:       config.Host,
		user:       config.Username,
		passwd:     config.Password,
		httpClient: newHTTPClient(),
	}
	err := client.Login()
	if err != nil {
		logger.ApLog().Error(err)
		// even if fileServer is dead, start service as well
		// return nil, err
	}

	return client, nil
}

// Login ...
func (c *SquirrelClient) Login() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), fileServerLoginTimeout)
	defer cancel()

	loginReq := loginRequest{
		Name:     c.user,
		Password: c.passwd,
	}

	loginReqByte, err := json.Marshal(loginReq)
	if err != nil {
		logger.ApLog().Error(err)
		return
	}

	url := c.Host + loginURI
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(loginReqByte))
	if err != nil {
		logger.ApLog().Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// 防止多人同時登入
	c.mx.Lock()
	defer c.mx.Unlock()
	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.ApLog().Error(err)
		return
	}

	var body []byte

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.ApLog().Error(err)
		return
	}

	data := &token{}

	err = json.Unmarshal(body, data)
	if err != nil {
		logger.ApLog().Error(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		logger.ApLog().Error(resp.StatusCode, data)
		err = errs.FileServerResponseNotOK
		return
	}

	c.token = data.AccessToken

	return
}

// Upload ...
func (c *SquirrelClient) Upload(path string, fileName string, f io.Reader) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), fileServerUploadTimeout)
	defer cancel()

	url := c.Host + fileURI + path

	fields := map[string]string{
		"filename": fileName,
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("file", fields["filename"])
	if err != nil {
		logger.ApLog().Error(err)
		return
	}

	_, err = io.Copy(fw, f)
	if err != nil {
		logger.ApLog().Error(err)
		return
	}

	for k, v := range fields {
		_ = writer.WriteField(k, v)
	}

	err = writer.Close()
	if err != nil {
		logger.ApLog().Error(err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		logger.ApLog().Error(err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.ApLog().Error(err)
		return
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.ApLog().Error(err)
		return
	}

	data := &response{}

	err = json.Unmarshal(respBody, data)
	if err != nil {
		logger.ApLog().Error(err)
		return
	}

	// 如果resp的狀態碼不為200，紀錄錯誤代碼與訊息，並回傳自定義的錯誤
	if resp.StatusCode != http.StatusOK {
		logger.ApLog().Error(resp.StatusCode, data)
		err = errs.FileServerResponseNotOK
		return
	}

	return
}

// Delete ...
func (c *SquirrelClient) Delete(path string) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), fileServerOperationTimeout)
	defer cancel()

	url := c.Host + fileURI + path
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		logger.ApLog().Error(err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.ApLog().Error(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.ApLog().Error(err)
		return
	}

	data := &response{}
	err = json.Unmarshal(body, data)
	if err != nil {
		logger.ApLog().Error(err)
		return
	}

	// 如果resp的狀態碼不為200，紀錄錯誤代碼與訊息，並回傳自定義的錯誤
	if resp.StatusCode != http.StatusOK {
		logger.ApLog().Error(resp.StatusCode, data)
		err = errs.FileServerResponseNotOK
		return
	}

	ok = true
	return
}

// GetFileList ...
func (c *SquirrelClient) GetFileList(path string) (file FileInfo, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), fileServerOperationTimeout)
	defer cancel()

	url := c.Host + fileURI + path + "?op=info"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		logger.ApLog().Error(err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.ApLog().Error(err)
		return
	}

	var body []byte

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.ApLog().Error(err)
		return
	}

	data := &response{}
	err = json.Unmarshal(body, data)
	if err != nil {
		logger.ApLog().Error(err)
		return
	}

	// 如果resp的狀態碼不為200，紀錄錯誤代碼與訊息，並回傳自定義的錯誤
	if resp.StatusCode != http.StatusOK {
		logger.ApLog().Error(resp.StatusCode, data)
		err = errs.FileServerResponseNotOK
		return
	}

	err = json.Unmarshal(data.Data, &file)
	if err != nil {
		logger.ApLog().Error(err)
		return
	}

	return
}

// HealthCheck ...
func (c *SquirrelClient) HealthCheck() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), fileServerHealthTimeout)
	defer cancel()

	url := c.Host + healthURI
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		logger.ApLog().Error(err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.ApLog().Error(err)
		return
	}
	defer resp.Body.Close()

	// 如果resp的狀態碼不為200，紀錄，並回傳自定義的錯誤
	if resp.StatusCode != http.StatusOK {
		logger.ApLog().Error(resp.StatusCode)
		err = errs.FileServerResponseNotOK
		return
	}

	return
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: fileServerUploadTimeout,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: time.Second,
			IdleConnTimeout:       90 * time.Second,
		},
	}
}
