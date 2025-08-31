// /qqmusic/base_api.go
package qqmusic

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"
	Cookie    = "NMTID=" + GetGuid()
)

// BaseNativeApi 模拟 C# 的基类
type BaseNativeApi struct {
	CookieFunc func() string
	Client     *http.Client
}

// NewBaseNativeApi 创建一个 BaseNativeApi 实例
func NewBaseNativeApi(cookieFunc func() string) *BaseNativeApi {
	return &BaseNativeApi{
		CookieFunc: cookieFunc,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// HttpRefer 定义请求的 Referer
func (b *BaseNativeApi) HttpRefer() string {
	// 默认值，可以在子类中“重写”
	return "https://c.y.qq.com/"
}

// sendPost 发送 POST 请求 (application/x-www-form-urlencoded)
func (b *BaseNativeApi) sendPost(targetURL string, data map[string]string) (string, error) {
	formData := url.Values{}
	for key, val := range data {
		formData.Set(key, val)
	}

	req, err := http.NewRequest("POST", targetURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", b.HttpRefer())
	req.Header.Set("User-Agent", UserAgent)
	cookie := ""
	if b.CookieFunc != nil {
		cookie = b.CookieFunc()
	}
	if cookie == "" {
		cookie = Cookie
	}
	req.Header.Set("Cookie", cookie)

	resp, err := b.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// sendJsonPost 发送 POST 请求 (application/json)
func (b *BaseNativeApi) sendJsonPost(targetURL string, data map[string]interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", targetURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", b.HttpRefer())
	req.Header.Set("User-Agent", UserAgent)
	cookie := ""
	if b.CookieFunc != nil {
		cookie = b.CookieFunc()
	}
	if cookie == "" {
		cookie = Cookie
	}
	req.Header.Set("Cookie", cookie)

	resp, err := b.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func GetGuid() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var sb strings.Builder
	for i := 0; i < 10; i++ {
		sb.WriteString(strconv.Itoa(r.Intn(10)))
	}
	return sb.String()
}
