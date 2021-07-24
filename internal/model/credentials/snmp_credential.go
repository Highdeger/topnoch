package credentials

import (
	"github.com/k-sone/snmpgo"
	"reflect"
	"time"
)

type SnmpCredential struct {
	Version          snmpgo.SNMPVersion   `json:","` // SNMP version to use
	Network          string               `json:","` // See net.Dial parameter (The default is `udp`)
	Timeout          time.Duration        `json:","` // Request timeout (The default is 5sec)
	Retries          uint                 `json:","` // Number of retries (The default is `0`)
	MessageMaxSize   int                  `json:","` // Maximum size of an SNMP message (The default is `1400`)
	Community        string               `json:","` // Community (V1 or V2c specific)
	UserName         string               `json:","` // Security name (V3 specific)
	SecurityLevel    snmpgo.SecurityLevel `json:","` // Security level (V3 specific)
	AuthPassword     string               `json:","` // Authentication protocol pass phrase (V3 specific)
	AuthProtocol     snmpgo.AuthProtocol  `json:","` // Authentication protocol (V3 specific)
	PrivPassword     string               `json:","` // Privacy protocol pass phrase (V3 specific)
	PrivProtocol     snmpgo.PrivProtocol  `json:","` // Privacy protocol (V3 specific)
	SecurityEngineId string               `json:","` // Security engine ID (V3 specific)
	ContextEngineId  string               `json:","` // Context engine ID (V3 specific)
	ContextName      string               `json:","` // Context name (V3 specific)
}

// MixWithAddress create snmpgo.SNMPArguments from mixing with a valid address
func (arg SnmpCredential) MixWithAddress(validAddress string) snmpgo.SNMPArguments {
	result := snmpgo.SNMPArguments{}
	res := reflect.ValueOf(&result).Elem()
	val := reflect.ValueOf(arg)
	valType := val.Type()

	for i := 0; i < val.NumField() && val.Field(i).CanInterface(); i++ {
		pName := valType.Field(i).Name
		pValue := val.Field(i)
		res.FieldByName(pName).Set(pValue)
	}
	res.FieldByName("Address").SetString(validAddress)

	return result
}
