package xhtml

import (
	"../xdatabase"
	structPerModel "../xstruct/model_related"
	"fmt"
)

func CreateModalFields(structName string) string {
	r := ""
	fieldNames, _ := structPerModel.GetStructFields(structName)
	switch structName {
	case "SnmpCredential":
		for _, fn := range fieldNames {
			fid := "modalField_" + fn
			switch fn {
			case "Version":
				r += ModalSelectSingleField(fn, "", fid, [][]string{
					{"1", "0"},
					{"2c", "1"},
					{"3", "2"},
				}, 0)
			case "Network":
				r += ModalSelectSingleField(fn, "", fid, [][]string{
					{"UDP", "udp"},
					{"TCP", "tcp"},
				}, 0)
			case "Timeout":
				r += ModalNumberField(fn, "(ns) (def:5,000,000,000)", fid, 20000000, 10000000000, 10000000)
			case "Retries":
				r += ModalNumberField(fn, "(def:0)", fid, 0, 20, 1)
			case "MessageMaxSize":
				r += ModalNumberField(fn, "(def:1400)", fid, 50, 20000, 1)
			case "Community":
				r += ModalTextField(fn, "", fid)
			case "UserName":
				r += ModalTextField(fn, "", fid)
			case "SecurityLevel":
				r += ModalSelectSingleField(fn, "", fid, [][]string{
					{"No Auth, No Private", "0"},
					{"Auth, No Private", "1"},
					{"Auth, Private", "2"},
				}, 0)
			case "AuthPassword":
				r += ModalTextField(fn, "", fid)
			case "AuthProtocol":
				r += ModalSelectSingleField(fn, "", fid, [][]string{
					{"MD5", "MD5"},
					{"SHA", "SHA"},
				}, 0)
			case "PrivPassword":
				r += ModalTextField(fn, "", fid)
			case "PrivProtocol":
				r += ModalSelectSingleField(fn, "", fid, [][]string{
					{"DES", "DES"},
					{"AES", "AES"},
				}, 0)
			case "SecurityEngineId":
				r += ModalTextField(fn, "", fid)
			case "ContextEngineId":
				r += ModalTextField(fn, "", fid)
			case "ContextName":
				r += ModalTextField(fn, "", fid)
			}
		}
	case "WindowsCredential":
		for _, fn := range fieldNames {
			fid := "modalField_" + fn
			switch fn {
			case "Username":
				r += ModalTextField(fn, "", fid)
			case "Password":
				r += ModalTextField(fn, "", fid)
			case "AuthType":
				r += ModalSelectSingleField(fn, "", fid, [][]string{
					{"Local Auth", "0"},
					{"Domain Auth", "1"},
				}, 0)
			}
		}
	case "LinuxCredential":
		for _, fn := range fieldNames {
			fid := "modalField_" + fn
			switch fn {
			case "Username":
				r += ModalTextField(fn, "", fid)
			case "Password":
				r += ModalTextField(fn, "", fid)
			case "AuthType":
				r += ModalSelectSingleField(fn, "", fid, [][]string{
					{"Local Auth", "0"},
					{"Domain Auth", "1"},
				}, 0)
			}
		}
	case "DeviceCredential":
		for _, fn := range fieldNames {
			fid := "modalField_" + fn
			switch fn {
			case "Username":
				r += ModalTextField(fn, "", fid)
			case "Password":
				r += ModalTextField(fn, "", fid)
			case "AuthType":
				r += ModalSelectSingleField(fn, "", fid, [][]string{
					{"Auth SSH-1", "0"},
					{"Auth SSH-2", "1"},
					{"Auth Telnet", "2"},
				}, 0)
			}
		}
	case "Node":
		for _, fn := range fieldNames {
			fid := "modalField_" + fn
			switch fn {
			case "IpPort":
				r += ModalTextField(fn, "", fid)
			case "AliveType":
				r += ModalSelectSingleField(fn, "", fid, [][]string{
					{"ICMP", "0"},
					{"SNMP", "1"},
				}, 0)
			case "AliveInterval":
				r += ModalNumberField(fn, "(ns)", fid, 10000000000, 60000000000, 5000000000)
			case "PingPacketSize":
				r += ModalNumberField(fn, "(ns)", fid, 1, 65535, 1)
			case "PingTimeout":
				r += ModalNumberField(fn, "(ns)", fid, 20000000, 300000000, 10000000)
			}
		}
	case "Discovery":
		for _, fn := range fieldNames {
			fid := "modalField_" + fn
			switch fn {
			case "Ranges":
				r += ModalTextField(fn, "", fid)
			case "DiscoveryType":
				r += ModalSelectSingleField(fn, "", fid, [][]string{
					{"ICMP", "0"},
					{"SNMP", "1"},
				}, 0)
			case "SnmpCredentialKey":
				temp := make([][]string, 0)
				vals, keys, _ := xdatabase.ObjectGetAll("SnmpCredential")
				for i := range keys {
					cred := vals[i].(map[string]interface{})
					title := fmt.Sprintf("Ver.%.0f (%s) (Community: %s)", cred["Version"].(float64), cred["Network"], cred["Community"])
					temp = append(temp, []string{title, keys[i]})
				}
				r += ModalSelectSingleField(fn, "", fid, temp, 0)
			case "SnmpPort":
				r += ModalNumberField(fn, "", fid, 1, 65535, 1)
			case "PingPacketSize":
				r += ModalNumberField(fn, "(bytes)", fid, 1, 65535, 1)
			}
		}
	}
	return r
}

func CreateRequestHeaders(structName string) string {
	r := ""
	r += fmt.Sprintf("xhttp.setRequestHeader('StructName', '%s')\n", structName)
	fieldNames, _ := structPerModel.GetStructFields(structName)
	switch structName {
	case "Discovery":
		fNames := []string{"Ranges", "DiscoveryType", "SnmpCredentialKey", "SnmpPort", "PingPacketSize"}
		for _, fn := range fNames {
			fid := "modalField_" + fn
			r += fmt.Sprintf("xhttp.setRequestHeader('%s', document.getElementById('%s').value)\n", fid, fid)
		}
	default:
		for _, fn := range fieldNames {
			fid := "modalField_" + fn
			r += fmt.Sprintf("xhttp.setRequestHeader('%s', document.getElementById('%s').value)\n", fid, fid)
		}
	}

	return r
}

func CreateColumnDefs(structName string) string {
	r := ""
	fieldNames, _ := structPerModel.GetStructFields(structName)
	for i, fn := range fieldNames {
		r += fmt.Sprintf("{title: '%s', data: '%s', targets: %d, adjust: true},\n", fn, fn, i+1)
	}

	switch structName {
	case "Discovery":
		r += fmt.Sprintf("{title: 'Start', targets: %d, adjust: true, defaultContent: '<a class=\"btn btn-primary btn-sm d-none d-sm-inline-block\" role=\"button\" onclick=\"discoverButton(this)\"><i class=\"fas fa-download fa-sm text-white-50\"></i>Start</a>'},\n", len(fieldNames)+1)
	case "Node":
		r += fmt.Sprintf("{title: 'Start', targets: %d, adjust: true, defaultContent: '<a class=\"btn btn-primary btn-sm d-none d-sm-inline-block\" role=\"button\" onclick=\"nodeButton(this)\"><i class=\"fas fa-download fa-sm text-white-50\"></i>Start</a>'},\n", len(fieldNames)+1)
	}

	return r
}
