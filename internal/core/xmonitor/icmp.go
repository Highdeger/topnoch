package xmonitor

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tatsushid/go-fastping"
	"io"
	"net"
	"time"
)

// GetPing ping ip with packetSize in bytes (win->32,most-unix->64)
func GetPing(ip string, timeout time.Duration) (rtt time.Duration, err error) {
	var result time.Duration = -1
	pinger := fastping.NewPinger()
	pinger.MaxRTT = timeout

	addr, err := net.ResolveIPAddr("ip4:icmp", ip)
	if err != nil {
		return -1, err
	}

	pinger.AddIPAddr(addr)
	pinger.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		//fmt.Printf("%s %d\n", addr.IP, rtt.Nanoseconds())
		result = rtt
	}
	pinger.OnIdle = func() {
	}

	err = pinger.Run()
	if err != nil {
		return -1, err
	} else {
		if result == -1 {
			return result, errors.New("destination host unreachable")
		} else {
			return result, nil
		}
	}
}

func GetPings(ips []string, timeout time.Duration, writer io.Writer) {
	pinger := fastping.NewPinger()
	pinger.MaxRTT = timeout

	for _, ip := range ips {
		addr, err := net.ResolveIPAddr("ip4:icmp", ip)
		if err != nil {
			byts := bytes.NewBufferString(fmt.Sprintf("resolving IP '%s' -> err: %s", ip, err.Error())).Bytes()
			_, _ = writer.Write(byts)
		}
		pinger.AddIPAddr(addr)
	}

	pinger.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		e := pinger.Err()
		if e != nil {
			byts := bytes.NewBufferString(fmt.Sprintf("%s -> err: %s\n", addr.IP, e.Error())).Bytes()
			_, _ = writer.Write(byts)
		} else {
			byts := bytes.NewBufferString(fmt.Sprintf("%s: %d\n", addr.IP, rtt)).Bytes()
			_, _ = writer.Write(byts)
		}
	}
	pinger.OnIdle = func() {
	}

	pinger.RunLoop()
}
