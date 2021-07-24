package main

import (
	"../internal/parse"
	"fmt"
	"strings"
)

func testTokenizing() {
	r := parse.TokenizeText("Linux devdeb 4.19.0-13-amd64 #1 SMP Debian 4.19.160-2 (2020-11-28) x86_64", []string{" "})
	out := "{"
	for i, part := range r {
		if i == len(r)-1 {
			out += "\"" + strings.Trim(part, " ") + "\""
		} else {
			out += "\"" + strings.Trim(part, " ") + "\"" + ", "
		}
	}
	out += fmt.Sprintf("}(len=%d)", len(r))
	fmt.Println(out)
}

func main() {
	testTokenizing()
}
