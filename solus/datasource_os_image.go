package solus

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func dataSourceOsImage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOsImageRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "id of the os image",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "name of the os image",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
		},
	}
}

func dataSourceOsImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*solus.Client)

	var (
		res solus.OsImage
		err error
	)

	id, hasID := d.GetOk("id")
	name, hasName := d.GetOk("name")

	// These two properties are mutually exclusive.
	switch {
	case hasID:
		res, err = client.OsImages.Get(ctx, id.(int))
		if err != nil {
			return diag.Errorf("failed to get os image by id %d: %s", id, err)
		}

	case hasName:
		p, err := client.OsImages.List(ctx, new(solus.FilterOsImages).ByName(name.(string)))
		if err != nil {
			return diag.Errorf("failed to get os image by name %q: %s", name, err)
		}

		if len(p.Data) == 0 {
			return diag.Errorf("os image not found")
		}

		if len(p.Data) > 1 {
			return diag.Errorf("find more than one os image")
		}

		res = p.Data[0]
	}

	err = (&schemaChainSetter{d: d}).
		SetID(res.ID).
		Set("name", res.Name).
		Error()
	if err != nil {
		return diag.Errorf("failed to map os image response to data: %s", err)
	}

	return nil
}
