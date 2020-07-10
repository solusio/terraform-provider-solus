package solus

type FilterUsers struct {
	filter
}

// ByStatus filter Users by specified status.
func (f *FilterUsers) ByStatus(status string) *FilterUsers {
	f.add("filter[status]", status)
	return f
}
