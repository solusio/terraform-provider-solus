package solus

// FilterLocations represent available filters for fetching list of Locations.
type FilterLocations struct {
	filter
}

// ByName filter Locations by specified name.
func (f *FilterLocations) ByName(name string) *FilterLocations {
	f.add("filter[search]", name)
	return f
}
