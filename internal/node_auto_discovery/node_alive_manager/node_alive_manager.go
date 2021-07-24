package node_alive_manager

import (
	modelNode "../../model/node"
)

var nodes = make([]*modelNode.Node, 0)
var keys = make([]string, 0)

func NodeIndex(key string) int {
	index := -1
	for i, v := range keys {
		if v == key {
			index = i
			break
		}
	}
	return index
}

func NodeAdd(key string) *modelNode.Node {
	node := modelNode.GetNodeByKey(key)
	nodes = append(nodes, node)
	keys = append(keys, key)
	return node
}

func NodeDelete(key string) {
	index := NodeIndex(key)
	if index != -1 {
		nodes = append(nodes[:index], nodes[index+1:]...)
		keys = append(keys[:index], keys[index+1:]...)

	}
}

func NodeGetOne(key string) *modelNode.Node {
	index := NodeIndex(key)
	if index != -1 {
		return nodes[index]
	} else {
		return nil
	}
}

func NodeGetAll() []*modelNode.Node {
	return nodes
}

func NodeAliveCount() int {
	return len(nodes)
}

func NodeAliveClear() {
	nodes = make([]*modelNode.Node, 0)
}
