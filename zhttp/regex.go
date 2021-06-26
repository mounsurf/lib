package zhttp

import "regexp"

var (
	titleRegex = regexp.MustCompile(`(?i)<title.*?>([\s\S]*?)</title>`)
	gbkRegex   = regexp.MustCompile(`(?i)<meta\s+[^>]*charset="?(gbk|gb2312)`)
)
