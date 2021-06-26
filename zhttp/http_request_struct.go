package zhttp

import "fmt"

/**
 * @Desc: http 请求结构。
 * @Date: 2021/3/6 16:40
 */

type HttpRequest struct {
	Scheme   string      // 协议： https、http
	Method   string      // 方法
	Path     string      // 路径
	Query    string      // query
	Protocol string      // http协议版本号
	Host     string      // Host
	Headers  [][2]string // 请求头。由于header可能有相同的key，因此没有设置为map格式
	Body     string      // 请求body
}

/*
 * 获取path和query的组合
 */
func (h *HttpRequest) GetPathAndQuery() string {
	if h.Query == "" {
		return h.Path
	}
	return fmt.Sprintf("%s?%s", h.Path, h.Query)
}

/*
 * 获取URL的组合
 */
func (h *HttpRequest) GetUrl() string {
	return fmt.Sprintf("%s://%s%s", h.Scheme, h.Host, h.GetPathAndQuery())
}

/*
 * 获取map结构的header: 如有重复header name ，则会覆盖只保留一个
 */
func (h *HttpRequest) GetHeadersMap() map[string]string {
	headerMap := map[string]string{}
	for i := 0; i < len(h.Headers); i++ {
		headerMap[h.Headers[i][0]] = h.Headers[i][1]
	}
	return headerMap
}
