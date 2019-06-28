package utils

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type HC struct {
	Client *http.Client
}

var HttpClient *HC

func MakeHttpClient() *HC {
	HttpClient = &HC{}
	HttpClient.Client = &http.Client{}
	return HttpClient
}

func UrlEncode(data map[string]string) string {
	var buf bytes.Buffer
	for k, v := range data {
		buf.WriteString(url.QueryEscape(k))
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(v))
		buf.WriteByte('&')
	}
	s := buf.String()
	return s[0 : len(s)-1]
}

func (httpClient *HC) HttpGet(url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}
	// 添加头部信息
	if len(headers) > 0 {
		for hk, hv := range headers {
			req.Header.Add(hk, hv)
		}
	}
	response, err := httpClient.Client.Do(req) //提交
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func (httpClient *HC) HttpDelete(url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return []byte{}, err
	}
	// 添加头部信息
	if len(headers) > 0 {
		for hk, hv := range headers {
			req.Header.Add(hk, hv)
		}
	}
	response, err := httpClient.Client.Do(req) //提交
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func (httpClient *HC) HttpPost(url string, data interface{}, headers map[string]string) ([]byte, error) {
	dataMap, okMap := data.(map[string]string)
	dataJson, okJson := data.(string)
	var req *http.Request
	var err error
	if okMap {
		req, err = http.NewRequest("POST", url, strings.NewReader(UrlEncode(dataMap)))
	} else if okJson {
		req, err = http.NewRequest("POST", url, strings.NewReader(dataJson))
	} else {
		return []byte{}, errors.New("type error")
	}
	// 添加头部信息
	if len(headers) > 0 {
		for hk, hv := range headers {
			req.Header.Add(hk, hv)
		}
	}
	response, err := httpClient.Client.Do(req) //提交
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}
