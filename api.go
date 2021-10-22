package biligo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type baseClient struct {
	debug  bool
	client *http.Client
	ua     string
	logger *log.Logger
}
type baseSetting struct {
	// 自定义http client
	//
	// 默认为 http.http.DefaultClient
	Client *http.Client
	// Debug模式 true将输出请求信息 false不输出
	//
	// 默认false
	DebugMode bool
	// 自定义UserAgent
	//
	// 默认Chrome随机Agent
	UserAgent string
	// Logger 的输出前缀，区分Client
	Prefix string
}

func newBaseClient(setting *baseSetting) *baseClient {
	client := setting.Client
	if client == nil {
		client = http.DefaultClient
	}

	ua := setting.UserAgent
	if ua == "" {
		rand.Seed(time.Now().UnixNano())
		ua = userAgent[rand.Intn(len(userAgent))]
	}

	return &baseClient{
		debug:  setting.DebugMode,
		client: client,
		ua:     ua,
		logger: log.New(os.Stdout, setting.Prefix, log.LstdFlags),
	}
}

func (h *baseClient) raw(base, endpoint, method string, payload map[string]string, dAfter func(d *url.Values), reqAfter func(r *http.Request)) ([]byte, error) {
	var (
		req  *http.Request
		err  error
		data url.Values
	)
	link := base + endpoint

	data = url.Values{}
	for k, v := range payload {
		data.Add(k, v)
	}

	// 侵入处理values
	if dAfter != nil {
		dAfter(&data)
	}

	switch method {
	case "GET":
		req, err = http.NewRequest(method, link, nil)
		if err != nil {
			return nil, err
		}
		req.URL.RawQuery = data.Encode()
	case "POST":
		req, err = http.NewRequest(method, link, strings.NewReader(data.Encode()))
		if err != nil {
			return nil, err
		}
	}

	req.Header.Add("Origin", "https://www.bilibili.com")
	req.Header.Add("Referer", "https://www.bilibili.com")
	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", h.ua)

	// 侵入处理req
	if reqAfter != nil {
		reqAfter(req)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	resp.Close = true
	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if h.debug {
		h.logger.Printf("url: %s value: %+v", link, data)
		h.logger.Printf("resp: %+v", string(raw))
	}

	return raw, nil
}
func (h *baseClient) parse(raw []byte) (*Response, error) {
	var result = &Response{}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("(%d) %s", result.Code, result.Message)
	}
	return result, nil
}
