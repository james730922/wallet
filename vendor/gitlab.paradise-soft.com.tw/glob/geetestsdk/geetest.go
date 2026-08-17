package geetestsdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

type RegisterResp struct {
	Success    bool   `json:"success"`
	GT         string `json:"gt"`
	Challenge  string `json:"challenge"`
	NewCaptcha bool   `json:"new_captcha"`
}

func newHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 5 * time.Second,
	}
}

func NewGeeTestCaptcha(conf CaptchaConfig, opts ...GeeTestCaptchaOpt) ICaptcha {
	obj := &geeTestCaptcha{
		alert:      func() {},
		libPool:    newLibPool(),
		httpClient: newHttpClient(),
		byPassURL:  conf_byPassURL,
		config:     conf,
	}

	for _, opt := range opts {
		opt()
	}

	obj.RunBackground()

	return obj
}

type ICaptcha interface {
	Register(account string) *RegisterResp
	Validate(challenge, validate, secCode string) *LibResult
	GetByPass() bool
	RefreshConfig(conf CaptchaConfig)
}

type alert func()

type geeTestCaptcha struct {
	alert      alert
	libPool    *libPool
	httpClient *http.Client

	byPassURL    string // 向geetest发送获取bypass状态请求的url
	byPassStatus bool   // bypass 狀態

	config CaptchaConfig

	aes cipherAES
}

func (g *geeTestCaptcha) RefreshConfig(conf CaptchaConfig) {
	g.config = conf
}

func (g *geeTestCaptcha) RunBackground() {
	// 檢查極驗服務
	go g.getByPassBackground()
}

// 验证初始化接口，GET请求
func (g *geeTestCaptcha) Register(account string) *RegisterResp {
	// 取得後台 ID, KEY
	gtLib := g.libPool.Get(g.config.GeeTestID, g.config.GeeTestKey)
	defer g.libPool.Put(gtLib)

	// 加密用戶帳號
	userID, err := g.aes.EncryptString(g.config.GeeTestUserIDSalt, account)
	if err != nil {
		// 由於用戶ID只作為極驗統計用途，故只記log不中斷驗證流程
		fmt.Errorf("GEETEST encrypt.AesEncrypt() failed err: %s", err)
		userID = new(string)
	}

	params := map[string]string{
		"digestmod": hash_md5,
		"user_id":   *userID,
		"risk_type": "1",
	}

	var result *LibResult
	// 檢查極驗服務
	if g.byPassStatus {
		// 極驗服務正常
		result = gtLib.Register(hash_md5, params)
	} else {
		// 極驗服務異常
		result = gtLib.LocalRegister()
	}

	return g.arrangeRegisterResp(result)
}

func (g *geeTestCaptcha) arrangeRegisterResp(result *LibResult) *RegisterResp {
	resp := &RegisterResp{}

	if err := json.Unmarshal([]byte(result.Data), resp); err != nil {
		fmt.Errorf("json.Unmarshal error: %s", err)
	}

	return resp
}

// 二次验证接口，POST请求
func (g *geeTestCaptcha) Validate(challenge, validate, secCode string) *LibResult {
	var result *LibResult

	// 取得後台 ID, KEY
	gtLib := g.libPool.Get(g.config.GeeTestID, g.config.GeeTestKey)
	defer g.libPool.Put(gtLib)

	// 檢查極驗服務
	if g.byPassStatus {
		// 極驗服務正常
		result = gtLib.SuccessValidate(challenge, validate, secCode)
	} else {
		// 極驗服務異常
		result = gtLib.FailValidate(challenge, validate, secCode)
		// 不走當機驗證，強制失敗
		result.Status = 0
	}

	g.validateAlert(result)

	// result.Status (int)
	// 狀態 0 = 失敗
	//     1 = 成功
	return g.arrangeValidateResp(result)
}

func (g *geeTestCaptcha) validateAlert(result *LibResult) {
	if result.Status == 0 && result.Msg == "请求极验validate接口失败" {
		g.alert()
	}
}

func (g *geeTestCaptcha) arrangeValidateResp(result *LibResult) *LibResult {
	if result.Status == 0 {
		result.Msg = "验证不通过"
	}

	return result
}

func (g *geeTestCaptcha) getByPassBackground() {
	for {
		ok := g.tryGetPass()
		if !ok {
			//g.alert()
			fmt.Errorf("geetest checkBypassStatus() 当前极验服务状态异常。")
		}

		g.byPassStatus = ok

		time.Sleep(time.Duration(g.config.GeeTestByPassCycleTimeSec) * time.Second)
	}
}

func (g *geeTestCaptcha) tryGetPass() bool {
	// 失敗時再度嘗試，做多 3 次
	for try := 0; try < 3; try++ {
		if g.GetByPass() {
			return true
		}
		time.Sleep(time.Duration(g.config.GeeTestByPassCycleTimeSec) * time.Second)
	}

	return false
}

func (g *geeTestCaptcha) GetByPass() bool {
	params := make(map[string]string)
	params["gt"] = g.config.GeeTestID

	t := time.Now()
	resBody, err := g.getByPass(g.byPassURL, params)
	execTime := time.Now().Sub(t)
	if err != nil {
		fmt.Errorf("geetest getByPass() ,execTime: %s ,err: %s", execTime.String(), err)
		return false
	}

	resMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(resBody), &resMap); err != nil {
		return false
	}

	if resMap["status"] != "success" {
		return false
	}

	return true
}

// 发送GET请求
func (g *geeTestCaptcha) getByPass(getURL string, params map[string]string) (string, error) {
	q := url.Values{}
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
	}

	req, err := http.NewRequest(http.MethodGet, getURL, nil)
	if err != nil {
		return "", errors.New("NewRequest fail")
	}

	req.URL.RawQuery = q.Encode()

	res, err := g.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode == http.StatusOK {
		return string(body), nil
	}

	return "", nil
}
