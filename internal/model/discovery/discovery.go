package discovery

import (
	"../../core/xdatabase"
	"../../core/xlog"
	"../../core/xmonitor"
	"../../core/xnetwork"
	"../../core/xstruct/model_independent"
	modelNode "../node"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"sync"
	"time"
)

type DiscoveryType int

var DiscoveryTypeTitles = []string{"ICMP", "SNMP"}

const (
	ICMP DiscoveryType = iota
	SNMP
)

type Discovery struct {
	Id                string        `json:","` // uuid v4
	Ranges            string        `json:","` // ranges pattern -> iprange+*-*,ipsubnet+*-*,ipcidr+*,...
	DiscoveryType     DiscoveryType `json:","` // type of being alive
	SnmpCredentialKey string        `json:","` // key to snmp credential
	SnmpPort          int           `json:","` // port for snmp connection
	PingPacketSize    int           `json:","` // max 65535 bytes (with 28 bytes overhead of icmp)
	Finished          bool          `json:","` // finished status
	Percentage        float64       `json:","` // percentage of work which is done
	Error             string        `json:","` // error
	Result            string        `json:","` // result as a marshalled(string) json -> {"ip1":[true,23805265],"ip1":[false,"err1"],"ip2":[false,"err2"],...}
}

var waitGroup sync.WaitGroup

func CreateDiscovery(ranges string, discoveryType DiscoveryType, snmpCredentialKey string, snmpPort int) *Discovery {
	return &Discovery{
		Id:                uuid.NewString(),
		Ranges:            ranges,
		DiscoveryType:     discoveryType,
		SnmpCredentialKey: snmpCredentialKey,
		SnmpPort:          snmpPort,
		PingPacketSize:    56,
		Finished:          false,
		Percentage:        0,
		Error:             "",
		Result:            "",
	}
}

func GetDiscoveryByKey(key string) *Discovery {
	value, err := xdatabase.ObjectGetByKey(key)
	if err != nil {
		xlog.LogFatal(err.Error())
	}

	r := &Discovery{}
	for k, v := range value.(map[string]interface{}) {
		_ = model_independent.SetStructField(r, k, v)
	}

	return r
}

func GetDiscoveryKeyById(id string) string {
	values, keys, err := xdatabase.ObjectGetAll("Discovery")
	if err != nil {
		xlog.LogFatal(err.Error())
		return ""
	}

	for i, value := range values {
		r := &Discovery{}
		for k, v := range value.(map[string]interface{}) {
			_ = model_independent.SetStructField(r, k, v)
		}

		if r.Id == id {
			return keys[i]
		}
	}

	return ""
}

func (discovery *Discovery) Start() {

	fmt.Println("DiscoveryStart start")

	discovery.Finished = false
	discovery.Percentage = 0
	discovery.Error = ""
	discovery.Result = ""

	ranges := strings.Split(discovery.Ranges, ",")
	allIps := make([]string, 0)

	for _, v := range ranges {
		parts := strings.Split(v, "+")

		switch parts[0] {
		case "iprange":
			ips := strings.Split(parts[1], "-")
			from := ips[0]
			to := ips[1]
			ipList, e := xnetwork.IpListFromIpRange(from, to)
			if e != nil {
				discovery.Error = e.Error()
				return
			}
			allIps = append(allIps, ipList...)
			break

		case "ipsubnet":
			ips := strings.Split(parts[1], "-")
			ip := ips[0]
			subnet := ips[1]
			ipList, e := xnetwork.IpListFromIpSubnet(ip, subnet)
			if e != nil {
				discovery.Error = e.Error()
				return
			}
			allIps = append(allIps, ipList...)
			break

		case "ipcidr":
			cidr := parts[1]
			ipList, e := xnetwork.IpListFromIpSlash(cidr)
			if e != nil {
				discovery.Error = e.Error()
				return
			}
			allIps = append(allIps, ipList...)
			break
		}
	}

	total := len(allIps)
	counter := 0
	result := ""

	var mu sync.Mutex
	for _, ip := range allIps {
		ip := ip
		waitGroup.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			rtt, e := xmonitor.GetPing(ip, 100*time.Millisecond)
			mu.Lock()
			if e != nil {
				//result += fmt.Sprintf("\"%s\":[false,\"%s\"],", ip, e.Error())
			} else {
				if rtt == -1 {
					//result += fmt.Sprintf("\"%s\":[false,\"Destination Host Unreachable\"],", ip)
				} else {
					//result += fmt.Sprintf("\"%s\":[true,%d],", ip, rtt.Nanoseconds())
					fmt.Printf("ping %s \tis alive (%d)\n", ip, rtt.Nanoseconds())

					aliveType := modelNode.ICMP
					if discovery.DiscoveryType == SNMP {
						aliveType = modelNode.SNMP
					}

					node := modelNode.CreateNode(ip, discovery.SnmpCredentialKey, discovery.SnmpPort, aliveType)
					_ = xdatabase.ObjectStore(*node)
					node.AliveCheckStart(false)
				}
			}
			counter++
			mu.Unlock()
			discovery.Percentage = (float64(counter) / float64(total)) * 100
		}(&waitGroup)
	}
	waitGroup.Wait()
	discovery.Result = result
	discovery.Finished = true

	_ = xdatabase.ObjectUpdateOnKey(*discovery, GetDiscoveryKeyById(discovery.Id))

	fmt.Println("DiscoveryStart end")
}

func (discovery *Discovery) IsFinished() bool {
	return discovery.Finished
}

func (discovery *Discovery) GetResult() string {
	return discovery.Result
}
