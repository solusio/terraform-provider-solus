package solus

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func dataSourceOsImageVersion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOsImageVersionRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "id of the os image version",
				ValidateFunc:  validation.NoZeroValues,
				ConflictsWith: []string{"os_image_id", "version"},
			},
			"os_image_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "os image id",
				ValidateFunc:  validation.NoZeroValues,
				ConflictsWith: []string{"id"},
				RequiredWith:  []string{"version"},
			},
			"version": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "os image version name",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id"},
				RequiredWith: []string{"os_image_id"},
			},
		},
	}
}

func dataSourceOsImageVersionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*solus.Client)

	var (
		res solus.OsImageVersion
		err error
	)

	_, hasID := d.GetOk("id")
	_, hasVersion := d.GetOk("version")

	// These two properties are mutually exclusive.
	switch {
	case hasID:
		res, err = dataSourceOsImageVersionReadByID(ctx, client, d)

	case hasVersion:
		res, err = dataSourceOsImageVersionReadByVersion(ctx, client, d)
	}

	if err != nil {
		return diag.Errorf("failed to get os image version: %s", err)
	}

	if res.ID == 0 {
		return diag.Errorf("os image version not found")
	}

	err = (&schemaChainSetter{d: d}).
		SetID(res.ID).
		Set("version", res.Version).
		Error()
	if err != nil {
		return diag.Errorf("failed to map os image version response to data: %s", err)
	}

	return nil
}

func dataSourceOsImageVersionReadByID(
	ctx context.Context,
	c *solus.Client,
	d *schema.ResourceData,
) (solus.OsImageVersion, error) {
	return c.OsImageVersions.Get(ctx, d.Get("id").(int))
}

func dataSourceOsImageVersionReadByVersion(
	ctx context.Context,
	c *solus.Client,
	d *schema.ResourceData,
) (solus.OsImageVersion, error) {
	osImageID := d.Get("os_image_id").(int)
	version := d.Get("version").(string)

	vv, err := c.OsImages.ListVersion(ctx, osImageID)
	if err != nil {
		return solus.OsImageVersion{}, err
	}

	var res solus.OsImageVersion
	for _, v := range vv {
		if v.Version == version {
			res = v
			break
		}
	}
	return res, nil
}
