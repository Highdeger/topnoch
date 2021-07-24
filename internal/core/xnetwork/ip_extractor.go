package xnetwork

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func IpListFromIpRange(ipRangeFrom, ipRangeTo string) ([]string, error) {
	//if strings.Count(ipRange, "-") != 1 {
	//	return nil, errors.New(fmt.Sprintf("ip range '%s' is not properly formatted", ipRange))
	//} else {
	//parts := strings.Split(ipRange, "-")
	for _, part := range []string{ipRangeFrom, ipRangeTo} {
		r, version := IsIpVersion(part)
		if !r {
			return nil, errors.New(fmt.Sprintf("'%s' is not an ip", part))
		} else {
			if version == 6 {
				return nil, errors.New("ip range of ipv6 not yet supported")
			}
		}
	}

	ip1 := strings.Split(ipRangeFrom, ".")
	ip2 := strings.Split(ipRangeTo, ".")

	var (
		ipList = make([]string, 0)
		aList  = make([]int, 0)
		bList  = make([]int, 0)
	)
	for i := 0; i < len(ip1); i++ {
		if a, err := strconv.Atoi(ip1[i]); err != nil {
			return nil, err
		} else {
			aList = append(aList, a)
		}
		if b, err := strconv.Atoi(ip2[i]); err != nil {
			return nil, err
		} else {
			bList = append(bList, b)
		}
	}

	for i := 0; i < len(aList); i++ {
		if aList[i] > bList[i] {
			aList[i], bList[i] = bList[i], aList[i]
		}
	}

	for a0 := aList[0]; a0 < bList[0]+1; a0++ {
		for a1 := aList[1]; a1 < bList[1]+1; a1++ {
			for a2 := aList[2]; a2 < bList[2]+1; a2++ {
				for a3 := aList[3]; a3 < bList[3]+1; a3++ {
					ipList = append(ipList, fmt.Sprintf("%d.%d.%d.%d", a0, a1, a2, a3))
				}
			}
		}
	}

	return ipList, nil
	//}
}

func IpListFromIpSubnet(ipText, subnetText string) ([]string, error) {
	ip := net.ParseIP(ipText)
	if ip == nil {
		return nil, errors.New(fmt.Sprintf("ip '%s' is not properly formatted", ipText))
	}
	ipSub := net.ParseIP(subnetText)
	if ipSub == nil {
		return nil, errors.New(fmt.Sprintf("subnet '%s' is not properly formatted", subnetText))
	}
	ipMask := net.IPMask{0, 0, 0, 0}
	if (ip.To4() != nil) && (ipSub.To4() != nil) {
		ip = ip.To4()
		ipSub = ipSub.To4()
	} else if (ip.To16() != nil) && (ipSub.To16() != nil) {
		ip = ip.To16()
		ipSub = ipSub.To16()
		ipMask = net.IPMask{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	} else {
		return nil, errors.New(fmt.Sprintf("ip '%s' and subnet '%s' don't match together", ip.String(), ipSub.String()))
	}
	if !IsSubnet(ipSub) {
		return nil, errors.New(fmt.Sprintf("subnet '%s' is not in a correct format (leading ones tailing by zeros)", ipSub.String()))
	}
	for i, byt := range ipSub {
		ipMask[i] = byt
	}
	ipNet := &net.IPNet{
		IP:   ip,
		Mask: ipMask,
	}
	return IpListFromIpNet(ip, ipNet), nil
}

func IpListFromIpSlash(ipSlash string) ([]string, error) {
	ip, ipNet, err := net.ParseCIDR(ipSlash)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("ip cidr (ip/mask) '%s' is not properly formatted", ipSlash))
	}
	return IpListFromIpNet(ip, ipNet), nil
}

func IpListFromIpNet(ip net.IP, ipNet *net.IPNet) []string {
	result := make([]string, 0)
	ip = ip.Mask(ipNet.Mask)
	for ipNet.Contains(ip) {
		result = append(result, ip.String())
		for j := len(ip) - 1; j >= 0; j-- {
			if ip[j] == 255 {
				ip[j] = 0
				continue
			} else {
				ip[j]++
				break
			}
		}
	}
	return result
}
