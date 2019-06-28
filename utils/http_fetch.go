package utils

import (
	"bytes"
	"cms2cs1806-go/util"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
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
	startTime := time.Now()
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
	endTime := time.Now()
	util.AppendFile("/tmp/http_time.log", fmt.Sprintf("%s\t%s\t%s\t%d\r\n", url, startTime.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"), endTime.UnixNano()/1e6-startTime.UnixNano()/1e6))
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
