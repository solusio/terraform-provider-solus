package solus

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func dataSourceIcon() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIconRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "id of the icon",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "name of the icon",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
		},
	}
}

func dataSourceIconRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	metadata, ok := m.(metadata)
	if !ok {
		return diag.Errorf("invalid metadata type %T", m)
	}
	client := metadata.Client
	timeout := metadata.RequestTimeout

	var (
		res solus.Icon
		err error
	)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	rawID, hasRawID := d.GetOk("id")
	rawName, hasRawName := d.GetOk("name")

	// These two properties are mutually exclusive.
	switch {
	case hasRawID:
		id, ok := rawID.(int)
		if !ok {
			return diag.Errorf("id isn't an integer")
		}

		res, err = client.Icons.Get(ctx, id)
		if err != nil {
			return diag.Errorf("failed to get icon by id %d: %s", id, err)
		}

	case hasRawName:
		name, ok := rawName.(string)
		if !ok {
			return diag.Errorf("name isn't a string")
		}

		p, err := client.Icons.List(ctx, new(solus.FilterIcons).ByName(name))
		if err != nil {
			return diag.Errorf("failed to get icon by name %q: %s", name, err)
		}

		if len(p.Data) == 0 {
			return diag.Errorf("icon not found")
		}

		if len(p.Data) > 1 {
			return diag.Errorf("find more than one icon")
		}

		res = p.Data[0]
	}

	d.SetId(strconv.Itoa(res.ID))

	return nil
}
