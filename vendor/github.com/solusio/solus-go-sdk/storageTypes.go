package solus

import (
	"context"
)

// StorageTypesService handles all available methods with storage types.
type StorageTypesService service

// StorageTypeName represents available storage types.
type StorageTypeName string

const (
	// StorageTypeNameFB indicates File Based storage.
	StorageTypeNameFB StorageTypeName = "fb"

	// StorageTypeNameLVM indicates LVM storage.
	StorageTypeNameLVM StorageTypeName = "lvm"

	// StorageTypeNameThinLVM indicates ThinLVM storage.
	StorageTypeNameThinLVM StorageTypeName = "thinlvm"

	// StorageTypeNameNFS indicates NFS storage.
	StorageTypeNameNFS StorageTypeName = "nfs"

	// StorageTypeNameVZ indicates VZ storage.
	StorageTypeNameVZ StorageTypeName = "vz"
)

// ImageFormat represents available image formats.
// Image formats is format of VM's disk.
type ImageFormat string

const (
	// ImageFormatQCOW2 indicates QCOW2 disk image.
	ImageFormatQCOW2 ImageFormat = "qcow2"

	// ImageFormatRaw indicates RAW disk image.
	ImageFormatRaw ImageFormat = "raw"

	// ImageFormatPLOOP indicates PLOOP disk image.
	ImageFormatPLOOP ImageFormat = "ploop"
)

// StorageType represents a storage type.
type StorageType struct {
	ID      int             `json:"id"`
	Name    StorageTypeName `json:"name"`
	Formats []ImageFormat   `json:"formats"`
}

// List lists storage types.
func (s *StorageTypesService) List(ctx context.Context) ([]StorageType, error) {
	var resp struct {
		Data []StorageType `json:"data"`
	}
	return resp.Data, s.client.get(ctx, "storage_types", &resp)
}
