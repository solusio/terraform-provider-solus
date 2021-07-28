package solus

import "strconv"

// FilterPlans represent available filters for fetching list of plans.
type FilterPlans struct {
	filter
}

// ByStorageType filter plans by specified storage type.
func (f *FilterPlans) ByStorageType(t StorageTypeName) *FilterPlans {
	f.add("filter[storage_type]", string(t))
	return f
}

// ByImageFormat filter plans by specified image format.
func (f *FilterPlans) ByImageFormat(v ImageFormat) *FilterPlans {
	f.add("filter[image_format]", string(v))
	return f
}

// ByName filter plans by specified name.
func (f *FilterPlans) ByName(name string) *FilterPlans {
	f.add("filter[search]", name)
	return f
}

// ByDiskSize filter plans by specified disk size.
func (f *FilterPlans) ByDiskSize(v int) *FilterPlans {
	f.add("filter[disk]", strconv.Itoa(v))
	return f
}
