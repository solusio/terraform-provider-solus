package solus

// FilterTasks represent available filters for fetching list of tasks.
type FilterTasks struct {
	filter
}

// ByAction filter Tasks by specified action.
func (f *FilterTasks) ByAction(action string) *FilterTasks {
	f.add("filter[action]", action)
	return f
}

// ByStatus filter Tasks by specified status.
func (f *FilterTasks) ByStatus(status string) *FilterTasks {
	f.add("filter[status]", status)
	return f
}

// ByComputeResourceID filter Tasks by specified Compute Resource ID.
func (f *FilterTasks) ByComputeResourceID(id int) *FilterTasks {
	f.addInt("filter[compute_resource_id]", id)
	return f
}

// ByComputeResourceVMID filter Tasks by specified Compute Resource VM ID.
func (f *FilterTasks) ByComputeResourceVMID(id int) *FilterTasks {
	f.addInt("filter[compute_resource_vm_id]", id)
	return f
}
