package zhttp

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/mounsurf/lib/util"
	"io"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

var quoteEscaper = strings.NewReplacer(`\`, `\\`, `"`, `\"`)

type FileUpload struct {
	// Filename is the name of the file that you wish to upload. We use this to guess the mimetype as well as pass it onto the server
	FileName string

	// FileContents is happy as long as you pass it a io.ReadCloser (which most file use anyways)
	FileContents io.ReadCloser

	// FieldName is form field name
	FieldName string

	// FileMime represents which mimetime should be sent along with the file.
	// When empty, defaults to application/octet-stream
	FileMime string
}

type RequestOptions struct {
	//0、Path字段
	RawPath string

	//一、Query字段
	// 原始Query，不会被编码处理。该字段会覆盖原URL的query字段。
	RawQuery string

	// 请求Query字段。会自动编码成string。会和原有的URL query合并。POST、GET均在URL QUERY中
	Params map[string]string

	//二、Header字段
	// 自定义HOST。修改HOST仅能通过该字段实现，不能通过设置Headers实现。
	Host string

	// 自定义请求Header。优先级较高。会覆盖通过Cookies、RawCookie、ContentType、UserAgent设置的头。
	Headers map[string]string

	// 会自动将map转为string。Cookies会和RawCookie合并。但是均会被headers中设置的Cookie覆盖。
	Cookies map[string]string
	// 原始Cookie，不会被编码处理
	RawCookie string

	ContentType string
	UserAgent   string

	// 是否为ajax请求。true会设置Header {X-Requested-With: XMLHttpRequest}。
	IsAjax bool

	// basic认证 []string{username, password}
	Auth []string

	// 请求禁用gzip压缩
	DisableCompression bool

	//三、Data字段(只有一个生效，优先级自上而下)
	// 原始Data，不会被编码处理
	RawData string

	// JSON数据，不会被编码。会自动设置Content-Type: application/json
	JSON string

	// XML数据，不会被编码。 会自动设置Content-Type: application/xml
	XML string

	// 文件上传
	Files []FileUpload

	// 请求Body字段。会自动编码成string。POST、GET均在BODY中。
	Data map[string]string

	// 四、HTTP处理字段
	// 是否校验证书。默认校验。注意：Go TLS机制无法验证证书是否被撤销。
	IgnoreCertError bool

	// 自定义HOSTS（DNS解析）
	Hosts string

	// 设置代理：格式为[协议]URL
	// *protocol* => proxy address e.g http => http://127.0.0.1:8080
	Proxies map[string]*url.URL

	// 拨号超时时间。默认5秒。
	DialTimeout time.Duration

	// 整个请求（包括拨号/请求/重定向）超时时间。默认10秒
	RequestTimeout time.Duration

	// 禁止跟随跳转
	DisableRedirect bool
}

func (ro *RequestOptions) closeFiles() {
	for _, f := range ro.Files {
		_ = f.FileContents.Close()
	}
}

func (ro RequestOptions) proxySettings(req *http.Request) (*url.URL, error) {
	if len(ro.Proxies) > 0 {
		// There was a proxy specified – do we support the protocol?
		if p, ok := ro.Proxies[req.URL.Scheme]; ok {
			return p, nil
		}
	}

	// Proxies were specified but not for any protocol that we use
	return http.ProxyFromEnvironment(req)
}

func FileUploadFromDisk(fieldName, filePath string) ([]FileUpload, error) {
	filePath = filepath.ToSlash(filepath.Clean(filePath))
	fd, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	_, fileName := filepath.Split(filePath)
	return []FileUpload{
		FileUpload{
			FieldName:    fieldName,
			FileName:     fileName,
			FileContents: fd,
		},
	}, nil
}

func createTransport(ro RequestOptions) *http.Transport {
	transport := &http.Transport{
		MaxIdleConns:          100,
		IdleConnTimeout:       5 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		Proxy:                 ro.proxySettings,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: ro.IgnoreCertError},
		DisableCompression:    ro.DisableCompression,
		DisableKeepAlives:     true,
	}

	transport.Dial = func(network, address string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, address, ro.DialTimeout)
		if err != nil {
			return nil, err
		}
		return newTimeoutConn(conn, ro.RequestTimeout), nil
	}

	return transport
}

func buildClient(ro RequestOptions, cookieJar http.CookieJar) *http.Client {

	// The function does not return an error ever... so we are just ignoring it
	if cookieJar == nil {
		cookieJar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	}
	client := &http.Client{
		Jar:       cookieJar,
		Transport: createTransport(ro),
		Timeout:   ro.RequestTimeout,
	}

	if ro.DisableRedirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return client
}

func buildRequest(method, urlStr string, ro *RequestOptions) (*http.Request, error) {
	if ro.RawData != "" {
		return http.NewRequest(method, urlStr, strings.NewReader(ro.RawData))
	}

	if ro.JSON != "" {
		return createBasicJSONRequest(method, urlStr, ro)
	}

	if ro.XML != "" {
		return createBasicXMLRequest(method, urlStr, ro)
	}

	if ro.Files != nil {
		return createFileUploadRequest(method, urlStr, ro)
	}

	if ro.Data != nil {
		return createBasicRequest(method, urlStr, ro)
	}
	return http.NewRequest(method, urlStr, nil)
}

func createBasicJSONRequest(method, urlStr string, ro *RequestOptions) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, strings.NewReader(ro.JSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func createBasicXMLRequest(method, urlStr string, ro *RequestOptions) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, strings.NewReader(ro.XML))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/xml")

	return req, nil
}

func createFileUploadRequest(method, urlStr string, ro *RequestOptions) (*http.Request, error) {
	if method == "POST" {
		return createMultiPartPostRequest(method, urlStr, ro)
	}

	// This may be a PUT or PATCH request so we will just put the raw
	// io.ReadCloser in the request body
	// and guess the MIME type from the file name

	// At the moment, we will only support 1 file upload as a time
	// when uploading using PUT or PATCH

	req, err := http.NewRequest(method, urlStr, ro.Files[0].FileContents)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", mime.TypeByExtension(filepath.Ext(ro.Files[0].FileName)))

	return req, nil
}

func createMultiPartPostRequest(method, urlStr string, ro *RequestOptions) (*http.Request, error) {
	body := &bytes.Buffer{}

	multipartWriter := multipart.NewWriter(body)

	for i, f := range ro.Files {
		if f.FileContents == nil {
			return nil, errors.New("Pointer FileContents cannot be nil")
		}

		fieldName := f.FieldName

		if fieldName == "" {
			if len(ro.Files) > 1 {
				fieldName = "file" + strconv.Itoa(i+1)
			} else {
				fieldName = "file"
			}
		}

		var writer io.Writer
		var err error

		if f.FileMime != "" {
			h := make(textproto.MIMEHeader)
			h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes(fieldName), escapeQuotes(f.FileName)))
			h.Set("Content-Type", f.FileMime)
			writer, err = multipartWriter.CreatePart(h)
		} else {
			writer, err = multipartWriter.CreateFormFile(fieldName, f.FileName)
		}

		if err != nil {
			return nil, err
		}

		if _, err = io.Copy(writer, f.FileContents); err != nil && err != io.EOF {
			return nil, err
		}
	}

	// Populate the other parts of the form (if there are any)
	for key, value := range ro.Data {
		multipartWriter.WriteField(key, value)
	}

	if err := multipartWriter.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, urlStr, body)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", multipartWriter.FormDataContentType())

	return req, err
}

func createBasicRequest(method, urlStr string, ro *RequestOptions) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, strings.NewReader(encodePostValues(ro.Data)))

	if err != nil {
		return nil, err
	}

	// The content type must be set to a regular form
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func doRequest(method, urlStr string, ro *RequestOptions, cookieJar http.CookieJar) (*Response, error) {
	if ro == nil {
		ro = &RequestOptions{}
	}
	defer ro.closeFiles()
	if ro.DialTimeout == 0 {
		ro.DialTimeout = 10 * time.Second
	}
	if ro.RequestTimeout == 0 {
		ro.RequestTimeout = 20 * time.Second
	}
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	// 强制不连接某些站点
	if err = checkAvailable(parsedURL.Hostname()); err != nil {
		return nil, err
	}
	// 自动忽略IP的证书错误
	// 不过引入一个问题，如果设置了autoIgnoreCertErr，那么调用方将无法再明确的校验IP的证书
	// TODO: 鉴于没有想到具体场景，先不解决该问题
	if autoIgnoreCertErr && util.IpRegex.MatchString(parsedURL.Hostname()) {
		ro.IgnoreCertError = true
	}

	urlStr, err = buildURL(parsedURL, ro)
	if err != nil {
		return nil, err
	}
	req, err := buildRequest(method, urlStr, ro)
	if err != nil {
		return nil, err
	}
	if ro.RawPath != "" {
		//req.URL.Path = ro.RawPath 有bug，改为如下代码
		req.URL.RawPath = ro.RawPath
	}
	parseHosts(req, ro)
	addCookies(req, ro)
	addHeaders(req, ro)

	client := buildClient(*ro, cookieJar)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	response := &Response{RawResponse: resp}
	err = response.setBodyAndClose()
	if err != nil {
		return nil, err
	}
	if esFlag {
		logToEs(response)
	}
	return response, nil
}

// buildURLParams returns a URL with all of the params
// Note: This function will override current URL params if they contradict what is provided in the map
// That is what the "magic" is on the last line
func buildURL(parsedURL *url.URL, ro *RequestOptions) (string, error) {
	if ro.RawQuery != "" {
		parsedURL.RawQuery = ro.RawQuery
	} else {
		if parsedURL.RawQuery != "" {
			parsedURL.RawQuery = escapeRawQuery(parsedURL.RawQuery)
		}
		if len(ro.Params) > 0 {
			query := url.Values{}
			for key, value := range ro.Params {
				query.Set(key, value)
			}
			if parsedURL.RawQuery != "" {
				parsedURL.RawQuery += "&" + query.Encode()
			} else {
				parsedURL.RawQuery = query.Encode()
			}
		}

	}
	return parsedURL.String(), nil
}

func parseHosts(req *http.Request, ro *RequestOptions) {
	if ro.Hosts != "" {
		req.Host = req.URL.Host
		port := req.URL.Port()
		if port != "" {
			req.URL.Host = ro.Hosts + ":" + port
		} else {
			req.URL.Host = ro.Hosts
		}
	}
}

// addHTTPHeaders adds any additional HTTP headers that need to be added are added here including:
// 1. Authorization Headers
// 2. Any other header requested
func addHeaders(req *http.Request, ro *RequestOptions) {
	for key, value := range ro.Headers {
		req.Header.Set(key, value)
	}

	if ro.Host != "" {
		req.Host = ro.Host
	}

	if ro.Auth != nil {
		req.SetBasicAuth(ro.Auth[0], ro.Auth[1])
	}

	if ro.IsAjax {
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
	}

	if ro.ContentType != "" {
		req.Header.Set("Content-Type", ro.ContentType)
	}

	if ro.UserAgent != "" {
		req.Header.Set("User-Agent", ro.UserAgent)
	} else if ro.Headers["User-Agent"] != "" {
		req.Header.Set("User-Agent", ro.Headers["User-Agent"])
	} else {
		req.Header.Set("User-Agent", randUserAgent())
	}

	if ro.Headers["Accept"] != "" {
		req.Header.Set("Accept", ro.Headers["Accept"])
	} else {
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	}
}

func addCookies(req *http.Request, ro *RequestOptions) {
	if ro.RawCookie != "" {
		req.Header.Set("Cookie", ro.RawCookie)
	}
	for k, v := range ro.Cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}
}

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func encodePostValues(postValues map[string]string) string {
	urlValues := &url.Values{}

	for key, value := range postValues {
		urlValues.Set(key, value)
	}

	return urlValues.Encode() // This will sort all of the string values
}

func escapeRawQuery(s string) string {
	var r string
	for i := 0; i < len(s); {
		if shouldEscape(s[i]) {
			r += "%" + string("0123456789ABCDEF"[s[i]>>4]) + string("0123456789ABCDEF"[s[i]&0x0F])
		} else {
			r += string(int(s[i]))
		}
		i++
	}

	return r
}
func shouldEscape(c byte) bool {
	// 按照chrome的规则，!"$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_`abcdefghijklmnopqrstuvwxyz{|}~中除了' " < > ` 之外全部不编码
	if ('?' <= c && c <= '_') || ('a' <= c && c <= '~') || c == '=' || ('(' <= c && c <= ';') || ('#' <= c && c <= '&') || c == '!' {
		return false
	}
	return true
}
