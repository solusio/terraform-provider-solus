package solus

// FilterComputeResources represent available filters for fetching list of compute
// resources.
type FilterComputeResources struct {
	filter
}

// ByName filter compute resources by specified name.
func (f *FilterComputeResources) ByName(name string) *FilterComputeResources {
	f.add("filter[search]", name)
	return f
}
