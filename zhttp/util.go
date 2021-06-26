package zhttp

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/url"
)

func GetTitle(body string) string {
	title := titleRegex.FindStringSubmatch(body)
	if title != nil && len(title) >= 1 {
		return title[1]
	}
	return ""
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

/**
 * 指定端口，生成本地代理配置
 */
func GetLocalProxy(port int) map[string]*url.URL {
	proxy, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", port))
	if err != nil {
		return nil
	}
	return map[string]*url.URL{
		"http":  proxy,
		"https": proxy,
	}
}
