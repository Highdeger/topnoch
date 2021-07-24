package main

import (
	"../internal/core"
	"../internal/server"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func generateIPs(length int) *[]string {
	list := make([]string, 0)
	// generate valid ipv4
	t := 0
	for i := 0; i < ((length / 4) + t); i++ {
		part := make([]byte, 1)
		ip := ""
		for j := 0; j < 4; j++ {
			rand.Read(part)
			if j != 0 {
				ip += "."
			}
			ip += fmt.Sprintf("%d", part[0])
		}
		if core.ItemIndex(ip, &list) == -1 {
			list = append(list, ip)
		} else {
			t += 1
		}
	}
	// generate valid/invalid ipv4
	t = 0
	for i := 0; i < ((length / 4) + t); i++ {
		ip := ""
		for j := 0; j < 4; j++ {
			rand.Seed(time.Now().UnixNano())
			part := rand.Intn(400) - 50
			if j != 0 {
				ip += "."
			}
			ip += fmt.Sprintf("%d", part)
		}
		if core.ItemIndex(ip, &list) == -1 {
			list = append(list, ip)
		} else {
			t += 1
		}
	}
	// generate valid ipv6
	t = 0
	for i := 0; i < ((length / 4) + t); i++ {
		part := make([]byte, 2)
		ip := ""
		for j := 0; j < 8; j++ {
			rand.Read(part)
			if j != 0 {
				ip += ":"
			}
			ip += fmt.Sprintf("%x", part)
		}
		if core.ItemIndex(ip, &list) == -1 {
			list = append(list, ip)
		} else {
			t += 1
		}
	}
	// generate valid/invalid ipv6
	t = 0
	var charList []string
	charList = append(charList, "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f")
	for i := 0; i < ((length / 4) + t); i++ {
		ip := ""
		for j := 0; j < 8; j++ {
			part := ""
			for k := 0; k < 6; k++ {
				rand.Seed(time.Now().UnixNano())
				if k == 0 {
					if rand.Intn(100) < 70 {
						part += charList[rand.Intn(len(charList))]
					}
				} else if k == 1 {
					if rand.Intn(100) < 70 {
						part += charList[rand.Intn(len(charList))]
					}
				} else {
					part += charList[rand.Intn(len(charList))]
				}
			}
			if j != 0 {
				ip += ":"
			}
			ip += part
		}
		if core.ItemIndex(ip, &list) == -1 {
			list = append(list, ip)
		} else {
			t += 1
		}
	}
	return &list
}

func generateHosts() *[]string {
	list := make([]string, 0)
	list = append(list, "98.218.135.7:65353")
	list = append(list, "98.218.135.7:65354")
	list = append(list, "255.4.131.68:22")
	list = append(list, "256.4.131.68:22")
	list = append(list, "217.4.131.68:a")
	list = append(list, "217.4.131.68")
	list = append(list, "217.4.131 68")
	list = append(list, "217.4.13a.68")
	list = append(list, "2001:db8:3c4f:15:0::1a2f:1a2b")
	list = append(list, "[2001:db8:3c4f:0015:0::1a2f:1a2b]:23")
	list = append(list, "2001:db8:3c4f:0015:0::1a2f:1a2b]:23")
	list = append(list, "2001:db8:3c4g:0015:0::1a2f:1a2b]:23")
	list = append(list, "http://www.google.com")
	list = append(list, "https://www.google.com")
	list = append(list, "https://google.com")
	list = append(list, "google.com")
	list = append(list, "yahoo.com")
	list = append(list, "gfdgfdssfdgfdsgsdfgsdfdsfds.com")
	return &list
}

func testIpCheckSpeed(length int, epoch int) {
	var (
		isIpTotal        int64
		isIpVersionTotal int64
		isIpRegexTotal   int64
	)
	isIpTrue := 0
	isIpFalse := 0
	isIpVersionTrue := 0
	isIpVersionFalse := 0
	isIpRegexTrue := 0
	isIpRegexFalse := 0

	for x := 0; x < epoch; x++ {
		testIPs := make([]string, 0)
		for _, v := range *generateIPs(10000) {
			testIPs = append(testIPs, v)
		}
		log.Printf("number of IPs generated is %d, epoch: %d/%d\n", len(testIPs)+1, x, epoch)

		t0 := time.Now()
		for i := 0; i < len(testIPs); i++ {
			r := server.IsIp(&testIPs[i])
			if r != nil {
				isIpTrue += 1
			} else {
				isIpFalse += 1
			}
		}
		isIpTotal += time.Since(t0).Nanoseconds()

		t0 = time.Now()
		for i := 0; i < len(testIPs); i++ {
			r, _ := server.IsIpVersion(testIPs[i])
			if r {
				isIpVersionTrue += 1
			} else {
				isIpVersionFalse += 1
			}
		}
		isIpVersionTotal += time.Since(t0).Nanoseconds()

		t0 = time.Now()
		for i := 0; i < len(testIPs); i++ {
			r, _ := server.IsIpRegex(testIPs[i])
			if r {
				isIpRegexTrue += 1
			} else {
				isIpRegexFalse += 1
			}
		}
		isIpRegexTotal += time.Since(t0).Nanoseconds()
	}

	log.Printf("total number of IPs tested: %d\n", length*epoch)
	log.Printf("method IsIp        : % 10d ns -> Result(true=%d, false=%d)\n", isIpTotal/int64(epoch), isIpTrue, isIpFalse)
	log.Printf("method IsIpVersion : % 10d ns -> Result(true=%d, false=%d)\n", isIpVersionTotal/int64(epoch), isIpVersionTrue, isIpVersionFalse)
	log.Printf("method IsIpRegex   : % 10d ns -> Result(true=%d, false=%d)\n", isIpRegexTotal/int64(epoch), isIpRegexTrue, isIpRegexFalse)
}

func testHostsCheck() {
	testHosts := make([]string, 0)
	for _, v := range *generateIPs(10000) {
		testHosts = append(testHosts, v)
	}
	for _, v := range *generateHosts() {
		testHosts = append(testHosts, v)
	}
	log.Printf("total json: %d\n", len(testHosts))
	for _, host := range testHosts {
		if r, t, err := server.IsHost(host); r {
			log.Printf("host '%s' is a host of type '%s'\n", host, t)
		} else if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	testIpCheckSpeed(10000, 10)
	testHostsCheck()
}
