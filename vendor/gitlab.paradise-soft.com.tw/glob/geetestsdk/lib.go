package geetestsdk

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

/**
 * sdk lib包，核心逻辑。
 *
 * @author liuquan@geetest.com
 */
const (
	API_URL              string = "http://api.geetest.com"
	REGISTER_URL         string = "/register.php"
	VALIDATE_URL         string = "/validate.php"
	JSON_FORMAT          string = "1"
	NEW_CAPTCHA          bool   = true
	HTTP_TIMEOUT_DEFAULT int    = 5 // 单位：秒
	VERSION              string = "golang-gin:3.1.1"
	GEETEST_CHALLENGE    string = "geetest_challenge" // 极验二次验证表单传参字段 chllenge
	GEETEST_VALIDATE     string = "geetest_validate"  // 极验二次验证表单传参字段 validate
	GEETEST_SECCODE      string = "geetest_seccode"   // 极验二次验证表单传参字段 seccode
)

type lib struct {
	geeTestID  string // 公钥
	geeTestKey string // 私钥
	libResult  *LibResult
}

//func NewGeetestLib(geetest_id string, geetest_key string) *GeetestLib {
//	return &GeetestLib{geetest_id, geetest_key, NewGeetestLibResult()}
//}

//--------------------------------------
//	將 NewGeetestLib 替換成 Pool 形式
//	避免驗證場景過量造成GC負擔
//--------------------------------------

type libPool struct {
	pool sync.Pool
}

func newLibPool() *libPool {
	return &libPool{
		pool: sync.Pool{New: func() interface{} {
			return &lib{
				geeTestID:  "",
				geeTestKey: "",
				libResult:  &LibResult{},
			}
		}},
	}
}

func (gp *libPool) Get(geeTestID, geeTestKey string) *lib {
	l := gp.pool.Get().(*lib)
	l.geeTestID = geeTestID
	l.geeTestKey = geeTestKey
	return l
}

func (gp *libPool) Put(lib *lib) {
	lib.Reset()
	gp.pool.Put(lib)
}

func (g *lib) Reset() {
	g.geeTestID = ""
	g.geeTestKey = ""
	g.libResult = NewGeetestLibResult()
}

/**
 * 验证初始化
 */
func (g *lib) Register(digestmod string, params map[string]string) *LibResult {
	origin_challenge := g.requestRegister(params)
	g.buildRegisterResult(origin_challenge, digestmod)
	return g.libResult
}

func (g *lib) LocalRegister() *LibResult {
	g.buildRegisterResult("", "")
	return g.libResult
}

/**
 * 向极验发送验证初始化的请求，GET方式
 */
func (g *lib) requestRegister(params map[string]string) string {
	params["gt"] = g.geeTestID
	params["json_format"] = JSON_FORMAT
	params["sdk"] = VERSION
	register_url := API_URL + REGISTER_URL
	resBody, err := g.httpGet(register_url, params)
	if err != nil {
		return ""
	}
	resMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(resBody), &resMap)
	if err != nil {
		return ""
	}
	return resMap["challenge"].(string)
}

/**
 * 构建验证初始化返回数据
 */
func (g *lib) buildRegisterResult(origin_challenge string, digestmod string) {
	// origin_challenge为空或者值为0代表失败
	if origin_challenge == "" || origin_challenge == "0" {
		// 本地随机生成32位字符串
		characters := []byte("0123456789abcdefghijklmnopqrstuvwxyz")
		challenge := []byte{}
		for i := 0; i < 32; i++ {
			challenge = append(challenge, characters[rand.Intn(len(characters))])
		}
		dataMap := map[string]interface{}{
			"success":     false,
			"gt":          g.geeTestID,
			"challenge":   string(challenge),
			"new_captcha": NEW_CAPTCHA,
		}
		dataJson, _ := json.Marshal(dataMap)
		g.libResult.setAll(0, string(dataJson), "获取当前bypass状态为fail，后续流程走宕机模式")
	} else {
		challenge := ""
		if digestmod == "md5" {
			challenge = g.md5_encode(origin_challenge + g.geeTestKey)
		} else if digestmod == "sha256" {
			challenge = g.sha256_encode(origin_challenge + g.geeTestKey)
		} else if digestmod == "hmac-sha256" {
			challenge = g.hmac_sha256_encode(origin_challenge, g.geeTestKey)
		} else {
			challenge = g.md5_encode(origin_challenge + g.geeTestKey)
		}
		dataMap := map[string]interface{}{
			"success":     true,
			"gt":          g.geeTestID,
			"challenge":   challenge,
			"new_captcha": NEW_CAPTCHA,
		}
		dataJson, _ := json.Marshal(dataMap)
		g.libResult.setAll(1, string(dataJson), "")
	}
}

/**
 * 正常流程下（即验证初始化成功），二次验证
 */
func (g *lib) SuccessValidate(challenge string, validate string, seccode string) *LibResult {
	if !g.checkParam(challenge, validate, seccode) {
		g.libResult.setAll(0, "", "正常模式，本地校验，参数challenge、validate、seccode不可为空")
	} else {
		response_seccode := g.requestValidate(challenge, validate, seccode)
		if response_seccode == "" {
			g.libResult.setAll(0, "", "请求极验validate接口失败")
		} else if response_seccode == "false" {
			g.libResult.setAll(0, "", "极验二次验证不通过")
		} else {
			g.libResult.setAll(1, "", "")
		}
	}
	return g.libResult
}

/**
 * 异常流程下（即验证初始化失败，宕机模式），二次验证
 * 注意：由于是宕机模式，初衷是保证验证业务不会中断正常业务，所以此处只作简单的参数校验，可自行设计逻辑。
 */
func (g *lib) FailValidate(challenge string, validate string, seccode string) *LibResult {
	if !g.checkParam(challenge, validate, seccode) {
		g.libResult.setAll(0, "", "宕机模式，本地校验，参数challenge、validate、seccode不可为空.")
	} else {
		g.libResult.setAll(1, "", "")
	}
	return g.libResult
}

/**
 * 向极验发送二次验证的请求，POST方式
 */
func (g *lib) requestValidate(challenge string, validate string, seccode string) string {
	params := map[string]string{}
	params["seccode"] = seccode
	params["json_format"] = JSON_FORMAT
	params["challenge"] = challenge
	params["sdk"] = VERSION
	params["captchaid"] = g.geeTestID
	validate_url := API_URL + VALIDATE_URL
	resBody, err := g.httpPost(validate_url, params)
	if err != nil {
		return ""
	}
	resMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(resBody), &resMap)
	if err != nil {
		return ""
	}
	return resMap["seccode"].(string)
}

/**
 * 校验二次验证的三个参数，校验通过返回true，校验失败返回false
 */
func (g *lib) checkParam(challenge string, validate string, seccode string) bool {
	return !(challenge == "" || strings.TrimSpace(challenge) == "" || validate == "" || strings.TrimSpace(validate) == "" || seccode == "" || strings.TrimSpace(seccode) == "")
}

/**
 * 发送GET请求，获取服务器返回结果
 */
func (g *lib) httpGet(getUrl string, params map[string]string) (string, error) {
	q := url.Values{}
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
	}
	req, err := http.NewRequest(http.MethodGet, getUrl, nil)
	if err != nil {
		return "", errors.New("NewRequest fail")
	}
	req.URL.RawQuery = q.Encode()
	client := &http.Client{Timeout: time.Duration(HTTP_TIMEOUT_DEFAULT) * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode == 200 {
		return string(body), nil
	}
	return "", nil
}

/**
 * 发送POST请求，获取服务器返回结果
 */
func (g *lib) httpPost(postUrl string, params map[string]string) (string, error) {
	q := url.Values{}
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
	}
	req, err := http.NewRequest(http.MethodPost, postUrl, strings.NewReader(q.Encode()))
	if err != nil {
		return "", errors.New("NewRequest fail")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{Timeout: time.Duration(HTTP_TIMEOUT_DEFAULT) * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode == 200 {
		return string(body), nil
	}
	return "", nil
}

/**
 * md5 加密
 */
func (g *lib) md5_encode(value string) string {
	h := md5.New()
	h.Write([]byte(value))
	return fmt.Sprintf("%x", h.Sum(nil))
}

/**
 * sha256加密
 */
func (g *lib) sha256_encode(value string) string {
	h := sha256.New()
	h.Write([]byte(value))
	return fmt.Sprintf("%x", h.Sum(nil))
}

/**
 * hmac-sha256 加密
 */
func (g *lib) hmac_sha256_encode(value string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(value))
	return fmt.Sprintf("%x", h.Sum(nil))
}
