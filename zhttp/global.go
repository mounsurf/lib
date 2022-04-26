package zhttp

import "net/url"

/**
 * @Desc: ?
 * @Date: 2022/4/26 08:40
 */
var globalProxies map[string]*url.URL

func SetGlobalProxies(proxies map[string]*url.URL) {
	globalProxies = proxies
}
