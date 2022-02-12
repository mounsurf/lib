package zhttp

import (
	"encoding/json"
	"fmt"
	"github.com/mounsurf/lib/util"
	"github.com/mounsurf/lib/zlog"
	"strings"
	"sync"
	"time"
)

const (
	defaultServer      = "http://localhost:9200/zhttp/log/" //默认server
	defaultQueueLength = 50                                 //默认50长度
)

var (
	esFlag    bool           //是否自动记录到ES
	esServer  string         //"http://localhost:9200/zhttp/log/"
	esTag     string         //标签
	esTagHash string         //标签hash
	esQueue   chan *EsData   //排队队列
	esWg      sync.WaitGroup //wait group
)

type EsData struct {
	Tag             string            `json:"tag"` // 标签
	TagHash         string            `json:"tag_hash"`
	Hash            string            `json:"hash"`   // 请求hash，请求body的hash
	Scheme          string            `json:"scheme"` // 协议
	Method          string            `json:"method"`
	Hostname        string            `json:"hostname"`
	Host            string            `json:"host"`
	Port            string            `json:"port"`
	Url             string            `json:"url"`
	Path            string            `json:"path"`
	ResponseBody    string            `json:"response_body"`
	Status          int               `json:"status"`
	RequestHeaders  map[string]string `json:"request_headers"`
	RequestHeader   string            `json:"request_header"`
	ResponseHeaders map[string]string `json:"response_headers"`
	ResponseHeader  string            `json:"response_header"`
	Title           string            `json:"title"`
	Time            string            `json:"time"`
}

func SetEsServer(server string) {
	esServer = server
	esQueue = make(chan *EsData, 50)
}

func CloseAutoLogToEs() {
	esFlag = false
	close(esQueue)
}
func SetAutoLogToEs(server string, tag string, queueLength int) {
	esFlag = true
	if server == "" {
		esServer = defaultServer
	} else {
		esServer = server
	}
	//判断是否为空、是否已经关闭，未关闭则关闭
	if esQueue != nil {
		_, isClose := <-esQueue
		if !isClose {
			close(esQueue)
		}
	}
	if queueLength <= 0 {
		queueLength = defaultQueueLength
	}
	esQueue = make(chan *EsData, queueLength)

	esTag = tag
	esTagHash = util.Md5([]byte(tag + time.Now().String()))
	go func() {
		consumeEsData()
	}()
}

func logToEs(resp *Response) {
	// 防止递归死循环
	hostname := resp.RawResponse.Request.URL.Hostname()
	if hostname == "127.0.0.1" || hostname == "localhost" || strings.Contains(resp.GetFinalUrl(), esServer) {
		return
	}

	esData := &EsData{
		Tag:      esTag,
		TagHash:  esTagHash,
		Hash:     util.Md5([]byte(resp.GetFinalUrl())),
		Scheme:   resp.RawResponse.Request.URL.Scheme,
		Method:   resp.RawResponse.Request.Method,
		Hostname: resp.RawResponse.Request.URL.Hostname(),
		Host:     resp.RawResponse.Request.URL.Host,
		Url:      resp.GetFinalUrl(),
		Path:     resp.RawResponse.Request.URL.Path,
		Status:   resp.StatusCode(),
		Title:    resp.GetTitle(),
		Time:     time.Now().Format("2006-01-02 15:04:05"),
	}
	port := resp.RawResponse.Request.URL.Port()
	if port == "" {
		if esData.Scheme == "http" {
			port = "80"
		} else if esData.Scheme == "https" {
			port = "443"
		}
	}
	esData.Port = port
	requestHeaders := map[string]string{}
	requestHeader := ""
	for k, v := range resp.RawResponse.Request.Header {
		if len(v) >= 1 {
			requestHeaders[k] = v[0]
		}
		for _, header := range v {
			requestHeader += fmt.Sprintf("%s: %s\n", k, header)
		}
	}
	esData.RequestHeaders = requestHeaders
	esData.RequestHeader = requestHeader

	responseHeaders := map[string]string{}
	responseHeader := ""
	contentType := ""
	responseBody := ""
	for k, v := range resp.RawResponse.Header {
		if len(v) >= 1 {
			if contentType == "" && k == "Content-Type" {
				contentType = v[0]
			}
			responseHeaders[k] = v[0]
		}
		for _, header := range v {
			responseHeader += fmt.Sprintf("%s: %s\n", k, header)
		}
	}
	if strings.Contains(strings.ToLower(contentType), "image/") {
		responseBody = "[image]"
	} else {
		responseBody = resp.GetBodyString()
	}
	esData.ResponseHeaders = responseHeaders
	esData.ResponseHeader = responseHeader
	esData.ResponseBody = responseBody
	//应当是异步的
	go func() {
		addToEsQueue(esData)
	}()
}

func addToEsQueue(esData *EsData) {
	esWg.Add(1)
	esQueue <- esData
}

func consumeEsData() {
	for esData := range esQueue {
		jsonData, err := json.Marshal(esData)
		if err == nil {
			resp, err := Post(esServer, &RequestOptions{
				ContentType: "application/json",
				JSON:        string(jsonData),
			})
			if err != nil {
				zlog.Error(err)
			} else if strings.HasPrefix(resp.GetBodyString(), `{"error"`) {
				zlog.Error(resp.GetBodyString())
			}
		}
		esWg.Done()
	}
}
func WaitEs() {
	esWg.Wait()
}

/**
demo:
func main(){
	zhttp.SetAutoLogToEs("", "test")
	zhttp.Get("http://www.baidu.com", &zhttp.RequestOptions{
		UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_1 like Mac OS X) AppleWebKit/603.1.30 (KHTML, like Gecko) Version/10.0 Mobile/14E304 Safari/602.1 Edg/98.0.4758.80",
	})
	zlog.Info(123)
	zhttp.WaitEs()
}
*/
