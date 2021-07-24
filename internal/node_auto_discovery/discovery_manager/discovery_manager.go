package node_alive_manager

import (
	modelDiscovery "../../model/discovery"
)

var discoveries = make([]*modelDiscovery.Discovery, 0)
var keys = make([]string, 0)

func DiscoveryIndex(key string) int {
	index := -1
	for i, v := range keys {
		if v == key {
			index = i
			break
		}
	}
	return index
}

func DiscoveryAdd(key string) {
	d := modelDiscovery.GetDiscoveryByKey(key)
	discoveries = append(discoveries, d)
	keys = append(keys, key)
}

func DiscoveryDelete(key string) {
	index := DiscoveryIndex(key)
	if index != -1 {
		discoveries = append(discoveries[:index], discoveries[index+1:]...)
		keys = append(keys[:index], keys[index+1:]...)
	}
}

func DiscoveryGetOne(key string) *modelDiscovery.Discovery {
	index := DiscoveryIndex(key)
	if index != -1 {
		return discoveries[index]
	} else {
		return nil
	}
}

func DiscoveryGetAll() []*modelDiscovery.Discovery {
	return discoveries
}

func DiscoveryCount() int {
	return len(discoveries)
}

func DiscoveryClear() {
	discoveries = make([]*modelDiscovery.Discovery, 0)
}
