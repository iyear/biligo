package util

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	qrcode "github.com/skip2/go-qrcode"
)

type Cookie struct {
	SESSDATA          string
	Bili_jct          string
	DedeUserID        string
	DedeUserID__ckMd5 string
	Sid               string
}

type GetQrcodeResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		URL       string `json:"url"`
		QrcodeKey string `json:"qrcode_key"`
	} `json:"data"`
}

type ScanResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		URL          string `json:"url"`
		RefreshToken string `json:"refresh_token"`
		Timestamp    int    `json:"timestamp"`
		Code         int    `json:"code"`
		Message      string `json:"message"`
	} `json:"data"`
}

func getQrcode() *GetQrcodeResponse {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: time.Second * 30,
	}

	req, err := http.NewRequest("GET", "https://passport.bilibili.com/x/passport-login/web/qrcode/generate", nil)
	if err != nil {
		panic(err)
	}

	req.Header = http.Header{
		"Referer":    []string{"https://passport.bilibili.com/login?from_spm_id=333.851.top_bar.login_window"},
		"User-Agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.51 Safari/537.36 Edg/93.0.961.27"},
	}

	// Send the GET request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Parse the JSON data
	var getQrcodeResp *GetQrcodeResponse
	if err := json.Unmarshal(body, &getQrcodeResp); err != nil {
		panic(err)
	}

	qr, err := qrcode.New(getQrcodeResp.Data.URL, qrcode.Medium)
	if err != nil {
		fmt.Println("无法生成二维码：", err)
		return nil
	}

	// 打印二维码到控制台
	fmt.Println(qr.ToSmallString(false))

	return getQrcodeResp
}

func getCookie(getQrcodeResp *GetQrcodeResponse) *Cookie {
	// 等待扫码
	fmt.Println("请在8秒内完成扫码...")
	time.Sleep(time.Second * 8)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: time.Second * 30,
	}

	url := "https://passport.bilibili.com/x/passport-login/web/qrcode/poll?qrcode_key=" + getQrcodeResp.Data.QrcodeKey
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header = http.Header{
		"Referer":    []string{"https://passport.bilibili.com/login?from_spm_id=333.851.top_bar.login_window"},
		"User-Agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.51 Safari/537.36 Edg/93.0.961.27"},
	}

	// Send the GET request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 解析 JSON 数据
	var scanResp *ScanResponse
	if err := json.Unmarshal(body, &scanResp); err != nil {
		panic(err)
	}

	if scanResp.Message != "0" {
		fmt.Println("扫码失败")
	}

	cookies := resp.Cookies()
	result := &Cookie{}
	for _, cookie := range cookies {
		switch cookie.Name {
		case "SESSDATA":
			result.SESSDATA = cookie.Value
		case "bili_jct":
			result.Bili_jct = cookie.Value
		case "DedeUserID":
			result.DedeUserID = cookie.Value
		case "DedeUserID__ckMd5":
			result.DedeUserID__ckMd5 = cookie.Value
		case "sid":
			result.Sid = cookie.Value
		}
	}
	return result
}

func GetCookie() *Cookie {
	resp := getQrcode()
	cookies := getCookie(resp)
	return cookies
}
