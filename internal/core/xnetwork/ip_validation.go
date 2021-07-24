package xnetwork

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

var (
	ipv4Regex, _ = regexp.Compile(`^((25[0-5]|(2[0-4]|1?[0-9])?[0-9])\.){3}(25[0-5]|(2[0-4]|1?[0-9])?[0-9])$`)
	ipv6Regex, _ = regexp.Compile(`^(([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]+|::(ffff(:0{1,4})?:)?((25[0-5]|(2[0-4]|1?[0-9])?[0-9])\.){3}(25[0-5]|(2[0-4]|1?[0-9])?[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1?[0-9])?[0-9])\.){3}(25[0-5]|(2[0-4]|1?[0-9])?[0-9]))$`)
)

func IsIp(ip *string) net.IP {
	return net.ParseIP(*ip)
}

func IsIpVersion(ip string) (result bool, version int) {
	if parsedIp := net.ParseIP(ip); parsedIp == nil {
		return false, 0
	} else {
		if parsedIp.To4() != nil {
			return true, 4
		} else {
			return true, 6
		}
	}
}

func IsIpRegex(ip string) (bool, int) {
	ip = strings.Trim(ip, " ")
	if ipv4Regex.MatchString(ip) {
		return true, 4
	}
	if ipv6Regex.MatchString(ip) {
		return true, 6
	}
	return false, 0
}

func IsHost(host string) (result bool, addressType string, connectionError error) {
	result = false
	addressType = ""
	connectionError = nil
	switch strings.Count(host, ":") {
	case 0:
		if r, version := IsIpVersion(host); r {
			result = true
			addressType = fmt.Sprintf("ipv%d", version)
		} else {
			if addrs, err := net.LookupHost(host); err != nil {
				connectionError = err
			} else {
				if len(addrs) > 0 {
					result = true
					addressType = "hostname"
				} else {
					connectionError = errors.New(fmt.Sprintf("check the host '%s': dns could not resolve any address for the host", host))
				}
			}
		}
	case 1:
		if strings.Count(host, "://") > 0 {
			if strings.HasPrefix(host, "http://") {
				addressType = "hostname.http"
			} else if strings.HasPrefix(host, "https://") {
				addressType = "hostname.https"
			}
			if addrs, err := net.LookupHost(strings.Split(host, "://")[1]); err == nil {
				if len(addrs) > 0 {
					result = true
				} else {
					addressType = ""
					connectionError = errors.New(fmt.Sprintf("check the host '%s': dns can not resolve any address for the hostname", host))
				}
			} else {
				addressType = ""
				connectionError = err
			}
		} else {
			parts := strings.Split(host, ":")
			port := parts[1]
			if n, err := strconv.Atoi(port); err == nil {
				if (n >= 0) || (n <= 65353) {
					if r, _ := IsIpVersion(parts[0]); r {
						result = true
						addressType = "ipv4+port"
					} else {
						if addrs, err := net.LookupHost(parts[0]); err == nil {
							if len(addrs) > 0 {
								result = true
								addressType = "hostname+port"
							} else {
								connectionError = errors.New(fmt.Sprintf("check the host '%s': wrong format of ip or a dead-end hostname", parts[0]))
							}
						} else {
							connectionError = errors.New(err.Error())
						}
					}
				} else {
					connectionError = errors.New(fmt.Sprintf("check the port '%d': not in the range (0<p<65353)", n))
				}
			} else {
				connectionError = errors.New(fmt.Sprintf("check the port '%s': not an integer", port))
			}
		}
	case 7:
		if r, _ := IsIpVersion(host); r {
			result = true
			addressType = fmt.Sprintf("ipv6")
		} else {
			connectionError = errors.New(fmt.Sprintf("check the host '%s': wrong format of ipv6", host))
		}
	case 8:
		if strings.HasPrefix(host, "[") {
			host = host[1:]
		}
		if strings.Count(host, "]:") == 1 {
			parts := strings.Split(host, "]:")
			if r, _ := IsIpVersion(parts[0]); r {
				port := parts[1]
				if n, err := strconv.Atoi(port); err == nil {
					if (n >= 0) || (n <= 65353) {
						result = true
						addressType = "ipv6+port"
					} else {
						connectionError = errors.New(fmt.Sprintf("check the port '%d': not in the range (0<p<65353)", n))
					}
				} else {
					connectionError = errors.New(fmt.Sprintf("check the port '%s': not a number", port))
				}
			} else {
				connectionError = errors.New(fmt.Sprintf("check the host '%s': wrong format of ipv6", parts[0]))
			}
		} else {
			connectionError = errors.New(fmt.Sprintf("check the host '%s': wrong format of [ipv6]:port", host))
		}
	}
	return
}

func IsSubnet(subnet net.IP) bool {
	result := true
	zeros := false
subnetCheck:
	for _, byt := range subnet {
		for _, ch := range fmt.Sprintf("%08b", byt) {
			if zeros && (string(ch) == "1") {
				result = false
				break subnetCheck
			}
			if !zeros && (string(ch) == "0") {
				zeros = true
			}
		}
	}
	return result
}

func IpMaskFrom(ip net.IP) net.IPMask {
	mask := make([]byte, len(ip))
	for i, byt := range ip {
		mask[i] = byt
	}
	return mask
}

func IsIpPort(ipPort string) bool {
	if strings.Contains(ipPort, ":") {
		ip_port := strings.Split(ipPort, ":")
		r, _ := IsIpVersion(ip_port[0])
		if r {
			n, err := strconv.Atoi(ip_port[1])
			if err != nil {
				return false
			} else {
				if (n >= 0) && (n <= 65353) {
					return true
				} else {
					return false
				}
			}
		} else {
			return false
		}
	} else {
		return false
	}
}
