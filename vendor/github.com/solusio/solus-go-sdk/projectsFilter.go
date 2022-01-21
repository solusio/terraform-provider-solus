package solus

// FilterProjects represent available filters for fetching list of Projects.
type FilterProjects struct {
	filter
}

// ByName baseFilter Projects by specified name.
func (f *FilterProjects) ByName(name string) *FilterProjects {
	f.add("filter[search]", name)
	return f
}
