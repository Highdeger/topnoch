package main

import (
	"../internal/monitor"
	"fmt"
	"github.com/k-sone/snmpgo"
)

var tesingSnmpArgs = snmpgo.SNMPArguments{
	Version:   snmpgo.V2c,
	Address:   "127.0.0.1:161",
	Community: "public",
}

func testSysDescr() {
	fmt.Printf("Description: %s\n", monitor.FetchSysDescr(tesingSnmpArgs))
}

func testSysObjectId() {
	fmt.Printf("System's OID: %s\n", monitor.FetchSysObjectId(tesingSnmpArgs))
}

func testSysUptime() {
	fmt.Printf("Uptime: %.2f s\n", float32(monitor.FetchSysUptime(tesingSnmpArgs))/100)
}

func testSysContactNameLocation() {
	contact, name, location := monitor.FetchSysContactNameLocation(tesingSnmpArgs)
	fmt.Printf("Contact: %s\nName: %s\nLocation: %s\n", contact, name, location)
}

func testSysServices() {
	fmt.Print("Layers potentially offering service: ")
	for _, v := range monitor.FetchSysServices(tesingSnmpArgs) {
		fmt.Print(monitor.OsiLayers[v], " ")
	}
	fmt.Println()
}

func testGetTable() {
	r := monitor.GetTableForOid("1.3.6.1.2.1.1.9", 1, []int{2, 3, 4}, tesingSnmpArgs)
	fmt.Println("Ind.\tColumns...")
	for _, row := range r {
		for _, v := range row {
			fmt.Printf("\t%s", v.GetValue())
		}
		fmt.Println()
	}
}

func testWalkOids() {
	r := monitor.WalkValuesForOids([]string{"1.3.6.1.4.1.2021.10"}, tesingSnmpArgs)
	for _, v := range r {
		fmt.Printf("%s\t\t%s\n", v.GetOid(), v.GetValue())
	}
}

func testAverageCpuLoadTable() {
	r := monitor.FetchAverageCpuLoadTable(tesingSnmpArgs)
	for _, row := range r {
		for _, v := range row {
			fmt.Printf("%s\t", v.GetValue())
		}
		fmt.Println()
	}
}

func testAverageCpuLoads() {
	for _, v := range monitor.FetchAverageCpuLoads(tesingSnmpArgs) {
		fmt.Println(v)
	}
}

func testSnmpCapabilities() {
	r := monitor.FetchSnmpCapabilities(tesingSnmpArgs)
	for _, row := range r {
		for _, v := range row {
			fmt.Printf("%s\t", v.GetValue())
		}
		fmt.Println()
	}
}

func main() {
	testSysDescr()
	testSysObjectId()
	testSysUptime()
	testSysContactNameLocation()
	testSysServices()
	testGetTable()
	testWalkOids()
	testAverageCpuLoadTable()
	testAverageCpuLoads()
	testSnmpCapabilities()
}
