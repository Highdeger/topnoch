package xlist

//ItemIndex finds the index of an item in a slice.
func ItemIndex(item string, list *[]string) int {
	index := -1
	for i, v := range *list {
		if v == item {
			index = i
			break
		}
	}
	return index
}
