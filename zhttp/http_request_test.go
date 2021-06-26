package zhttp

import (
	"fmt"
	"lib/util"
)

/**
 * @Desc:
 * @Date: 2021/3/6 18:52
 */
func demo() {
	data := `GET / HTTP/1.1
Host: www.baidu.com
Connection: close
Upgrade-Insecure-Requests: 1
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.192 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
Purpose: prefetch
Sec-Fetch-Site: none
Sec-Fetch-Mode: navigate
Sec-Fetch-User: ?1
Sec-Fetch-Dest: document
Accept-Encoding: gzip, deflate
Accept-Language: zh-CN,zh;q=0.9

`
	httpRequest, err := ParseHttpRequest(data, "https")
	util.CheckError(err)
	fmt.Println()
	resp, err := Request(httpRequest.Method, httpRequest.GetUrl(), &RequestOptions{
		RawData: httpRequest.Body,
		Headers: httpRequest.GetHeadersMap(),
		Proxies: GetLocalProxy(8889),
	})
	util.CheckError(err)
	fmt.Println(resp.GetBodyString())
}
