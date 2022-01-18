package solus

// FilterOsImages represent available filters for fetching list of os images.
type FilterOsImages struct {
	filter
}

// ByName filter os images by specified name.
func (f *FilterOsImages) ByName(name string) *FilterOsImages {
	f.add("filter[search]", name)
	return f
}
