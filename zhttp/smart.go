package zhttp

import (
	"errors"
	"github.com/mounsurf/lib/zlog"
	"strings"
)

// 智能模式

var (
	autoIgnoreCertErr bool // 如果请求的host是IP，那么自动忽略证书错误，默认不开启
	autoRefuseGovSite bool // gov域名默认不做任何请求
	autoRefuseEduSite bool // edu域名默认不做任何请求
)

func init() {
	SetBestPractice()
}

func SetAutoIgnoreCertErr(value bool) {
	autoIgnoreCertErr = value
	if value {
		zlog.Info("当host为IP时，会自动忽略证书错误")
	} else {
		zlog.Info("当host为IP时，不会自动忽略证书错误")
	}
}

func SetAutoRefuseGovSite(value bool) {
	autoRefuseGovSite = value
	if value {
		zlog.Info("自动拒绝连接gov站点")
	} else {
		zlog.Info("可连接gov站点")
	}
}

func SetAutoRefuseEduSite(value bool) {
	autoRefuseEduSite = value
	if value {
		zlog.Info("自动拒绝连接edu站点")
	} else {
		zlog.Info("可连接edu站点")
	}
}

// SetBestPractice 最佳实践
func SetBestPractice() {
	SetAutoIgnoreCertErr(true)
	SetAutoRefuseGovSite(true)
	SetAutoRefuseEduSite(true)
}

func checkAvailable(host string) error {
	if autoRefuseGovSite && strings.Contains(host, ".gov.") {
		return errors.New("gov site is not available")
	}
	if autoRefuseEduSite && strings.Contains(host, ".edu.") {
		return errors.New("edu site is not available")
	}
	return nil
}
