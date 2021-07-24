package solus

// FilterIPBlocks represent available filters for fetching list of ip blocks.
type FilterIPBlocks struct {
	filter
}

// ByName filter ip blocks by specified name.
func (f *FilterIPBlocks) ByName(name string) *FilterIPBlocks {
	f.add("filter[search]", name)
	return f
}
