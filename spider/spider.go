package spider

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/mounsurf/lib/zhttp"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var spiderUrlRegex = regexp.MustCompile(`['"](http://[^"']*|https://[^"']*|/[^"']*)['"]`)

func GetUrlPageLinks(urlStr string, cookie string, sameOrigin bool) (result []string) {
	result = []string{}
	resultMap := map[string]struct{}{}
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		return
	}
	resp, err := zhttp.Get(urlStr, &zhttp.RequestOptions{
		RequestTimeout:  time.Second * 5,
		DialTimeout:     time.Second * 5,
		IgnoreCertError: true,
		Headers: map[string]string{
			"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
			"Cookie": cookie,
		}})
	if err != nil {
		return
	}
	respBody := resp.GetBodyString()
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(respBody))
	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		link, exists := selection.Attr("href")
		if exists {
			link, err = zhttp.GetAbsUrl(link, urlObj.Host, urlObj.Scheme)
			if err != nil {
				return
			}
			if sameOrigin {
				if strings.HasPrefix(link, urlObj.Scheme+"://"+urlObj.Host+"/") {
					resultMap[link] = struct{}{}
				}
			} else {
				resultMap[link] = struct{}{}
			}
		}
	})
	doc.Find("script").Each(func(i int, selection *goquery.Selection) {
		urlArrList := spiderUrlRegex.FindAllStringSubmatch(selection.Text(), -1)
		for _, urlArr := range urlArrList {
			if len(urlArr) == 2 {
				link, err := zhttp.GetAbsUrl(urlArr[1], urlObj.Host, urlObj.Scheme)
				if err != nil {
					return
				}
				linkObj, err := url.Parse(link)
				if err != nil || linkObj.Host == "" {
					return
				}
				for _, suffix := range []string{".js", ".css", ".map", ".jpg", ".gif", ".svg", ".png"} {
					if strings.HasSuffix(linkObj.Path, suffix) {
						return
					}
				}
				if sameOrigin {
					if strings.HasPrefix(link, urlObj.Scheme+"://"+urlObj.Host+"/") {
						resultMap[link] = struct{}{}
					}
				} else {
					resultMap[link] = struct{}{}
				}
			}
		}

	})
	for link, _ := range resultMap {
		result = append(result, link)
	}
	return result
}

func scrawlUrlMap(urlMap map[string]struct{}, cookie string, sameOrigin bool) map[string]struct{} {
	allUrl := map[string]struct{}{}
	for urlStr, _ := range urlMap {
		urlStrLinks := GetUrlPageLinks(urlStr, cookie, sameOrigin)
		for i := 0; i < len(urlStrLinks) && i < 20; i++ {
			allUrl[urlStrLinks[i]] = struct{}{}
		}
	}
	return allUrl
}

func ScrawlUrl(urlStr string, cookie string, depth int, sameOrigin bool) []string {
	allUrl := []string{}
	var layerUrlMap = make([]map[string]struct{}, depth+1)
	layerUrlMap[0] = map[string]struct{}{
		urlStr: {},
	}
	allUrl = append(allUrl, urlStr)
	for i := 1; i < depth+1; i++ {
		layerUrlMap[i] = map[string]struct{}{}
		urlLinksMap := scrawlUrlMap(layerUrlMap[i-1], cookie, sameOrigin)
		for urlLink, _ := range urlLinksMap {
			exist := false
			for j := 0; j < i; j++ {
				if _, ok := layerUrlMap[j][urlLink]; ok {
					exist = true
				}
			}
			if !exist {
				layerUrlMap[i][urlLink] = struct{}{}
				allUrl = append(allUrl, urlLink)
			}
		}
	}
	return allUrl
}
