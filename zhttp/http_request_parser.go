package zhttp

import (
	"errors"
	"regexp"
	"strings"
)

/**
 * @Desc:
 * @Date: 2021/3/6 11:24
 */

/**
 * 解析HTTP请求的第一行，包括method、path、query、protocol version
 */
func parseFirstLine(request string) (*HttpRequest, int, error) {
	if request == "" {
		return nil, -1, errors.New("http request is empty")
	}
	firstLineRegex := regexp.MustCompile(`^(\w+) ([^\s]+) HTTP/([\d.]+)[\r\n]`)
	m := firstLineRegex.FindStringSubmatch(request)
	if m == nil {
		return nil, -1, errors.New("invalid http request(check first line)")
	}
	splitFlagIndex := strings.Index(m[2], "?")
	var path, query string
	if splitFlagIndex >= 0 {
		path = m[2][:splitFlagIndex]
		query = m[2][splitFlagIndex+1:]
	} else {
		path = m[2]
		query = ""
	}
	httpRequest := &HttpRequest{
		Method:   m[1],
		Path:     path,
		Query:    query,
		Protocol: m[3],
	}
	return httpRequest, len(m[0]), nil
}

/*
 * 解析一个Header，返回header name, header value
 */
func parseOneHeader(headerLine string) (string, string) {
	index := strings.Index(headerLine, ":")
	if index < 0 {
		return "", ""
	}
	var name, value string
	name = headerLine[:index]
	if index+1 < len(headerLine) && headerLine[index+1] == ' ' {
		value = headerLine[index+2:]
	} else {
		value = headerLine[index+1:]
	}
	return name, value
}

/**
 * 解析headers
 * 返回 headers、host、body index
 * http中换行符是\r\n，尽管如此，因为不同OS在文本复制时可能会出现\r\n转换为\n的情况，因此兼容两种形式处理
 *
 */
func parseHeaders(request string, index int) ([][2]string, string, int) {
	headers := [][2]string{}
	start := index
	var host, headerLine string
	for i := index; i < len(request); i++ {
		if request[i] == '\n' {
			if i-1 > 0 && request[i-1] == '\r' {
				headerLine = request[start : i-1]
			} else {
				headerLine = request[start:i]
			}
			start = i + 1
			if headerLine == "" {
				break
			}
			headerName, headerValue := parseOneHeader(headerLine)
			// Host行不计入headers
			if "Host" == headerName {
				host = headerValue
			} else {
				headers = append(headers, [2]string{headerName, headerValue})
			}

		}
	}
	return headers, host, start
}

func ParseHttpRequest(request string, scheme string) (*HttpRequest, error) {
	httpRequest, index, err := parseFirstLine(request)
	if err != nil {
		return nil, err
	}
	httpRequest.Headers, httpRequest.Host, index = parseHeaders(request, index)
	httpRequest.Body = request[index:]
	httpRequest.Scheme = scheme
	return httpRequest, nil
}

/*
 * 重放请求
 * 需要指定原请求、scheme、requestOptions
 * 其中requestOptions 的优先级较高，会覆盖原请求中的值
 */
func RepeatRequest(request string, scheme string, requestOptions *RequestOptions) (*Response, error) {
	scheme = strings.ToLower(scheme)
	if "http" != scheme && "https" != scheme {
		return nil, errors.New("invalid scheme: " + scheme)
	}
	httpRequest, err := ParseHttpRequest(request, scheme)
	if err != nil {
		return nil, err
	}
	if requestOptions == nil {
		requestOptions = &RequestOptions{
			RawData: httpRequest.Body,
			Headers: httpRequest.GetHeadersMap(),
		}
	} else {
		if requestOptions.RawData == "" {
			requestOptions.RawData = httpRequest.Body
		}
		if requestOptions.Headers == nil {
			requestOptions.Headers = httpRequest.GetHeadersMap()
		} else {
			tmpHeaderMap := httpRequest.GetHeadersMap()
			// 这里体现了requestOptions中的优先级更高
			for k, v := range requestOptions.Headers {
				tmpHeaderMap[k] = v
			}
			requestOptions.Headers = tmpHeaderMap
		}
	}
	return Request(httpRequest.Method, httpRequest.GetUrl(), requestOptions)
}
