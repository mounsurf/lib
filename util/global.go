package util

import "regexp"

var (
	Ipv4Regex     = regexp.MustCompile(`^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$`)
	IpRegex       = Ipv4Regex
	IpSubnetRegex = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}/\d{1,2}$`)
)
