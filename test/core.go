package main

import (
	"../internal/core"
	mnc "../internal/manage_node_credentials"
	"fmt"
	"github.com/k-sone/snmpgo"
	"reflect"
	"time"
)

func testLog(epoch int) {
	t := time.Now()
	for i := 0; i < epoch; i++ {
		e := core.LogDebug(fmt.Sprintf("Message #%d", i+1))
		if e != nil {
			panic(e)
		}
	}
	t1 := time.Since(t).Seconds()
	fmt.Printf("coreLog records written: %d - time: %.9f s -> ops/sec: %.2f", epoch, t1, float64(epoch)/t1)
}

type Student struct {
	Age  int
	Name string
}
type Course struct {
	Name string
	Ta1  Student
	Ta2  Student
}

func testStructToJsonToStruct() {
	arg := mnc.SnmpCredential{
		Version:   snmpgo.V1,
		Community: "public",
		Retries:   5,
	}
	fmt.Printf("Struct -> Type: %s, Data: %s\n", reflect.TypeOf(arg), arg)

	a := core.StructToJson(arg)
	fmt.Printf("Json   -> Type: %s, Data: %s\n", reflect.TypeOf(a), a)

	snmpcred := mnc.SnmpCredential{}
	core.JsonToStruct(a, &snmpcred)
	fmt.Printf("Struct -> Type: %s, Data: %s\n", reflect.TypeOf(snmpcred), snmpcred)

	snmpargs := snmpcred.MixWithAddress("127.0.0.1:161")
	fmt.Printf("Struct -> Type: %s, Data: %s\n", reflect.TypeOf(snmpargs), snmpargs)

	fmt.Println()

	s1 := Student{
		Age:  29,
		Name: "ali",
	}
	s2 := Student{
		Age:  25,
		Name: "mohammad",
	}
	cr := Course{
		Name: "classroom-1",
		Ta1:  s1,
		Ta2:  s2,
	}
	fmt.Printf("Struct -> Type: %s, Data: %s\n", reflect.TypeOf(cr), cr)

	b := core.StructToJson(cr)
	fmt.Printf("Json   -> Type: %s, Data: %s\n", reflect.TypeOf(b), b)
}

func main() {
	testLog(10000)
	testStructToJsonToStruct()
}
