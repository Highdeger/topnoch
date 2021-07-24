package main

import (
	"../internal/core/xcache"
	model "../internal/model/node"
	"fmt"
	"time"
)

func main() {
	total := 100000
	counter := 0

	fmt.Printf("Start to set %d of items [key->struct(interface)]\n", total)
	t := time.Now()
	for i := 0; i < total; i++ {
		err := xcache.Set(fmt.Sprintf("Key-%d", i), model.Node{
			Name: fmt.Sprintf("Name-%d", i)})
		if err == nil {
			// nothing
		} else {
			// nothing
		}
	}
	fmt.Printf("Setting ended in %.3fs\n", time.Since(t).Seconds())

	fmt.Printf("Start to get %d of items [key->struct(interface)]\n", total)
	t = time.Now()
	for i := 0; i < total; i++ {
		v, e := xcache.Get(fmt.Sprintf("Key-%d", i))
		if e == nil {
			if v != nil {
				counter++
			} else {
				// nothing
			}
		} else {
			// nothing
		}
	}

	fmt.Printf("Getting ended in %.3fs (success=%.2f%%)\n", time.Since(t).Seconds(), (float64(counter)/float64(total))*100)
}
