package model_related

import (
	modelCredentials "../../../model/credentials"
	modelDiscovery "../../../model/discovery"
	modelNode "../../../model/node"
	"fmt"
	"reflect"
)

func GetStructFields(structName string) (fieldNames []string, fieldKinds []string) {
	fieldNames = make([]string, 0)
	fieldKinds = make([]string, 0)

	v := GetValue(structName)

	vType := v.Type()
	for i := 0; i < v.NumField() && v.Field(i).CanInterface(); i++ {
		fieldNames = append(fieldNames, fmt.Sprintf("%s", vType.Field(i).Name))
		fieldKinds = append(fieldKinds, fmt.Sprintf("%s", vType.Field(i).Type.Kind().String()))
	}
	return
}

func GetValue(structName string) reflect.Value {
	switch structName {
	case "SnmpCredential":
		return reflect.ValueOf(&modelCredentials.SnmpCredential{}).Elem()
	case "WindowsCredential":
		return reflect.ValueOf(&modelCredentials.WindowsCredential{}).Elem()
	case "LinuxCredential":
		return reflect.ValueOf(&modelCredentials.LinuxCredential{}).Elem()
	case "DeviceCredential":
		return reflect.ValueOf(&modelCredentials.DeviceCredential{}).Elem()
	case "Node":
		return reflect.ValueOf(&modelNode.Node{}).Elem()
	case "Discovery":
		return reflect.ValueOf(&modelDiscovery.Discovery{}).Elem()
	default:
		return reflect.Value{}
	}
}
