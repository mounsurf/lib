package zhttp

import (
	"bytes"
	"compress/gzip"
	"errors"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Response struct {
	RawResponse *http.Response
	body        []byte
}

func (r *Response) StatusCode() int {
	return r.RawResponse.StatusCode
}

func (r *Response) RawHeaders() string {
	var rawHeader string
	for k, v := range r.RawResponse.Header {
		rawHeader += k + ": " + strings.Join(v, ",") + "\r\n"
	}
	return strings.TrimSuffix(rawHeader, "\r\n")
}

func (r *Response) HeadersMap() map[string]string {
	headers := map[string]string{}
	for k, v := range r.RawResponse.Header {
		headers[k] = strings.Join(v, ",")
	}
	return headers
}

func (r *Response) HasHeader(headerName string) bool {
	for k, _ := range r.RawResponse.Header {
		if k == headerName {
			return true
		}
	}
	return false
}

func (r *Response) HasHeaderAndValue(headerName, headerVaule string) bool {
	for name, values := range r.RawResponse.Header {
		if name == headerName {
			for _, value := range values {
				if value == headerVaule {
					return true
				}
			}
		}
	}
	return false
}

func (r *Response) GetHeader(headerName string) string {
	for k, v := range r.RawResponse.Header {
		if k == headerName {
			return strings.Join(v, ",")
		}
	}
	return ""
}

func (r *Response) RawCookies() string {
	var rawCookie string
	for _, cookie := range r.RawResponse.Cookies() {
		rawCookie += cookie.Name + "=" + cookie.Value + ";"
	}
	return strings.TrimSuffix(rawCookie, ";")
}

func (r *Response) CookiesMap() map[string]string {
	cookies := map[string]string{}
	for _, cookie := range r.RawResponse.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}
	return cookies
}

func (r *Response) HasCookieAndValue(cookieName, cookieVaule string) bool {
	for _, cookie := range r.RawResponse.Cookies() {
		if cookie.Name == cookieName && cookie.Value == cookieVaule {
			return true
		}
	}
	return false
}

func (r *Response) HasCookie(cookieName string) bool {
	for _, cookie := range r.RawResponse.Cookies() {
		if cookie.Name == cookieName {
			return true
		}
	}
	return false
}

func (r *Response) GetBodyString() string {
	if r == nil {
		return ""
	}
	return string(r.GetBodyBytes())
}

func (r *Response) GetBodyBytes() []byte {
	if r == nil {
		return []byte{}
	}
	return r.body
}

func (r *Response) GetTitle() string {
	if r == nil {
		return ""
	}
	titleResult := titleRegex.FindSubmatch(r.GetBodyBytes())
	if titleResult != nil && len(titleResult) >= 1 {
		return html.UnescapeString(strings.Replace(strings.Replace(string(titleResult[1]), "\n", "", -1), "\r", "", -1))
	}
	return ""
}

func (r *Response) ReadN(n int64) []byte {
	if n > int64(len(r.body)) {
		n = int64(len(r.body))
	}
	return r.body[:n-1]
}

func (r *Response) RawRequest() string {
	rawRequest := r.RawResponse.Request.Method + " " + r.RawResponse.Request.URL.RequestURI() + " " + r.RawResponse.Request.Proto + "\r\n"
	host := r.RawResponse.Request.Host
	if host == "" {
		host = r.RawResponse.Request.URL.Host
	}
	rawRequest += "Host: " + host + "\r\n"
	for key, val := range r.RawResponse.Request.Header {
		rawRequest += key + ": " + val[0] + "\r\n"
	}
	rawRequest += "\r\n" + r.reqBody()

	return rawRequest
}

func (r *Response) reqBody() string {
	var body string
	if r.RawResponse.Request.GetBody != nil {
		b, err := r.RawResponse.Request.GetBody()
		if err == nil {
			buf, err := ioutil.ReadAll(b)
			b.Close()
			if err == nil {
				body = string(buf)
			}
		}
	}
	return body
}

func (r *Response) setBodyAndClose() error {
	err := r.setBody()
	if err != nil {
		_ = r.RawResponse.Body.Close()
		return err
	}
	return r.RawResponse.Body.Close()
}

// 这个似乎没有意义
func (r *Response) updateContentLength() {
	if r.RawResponse.Header != nil {
		contentLengthHeader := "Content-Length"
		contentLengthValues := r.RawResponse.Header[contentLengthHeader]
		if contentLengthValues == nil {
			contentLengthValues = []string{}
		}
		contentLengthValues = append(contentLengthValues, strconv.Itoa(len(r.body)))
		r.RawResponse.Header[contentLengthHeader] = contentLengthValues
	}
}

func (r *Response) setBody() error {
	var err error
	// 此处如果返回的body是gzip格式，那么需要解码
	if r.GetHeader("Content-Encoding") == "gzip" {
		delete(r.RawResponse.Header, "Content-Encoding")
		gzipReader, err := gzip.NewReader(r.RawResponse.Body)
		if err != nil {
			return err
		}

		defer gzipReader.Close()
		r.body, err = ioutil.ReadAll(gzipReader)
		if err != nil {
			return err
		}
		r.updateContentLength()

	} else {
		r.body, err = ioutil.ReadAll(r.RawResponse.Body)
		if err != nil {
			return err
		}
	}

	//  直接解码，（很多情况下不会用到解码的数据，所以会牺牲效率）
	r.decodeBody()
	return nil
}

//解码
func (r *Response) decodeBody() {
	if r.body != nil {
		//只看前多少字符有没有编码
		maxLength := 5000
		length := len(r.body)
		if length > maxLength {
			length = maxLength
		}
		e, _, certain := charset.DetermineEncoding(r.body[0:length], "")
		var reader *transform.Reader
		if certain {
			reader = transform.NewReader(bytes.NewReader(r.body), e.NewDecoder())
		} else {
			reader = transform.NewReader(bytes.NewReader(r.body), encoding.Nop.NewDecoder())
		}
		utf8Body, err := ioutil.ReadAll(reader)
		if err == nil {
			r.body = utf8Body
			r.updateContentLength()
		}
	}
}

// 获取Location的绝对URL
func (r *Response) GetAbsoluteLocation() (string, error) {
	location := r.GetHeader("Location")
	if location == "" {
		return "", errors.New("location does not exist")
	} else {
		if strings.HasPrefix(location, "http://") || strings.HasPrefix(location, "https://") {
			return location, nil
		} else if strings.HasPrefix(location, "//") {
			return r.RawResponse.Request.URL.Scheme + ":" + location, nil
		} else if strings.HasPrefix(location, "/") {
			return r.RawResponse.Request.URL.Scheme + "://" + r.RawResponse.Request.URL.Host + location, nil
		} else if location[0] == '?' ||
			(location[0] >= '0' && location[0] <= '9') ||
			(location[0] >= 'A' && location[0] <= 'Z') ||
			(location[0] >= 'a' && location[0] <= 'z') {
			return r.RawResponse.Request.URL.Scheme + "://" + r.RawResponse.Request.URL.Host + "/" + location, nil
		} else if location[0] == '.' && location[1] == '/' {
			iniPath := r.RawResponse.Request.URL.Path
			lastIndex := strings.LastIndex(iniPath, "/")
			var iniDir string
			if lastIndex >= 0 {
				iniDir = iniPath[:lastIndex+1]
			} else {
				iniDir = "/"
			}
			return r.RawResponse.Request.URL.Scheme + "://" + r.RawResponse.Request.URL.Host + iniDir + location[2:], nil
		} else {
			return location, nil
		}
	}
}

//获取最终请求的URL，存在URL跳转的情况下，获取最终跳转的地址
func (r *Response) GetFinalUrl() string {
	if r == nil {
		return ""
	}
	return r.RawResponse.Request.URL.String()
}
