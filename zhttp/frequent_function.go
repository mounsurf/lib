package zhttp

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"regexp"
	"strings"
)

var (
	schemeUrlRegex    = regexp.MustCompile(`^([a-z]+:).*`)
	relativePathRegex = regexp.MustCompile(`^[\w\d].*`)
)

func GetResponse(urlStr, cookie string) (*Response, error) {
	resp, err := Get(urlStr, &RequestOptions{
		IgnoreCertError: true,
		RawCookie:       cookie,
		DisableRedirect: true,
		Headers: map[string]string{
			"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		},
	})
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetResponseInfo(urlStr, cookie string) (string, map[string]string, error) {
	resp, err := Get(urlStr, &RequestOptions{
		IgnoreCertError: true,
		RawCookie:       cookie,
		DisableRedirect: true,
		Headers: map[string]string{
			"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		},
	})
	if err != nil {
		return "", nil, err
	}
	return resp.GetBodyString(), resp.HeadersMap(), err
}

func GetRawResponse(urlStr, cookie string, disableRedirect bool) (*Response, error) {
	resp, err := Get(urlStr, &RequestOptions{
		IgnoreCertError: true,
		RawCookie:       cookie,
		DisableRedirect: disableRedirect,
		Headers: map[string]string{
			"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		},
	})
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetRawHeaders(urlStr, cookie string) (string, error) {
	resp, err := Get(urlStr, &RequestOptions{
		IgnoreCertError: true,
		RawCookie:       cookie,
		DisableRedirect: true,
		Headers: map[string]string{
			"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		},
	})
	if err != nil {
		return "", err
	}
	return resp.RawHeaders(), err
}

func GetPostResponseInfo(urlStr, cookie, data, contentType string) (string, map[string]string, error) {
	resp, err := Post(urlStr, &RequestOptions{
		IgnoreCertError: true,
		RawCookie:       cookie,
		DisableRedirect: true,
		RawData:         data,
		ContentType:     contentType,
		Headers: map[string]string{
			"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		},
	})
	if err != nil {
		return "", nil, err
	}
	return resp.GetBodyString(), resp.HeadersMap(), err
}

func GetAbsUrl(urlStr string, host string, scheme string) (string, error) {
	if strings.HasPrefix(urlStr, "http://") || strings.HasPrefix(urlStr, "https://") {
		return urlStr, nil
	} else if strings.HasPrefix(urlStr, "//") {
		return scheme + ":" + urlStr, nil
	} else if strings.HasPrefix(urlStr, "/") {
		return scheme + "://" + host + urlStr, nil
	} else if !schemeUrlRegex.MatchString(urlStr) && relativePathRegex.MatchString(urlStr) {
		return scheme + "://" + host + "/" + urlStr, nil
	} else {
		return "", errors.New("Not Valid Url")
	}
}

// GetATagLink 0:全部 1:同源 2:非同源
func GetATagLink(body string, urlObj *url.URL, mode int) map[string]struct{} {
	resultMap := map[string]struct{}{}
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(body))
	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		link, exists := selection.Attr("href")
		var err error
		if exists {
			link, err = GetAbsUrl(link, urlObj.Host, urlObj.Scheme)
			if err != nil {
				return
			}
			if mode == 0 {
				resultMap[link] = struct{}{}
			} else if mode == 1 {
				if strings.HasPrefix(link, urlObj.Scheme+"://"+urlObj.Host+"/") {
					resultMap[link] = struct{}{}
				}
			} else if mode == 2 {
				if !strings.HasPrefix(link, urlObj.Scheme+"://"+urlObj.Host+"/") {
					resultMap[link] = struct{}{}
				}
			}
		}
	})
	return resultMap
}
