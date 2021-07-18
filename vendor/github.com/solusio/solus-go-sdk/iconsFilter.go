package solus

// FilterIcons represent available filters for fetching list of icons.
type FilterIcons struct {
	filter
}

// ByName filter icons by specified name.
func (f *FilterIcons) ByName(name string) *FilterIcons {
	f.add("filter[search]", name)
	return f
}

// ByType filter icons by specified type.
func (f *FilterIcons) ByType(t IconType) *FilterIcons {
	f.add("filter[type]", string(t))
	return f
}
