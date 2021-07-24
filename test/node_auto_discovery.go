package main

import (
	nad "../internal/node_auto_discovery"
	"fmt"
)

func testIpRange() {
	input := "192.168.4.0-193.170.7.3"
	fmt.Printf("\nparsing the range '%s'\n\n", input)
	r, e := nad.IpListFromIpRange(input)
	if e != nil {
		panic(e)
	}
	printList(r)
}

func testIpSlash() {
	for _, v := range []string{"55.187.224.56/22", "34.260.199.11/20", "34.71.199.11/35", "34.71.199.11/20"} {
		fmt.Printf("\nparsing the cidr (ip/mask) '%s'\n", v)
		r, err := nad.IpListFromIpSlash(v)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			printList(r)
		}
	}
}

func testIpSubnet() {
	fmt.Printf("\nparsing the cidr (ip/subnet) '%s/%s'\n", "192.168.75.25", "255.255.192.0")
	r, err := nad.IpListFromIpSubnet("192.168.75.25", "255.255.192.0")
	if err != nil {
		fmt.Println(err)
	} else {
		printList(r)
	}
}

func printList(r []string) {
	fmt.Printf("found ip: %d\n", len(r))
	fmt.Println(r[0])
	fmt.Println(r[1])
	fmt.Println(r[2])
	fmt.Println("...")
	fmt.Println(r[len(r)-3])
	fmt.Println(r[len(r)-2])
	fmt.Println(r[len(r)-1])
}

func main() {
	testIpRange()
	testIpSlash()
	testIpSubnet()
}
