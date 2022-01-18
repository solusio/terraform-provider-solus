package solus

// FilterServers represent available filters for fetching list of Servers.
type FilterServers struct {
	filter
}

// ByUserID filter Servers by specified User ID.
func (f *FilterServers) ByUserID(id int) *FilterServers {
	f.addInt("filter[user_id]", id)
	return f
}

// ByComputeResourceID filter Servers by specified Compute Resource ID.
func (f *FilterServers) ByComputeResourceID(id int) *FilterServers {
	f.addInt("filter[compute_resource_id]", id)
	return f
}

// ByStatus filter Servers by specified status.
func (f *FilterServers) ByStatus(status string) *FilterServers {
	f.add("filter[status]", status)
	return f
}

// ByVirtualizationType filter Servers by virtualization type.
func (f *FilterServers) ByVirtualizationType(virtualizationType VirtualizationType) *FilterServers {
	f.add("filter[virtualization_type]", string(virtualizationType))
	return f
}
