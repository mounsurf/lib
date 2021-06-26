package util

import (
	"encoding/binary"
	"errors"
	"net"
	"strconv"
	"strings"
)

func Ip2UInt32(ipAddr string) (uint32, error) {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return 0, errors.New("wrong ip format")
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip), nil
}

func UInt322ip(ipLong uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, ipLong)
	ip := net.IP(ipByte)
	return ip.String()
}
func ParseIpSubnet(ipSubnet string) ([]string, error) {
	invalidErr := errors.New("wrong ip subnet format")
	if IpRegex.MatchString(ipSubnet) {
		ip := net.ParseIP(ipSubnet)
		if ip == nil {
			return nil, invalidErr
		}
		return []string{ip.String()}, nil
	}
	if IpSubnetRegex.MatchString(ipSubnet) {
		subnetInfo := strings.Split(ipSubnet, "/")
		ip, err := Ip2UInt32(subnetInfo[0])
		if err != nil {
			return nil, err
		}
		subnet, err := strconv.Atoi(subnetInfo[1])
		if err != nil || subnet > 32 || subnet < 0 {
			return nil, invalidErr
		}
		left := uint(32 - subnet)
		ipStart := ip >> left << left
		maxUInt32 := uint32(0xffffffff)
		var ipRange uint32
		if left == 0 {
			ipRange = maxUInt32
		} else {
			ipRange = maxUInt32 - maxUInt32>>left<<left + 1
		}
		result := []string{}
		for i := uint32(0); i < ipRange; i++ {
			result = append(result, UInt322ip(ipStart+i))
		}
		return result, nil
	}
	return nil, invalidErr
}

func GetTopDomain(domain string) string {
	pointCount := 0
	var start int
	for start = len(domain) - 1; start >= 0; start-- {
		if domain[start] != '.' {
			continue
		}
		pointCount++
		if pointCount == 2 {
			break
		}
	}
	return domain[start+1:]
}
