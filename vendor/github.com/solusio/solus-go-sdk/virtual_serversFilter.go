package solus

// FilterVirtualServers represent available filters for fetching list of virtual
// servers.
type FilterVirtualServers struct {
	filter
}

// ByUserID filter virtual servers by specified User ID.
func (f *FilterVirtualServers) ByUserID(id int) *FilterVirtualServers {
	f.addInt("filter[user_id]", id)
	return f
}

// ByComputeResourceID filter virtual servers by specified Compute Resource ID.
func (f *FilterVirtualServers) ByComputeResourceID(id int) *FilterVirtualServers {
	f.addInt("filter[compute_resource_id]", id)
	return f
}

// ByStatus filter virtual servers by specified status.
func (f *FilterVirtualServers) ByStatus(status string) *FilterVirtualServers {
	f.add("filter[status]", status)
	return f
}

// ByVirtualizationType filter virtual servers by virtualization type.
func (f *FilterVirtualServers) ByVirtualizationType(virtualizationType VirtualizationType) *FilterVirtualServers {
	f.add("filter[virtualization_type]", string(virtualizationType))
	return f
}
