package solus

// FilterSSHKeys represent available filters for fetching list of SSH keys.
type FilterSSHKeys struct {
	filter
}

// ByName filter SSH keys by specified name.
func (f *FilterSSHKeys) ByName(name string) *FilterSSHKeys {
	f.add("filter[search]", name)
	return f
}
