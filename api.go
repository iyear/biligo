package biligo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Raw base末尾带/
func (b *BiliClient) Raw(base, endpoint, method string, payload map[string]string) ([]byte, error) {
	var (
		req     *http.Request
		err     error
		reqData url.Values
	)
	link := base + endpoint

	reqData = url.Values{}
	for k, v := range payload {
		reqData.Add(k, v)
	}

	switch method {
	case "GET":
		req, err = http.NewRequest(method, link, nil)
		if err != nil {
			return nil, err
		}
		req.URL.RawQuery = reqData.Encode()
	case "POST":
		reqData.Add("csrf", b.auth.BiliJCT)
		req, err = http.NewRequest(method, link, strings.NewReader(reqData.Encode()))
		if err != nil {
			return nil, err
		}
	}

	req.Header.Add("Cookie",
		fmt.Sprintf("DedeUserID=%s;SESSDATA=%s;DedeUserID__ckMd5=%s",
			b.auth.DedeUserID, b.auth.SESSDATA, b.auth.DedeUserIDCkMd5))
	req.Header.Add("Origin", "https://www.bilibili.com")
	req.Header.Add("Referer", "https://www.bilibili.com")
	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", b.ua)

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	resp.Close = true
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if b.debug {
		b.logger.Printf("url: %s value: %+v", link, reqData)
		b.logger.Printf("resp: %+v", string(data))
	}

	return data, nil
}
func (b *BiliClient) RawParse(base, endpoint, method string, payload map[string]string) (*Response, error) {
	raw, err := b.Raw(base, endpoint, method, payload)
	if err != nil {
		return nil, err
	}
	var result = &Response{}
	if err = json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("(%d) %s", result.Code, result.Message)
	}
	return result, nil
}
func (c *CommClient) Raw(base, endpoint, method string, payload map[string]string) ([]byte, error) {
	var (
		req     *http.Request
		err     error
		reqData url.Values
	)
	link := base + endpoint

	reqData = url.Values{}
	for k, v := range payload {
		reqData.Add(k, v)
	}

	switch method {
	case "GET":
		req, err = http.NewRequest(method, link, nil)
		if err != nil {
			return nil, err
		}
		req.URL.RawQuery = reqData.Encode()
	case "POST":
		req, err = http.NewRequest(method, link, strings.NewReader(reqData.Encode()))
		if err != nil {
			return nil, err
		}
	}

	req.Header.Add("Origin", "https://www.bilibili.com")
	req.Header.Add("Referer", "https://www.bilibili.com")
	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", c.ua)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	resp.Close = true
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if c.debug {
		c.logger.Printf("url: %s value: %+v", link, reqData)
		c.logger.Printf("resp: %+v", string(data))
	}

	return data, nil
}
func (c *CommClient) RawParse(base, endpoint, method string, payload map[string]string) (*Response, error) {
	raw, err := c.Raw(base, endpoint, method, payload)
	if err != nil {
		return nil, err
	}
	var result = &Response{}
	if err = json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("(%d) %s", result.Code, result.Message)
	}
	return result, nil
}
