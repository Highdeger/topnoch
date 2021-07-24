package node

import (
	"../../core/xdatabase"
	"../../core/xlog"
	"../../core/xmonitor"
	"../../core/xstruct/model_independent"
	"../credentials"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"sync"
	"time"
)

type AliveType int

var AliveTypeTitle = []string{"ICMP", "SNMP"}

const (
	ICMP AliveType = iota
	SNMP
)

type MachineType int

var MachineTypeTitles = []string{"Unknown", "Windows", "Linux", "Cisco"}

const (
	Unknown MachineType = iota
	Windows
	Linux
	Cisco
)

type ArchType int

var ArchTypeTitles = []string{"ArchUnknown", "Arch32Bit", "Arch64Bit"}

const (
	ArchUnknown ArchType = iota
	Arch32Bit
	Arch64Bit
)

type HardwareState int

var HardwareStateTitles = []string{"Up", "Down", "Warning"}

const (
	NotDetermined HardwareState = iota
	Up
	Down
	Warning
)

type Node struct {
	Id                 string        `json:","` // uuid v4
	Ip                 string        `json:","` // ip of node
	Name               string        `json:","` // name of node (initialized to hostname)
	AliveType          AliveType     `json:","` // type of being alive
	AliveInterval      time.Duration `json:","` // interval to check for being alive
	IsCheckingAlive    bool          `json:","` // is checking
	AliveError         string        `json:","` // the error of checking alive, nil if no error
	SnmpCredentialKey  string        `json:","` // key to snmp credential
	SnmpPort           int           `json:","` // port for snmp connection
	OsCredentialKey    string        `json:","` // key to os credential
	MachineType        MachineType   `json:","` // machine type
	ArchType           ArchType      `json:","` // architecture type
	OsVersion          string        `json:","` // version of os in details
	PingPacketSize     int           `json:","` // ping packet size (max 65535 bytes with 28 bytes overhead of icmp)
	PingAverage        int64         `json:","` // average rtt
	PingQuantity       int64         `json:","` // number of all ping attempt
	PingQuantityFailed int64         `json:","` // number of failed ping attempt to calculate packet loss
	PingTimeout        time.Duration `json:","` // ping timeout (default=50ms)
	HardwareState      HardwareState `json:","` // hardware state
}

func CreateNode(ip string, snmpCredentialKey string, port int, aliveType AliveType) *Node {
	return &Node{
		Id:                 uuid.NewString(),
		Ip:                 ip,
		Name:               "",
		AliveType:          aliveType,
		AliveInterval:      1 * time.Second,
		IsCheckingAlive:    false,
		AliveError:         "",
		SnmpCredentialKey:  snmpCredentialKey,
		SnmpPort:           port,
		OsCredentialKey:    "",
		MachineType:        0,
		ArchType:           0,
		OsVersion:          "",
		PingPacketSize:     56,
		PingAverage:        0,
		PingQuantity:       0,
		PingQuantityFailed: 0,
		PingTimeout:        200 * time.Millisecond,
		HardwareState:      NotDetermined,
	}
}

func GetNodeByKey(key string) *Node {
	value, err := xdatabase.ObjectGetByKey(key)
	if err != nil {
		xlog.LogFatal(err.Error())
	}

	r := &Node{}
	for k, v := range value.(map[string]interface{}) {
		_ = model_independent.SetStructField(r, k, v)
	}

	return r
}

func GetNodeKeyById(id string) string {
	values, keys, err := xdatabase.ObjectGetAll("Node")
	if err != nil {
		xlog.LogFatal(err.Error())
		return ""
	}

	for i, value := range values {
		r := &Node{}
		for k, v := range value.(map[string]interface{}) {
			_ = model_independent.SetStructField(r, k, v)
		}

		if r.Id == id {
			return keys[i]
		}
	}

	return ""
}

func (node *Node) IsAlive() (ok bool, err string) {
	if node.AliveError == "" {
		return true, ""
	} else {
		return false, node.AliveError
	}
}

func (node *Node) AliveCheckStart(once bool) {
	node.IsCheckingAlive = true
	var mutex sync.Mutex
	go func(mu *sync.Mutex) {
		for node.IsCheckingAlive {
			switch node.AliveType {
			case ICMP:
				_, err := node.AlivePing()
				if err != nil {
					node.AliveError = err.Error()
					node.AliveCheckStop()
				} else {
					node.AliveError = ""
				}

			case SNMP:
				err := node.AliveSnmp()
				if err != nil {
					node.AliveError = err.Error()
					node.AliveCheckStop()
				} else {
					node.AliveError = ""
				}
			}

			if once {
				node.IsCheckingAlive = false
				_ = node.SnmpCategorize()
			}

			mu.Lock()
			_ = xdatabase.ObjectUpdateOnKey(node, GetNodeKeyById(node.Id))
			mu.Unlock()

			time.Sleep(node.AliveInterval)
		}
	}(&mutex)
}

func (node *Node) AliveCheckStop() {
	node.IsCheckingAlive = false
}

func (node *Node) AlivePing() (time.Duration, error) {
	rtt, err := xmonitor.GetPing(node.Ip, node.PingTimeout)
	if err != nil {
		node.PingQuantity++
		node.PingQuantityFailed++
		return rtt, err
	}

	node.PingQuantity++
	oldAverage := node.PingAverage
	node.PingAverage = (oldAverage + rtt.Nanoseconds()) / node.PingQuantity

	return rtt, err
}

func (node *Node) SnmpCategorize() error {
	if node.SnmpCredentialKey == "" {
		fmt.Println("ERROR:", "no snmp credential key found")
		return errors.New("no snmp credential key found")
	}

	value, err := xdatabase.ObjectGetByKey(node.SnmpCredentialKey)
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		return err
	}

	//snmpCredential := value.(credentials.SnmpCredential)
	snmpCredential := &credentials.SnmpCredential{}
	for k, v := range value.(map[string]interface{}) {
		_ = model_independent.SetStructField(snmpCredential, k, v)
	}

	sysDescText := xmonitor.FetchSysDescr(snmpCredential.MixWithAddress(fmt.Sprintf("%s:%d", node.Ip, node.SnmpPort)))
	if strings.Trim(sysDescText, " ") == "" {
		node.MachineType = Unknown
		node.Name = fmt.Sprintf("Unknown (%s)", node.Ip)
		node.OsVersion = ""
		//fmt.Println("can't fetch system description by snmp")
		//return errors.New("can't fetch system description by snmp")
	}
	sysDescription := strings.Split(sysDescText, " ")

	switch sysDescription[0] {
	// Linux devdeb 4.19.0-14-amd64 #1 SMP Debian 4.19.171-2 (2021-01-30) x86_64
	case "Linux":
		node.MachineType = Linux
		node.Name = sysDescription[1]

		version := ""
		for i := 2; i < len(sysDescription); i++ {
			if i == len(sysDescription)-1 {
				switch sysDescription[i] {
				case "x86_64":
					node.ArchType = Arch64Bit
				case "x86":
					node.ArchType = Arch32Bit
				}
			} else {
				version += sysDescription[i]
				if i != len(sysDescription)-2 {
					version += " "
				}
			}
		}
		node.OsVersion = version

	case "Windows":
		node.MachineType = Windows
		node.Name = sysDescription[1]

		version := ""
		for i := 2; i < len(sysDescription); i++ {
			if i == len(sysDescription)-1 {
				switch sysDescription[i] {
				case "x86_64":
					node.ArchType = Arch64Bit
				case "x86":
					node.ArchType = Arch32Bit
				}
			} else {
				version += sysDescription[i]
				if i != len(sysDescription)-2 {
					version += " "
				}
			}
		}
		node.OsVersion = version

	case "Cisco":
		node.MachineType = Cisco

	// Hardware: x86 Family 15 Model 12 Stepping 0 AT/AT COMPATIBLE - Software: Windows 2000 Version 5.0 (Build 2195 Uniprocessor Free)
	case "Hardware:":
		if strings.Contains(sysDescText, "Software: Windows") {
			node.MachineType = Windows

			sysDescText = strings.Trim(strings.TrimPrefix(sysDescText, "Hardware:"), " ")
			split := strings.Split(sysDescText, " - Software: ")
			version := strings.Trim(strings.TrimPrefix(split[1], "Windows"), " ")
			node.OsVersion = version
		}
	}
	return nil
}

func (node *Node) AliveSnmp() error {
	if node.SnmpCredentialKey == "" {
		return errors.New("no snmp credential key found")
	}

	value, err := xdatabase.ObjectGetByKey(node.SnmpCredentialKey)
	if err != nil {
		return err
	}

	snmpCredential := &credentials.SnmpCredential{}
	for k, v := range value.(map[string]interface{}) {
		_ = model_independent.SetStructField(snmpCredential, k, v)
	}
	//snmpCredential := value.(credentials.SnmpCredential)

	snmpArg := snmpCredential.MixWithAddress(fmt.Sprintf("%s:%d", node.Ip, node.SnmpPort))
	sysDescText := xmonitor.FetchSysDescr(snmpArg)

	cpu1m := xmonitor.FetchAverageCpuLoads(snmpArg)[0]
	nodeKey, e := xdatabase.ObjectGetKey(node)
	fmt.Println("error", e)
	if e != nil {
		return e
	}
	fmt.Println("nodeKey", nodeKey)
	xdatabase.ParamStore(fmt.Sprint(cpu1m), "cpu1m", nodeKey)

	if strings.Trim(sysDescText, " ") == "" {
		return errors.New("can't fetch system description by snmp")
	} else {
		return nil
	}
}
