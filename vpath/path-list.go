package vpath

// VirtualPathList ...
type VirtualPathList []VirtualPath

// Len is the number of elements in the collection.
func (list VirtualPathList) Len() int {
	return len(list)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (list VirtualPathList) Less(i, j int) bool {
	return list[i].Path < list[j].Path
}

// Swap swaps the elements with indexes i and j.
func (list VirtualPathList) Swap(i, j int) {
	tmp := list[i]
	list[i] = list[j]
	list[j] = tmp
}
