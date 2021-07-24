package model_independent

import (
	"fmt"
	"reflect"
)

func GetStructNameOfInterface(v interface{}) string {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func SetStructField(itemPointer interface{}, fieldName string, newValue interface{}) error {
	structValue := reflect.ValueOf(itemPointer).Elem()
	//fmt.Printf("value -> %s (iszero %t) (isvalid %t)\n", fieldName, structValue.IsZero(), structValue.IsValid())
	structFieldValue := structValue.FieldByName(fieldName)
	//fmt.Printf("field name -> %s (iszero %t) (isvalid %t)\n", fieldName, structFieldValue.IsZero(), structFieldValue.IsValid())

	if !structFieldValue.IsValid() {
		fmt.Printf("no such field (%s)\n", fieldName)
		return nil
		//return errors.New(fmt.Sprintf("no such field (%s)", fieldName))
	}

	if !structFieldValue.CanSet() {
		fmt.Printf("cannot set field (%s)\n", fieldName)
		return nil
		//return errors.New(fmt.Sprintf("cannot set field (%s)", fieldName))
	}
	structFieldType := structFieldValue.Type()
	structFieldKind := structFieldType.Kind()
	newVal := reflect.ValueOf(newValue)
	//if structFieldType != newVal.Type() {
	//	fmt.Printf("types not match for field (%s)\n", fieldName)
	//	//return errors.New(fmt.Sprintf("types not match for field (%s)", fieldName))
	//}
	if structFieldKind == reflect.Int ||
		structFieldKind == reflect.Int32 ||
		structFieldKind == reflect.Int64 {
		structFieldValue.SetInt(int64(newVal.Float()))
	} else if structFieldKind == reflect.Uint ||
		structFieldKind == reflect.Uint32 ||
		structFieldKind == reflect.Uint64 {
		structFieldValue.SetUint(uint64(newVal.Float()))
	} else if structFieldKind == reflect.String {
		structFieldValue.SetString(newVal.String())
	} else {
		structFieldValue.Set(newVal)
	}
	return nil
}
