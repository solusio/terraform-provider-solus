package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func dataSourceOsImageVersion() *schema.Resource {
	return &schema.Resource{
		ReadContext: adoptRead("OS Image Version", dataSourceOsImageVersionRead),
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "ID of the OS Image Version",
				ValidateFunc:  validation.IntAtLeast(1),
				ConflictsWith: []string{"os_image_id", "version"},
			},
			"os_image_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "ID of the OS Image",
				ValidateFunc:  validation.IntAtLeast(1),
				ConflictsWith: []string{"id"},
				RequiredWith:  []string{"version"},
			},
			"version": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "OS Image Version name",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id"},
				RequiredWith: []string{"os_image_id"},
			},
		},
	}
}

func dataSourceOsImageVersionRead(ctx context.Context, client *client, d *schema.ResourceData) error {
	var (
		res solus.OsImageVersion
		err error
	)

	id, hasID := d.GetOk("id")
	_, hasVersion := d.GetOk("version")

	switch {
	case hasID:
		res, err = client.OsImageVersions.Get(ctx, id.(int))
		err = normalizeAPIError(err)

	case hasVersion:
		res, err = dataSourceOsImageVersionReadByVersion(ctx, client, d)
	}

	if err != nil {
		return err
	}

	return newSchemaChainSetter(d).
		SetID(res.ID).
		Set("version", res.Version).
		Error()
}

func dataSourceOsImageVersionReadByVersion(
	ctx context.Context,
	c *client,
	d *schema.ResourceData,
) (solus.OsImageVersion, error) {
	osImageID := d.Get("os_image_id").(int)
	version := d.Get("version").(string)

	vv, err := c.OsImages.ListVersion(ctx, osImageID)
	if err != nil {
		return solus.OsImageVersion{}, normalizeAPIError(err)
	}

	var res solus.OsImageVersion
	for _, v := range vv {
		if v.Version == version {
			res = v
			break
		}
	}

	if res.ID == 0 {
		return solus.OsImageVersion{}, errResourceNotFound
	}
	return res, nil
}
