package biligo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
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

// request v为携带的参数，用于debug输出
func (h *baseClient) request(req *http.Request, v interface{}) ([]byte, error) {
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
		h.logger.Printf("%s %s  %+v", req.Method, req.URL, v)
		h.logger.Printf("%s", string(raw))
	}

	return raw, nil
}
func (h *baseClient) raw(base, endpoint, method string, payload map[string]string, dAfter func(d *url.Values), reqAfter func(r *http.Request)) ([]byte, error) {
	var (
		req *http.Request
		err error
	)

	data := url.Values{}
	for k, v := range payload {
		data.Add(k, v)
	}

	// 侵入处理values
	if dAfter != nil {
		dAfter(&data)
	}

	link := base + endpoint
	switch method {
	case http.MethodGet:
		if req, err = http.NewRequest(method, link, nil); err != nil {
			return nil, err
		}
		req.URL.RawQuery = data.Encode()
	case http.MethodPost:
		if req, err = http.NewRequest(method, link, strings.NewReader(data.Encode())); err != nil {
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

	return h.request(req, payload)
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
func (h *baseClient) upload(base, endpoint string, payload map[string]string, files []*FileUpload, mAfter func(m *multipart.Writer) error, reqAfter func(r *http.Request)) ([]byte, error) {
	var (
		req *http.Request
		err error
	)
	link := base + endpoint

	body := new(bytes.Buffer)
	mp := multipart.NewWriter(body)

	for _, f := range files {
		var ff io.Writer
		if ff, err = mp.CreateFormFile(f.field, f.name); err != nil {
			return nil, err
		}
		if _, err = io.Copy(ff, f.file); err != nil {
			return nil, err
		}
	}

	for k, v := range payload {
		if err = mp.WriteField(k, v); err != nil {
			return nil, err
		}
	}

	// 侵入处理field
	if mAfter != nil {
		if err = mAfter(mp); err != nil {
			return nil, err
		}
	}

	// 为mp添加结束符
	if err = mp.Close(); err != nil {
		return nil, err
	}

	// 只支持POST
	if req, err = http.NewRequest(http.MethodPost, link, body); err != nil {
		return nil, err
	}

	req.Header.Add("Origin", "https://www.bilibili.com")
	req.Header.Add("Referer", "https://www.bilibili.com")
	req.Header.Add("Content-type", mp.FormDataContentType())
	req.Header.Add("User-Agent", h.ua)

	// 侵入处理req
	if reqAfter != nil {
		reqAfter(req)
	}

	// 文件不输出，否则全是乱码
	return h.request(req, payload)
}
