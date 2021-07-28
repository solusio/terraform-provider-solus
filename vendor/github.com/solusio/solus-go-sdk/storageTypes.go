package solus

import (
	"context"
)

type StorageTypesService service

type StorageTypeName string

const (
	StorageTypeNameFB      StorageTypeName = "fb"
	StorageTypeNameLVM     StorageTypeName = "lvm"
	StorageTypeNameThinLVM StorageTypeName = "thinlvm"
	StorageTypeNameNFS     StorageTypeName = "nfs"
)

type ImageFormat string

const (
	ImageFormatQCOW2 ImageFormat = "qcow2"
	ImageFormatRaw   ImageFormat = "raw"
)

type StorageType struct {
	ID      int             `json:"id"`
	Name    StorageTypeName `json:"name"`
	Formats []ImageFormat   `json:"formats"`
}

func (s *StorageTypesService) List(ctx context.Context) ([]StorageType, error) {
	var resp struct {
		Data []StorageType `json:"data"`
	}
	return resp.Data, s.client.get(ctx, "storage_types", &resp)
}
