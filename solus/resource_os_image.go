package solus

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func resourceOsImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOsImageCreate,
		ReadContext:   resourceOsImageRead,
		UpdateContext: resourceOsImageUpdate,
		DeleteContext: resourceOsImageDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"icon_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"is_visible": {
				Type:         schema.TypeBool,
				Optional:     true,
				Default:      true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func resourceOsImageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, ok := m.(*solus.Client)
	if !ok {
		return diag.Errorf("invalid Solus client type %T", m)
	}

	name := d.Get("name").(string)
	iconID := d.Get("icon_id").(int)
	isVisible := d.Get("is_visible").(bool)

	i, err := client.OsImages.Create(ctx, solus.OsImageRequest{
		Name:      name,
		IconID:    newNullableIntForID(iconID),
		IsVisible: isVisible,
	})
	if err != nil {
		return diag.Errorf("failed to create new os image: %s", err)
	}

	d.SetId(strconv.Itoa(i.ID))
	return resourceOsImageRead(ctx, d, m)
}

func resourceOsImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, ok := m.(*solus.Client)
	if !ok {
		return diag.Errorf("invalid Solus client type %T", m)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	i, err := client.OsImages.Get(ctx, id)
	if err != nil {
		return diag.Errorf("failed to get os image by id %d: %s", id, err)
	}

	err = (&schemaChainSetter{d: d}).
		SetID(i.ID).
		Set("name", i.Name).
		Set("icon_id", i.Icon.ID).
		Set("is_visible", i.IsVisible).
		Error()
	if err != nil {
		return diag.Errorf("failed to map os image response to resource: %s", err)
	}

	return nil
}

func resourceOsImageUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, ok := m.(*solus.Client)
	if !ok {
		return diag.Errorf("invalid Solus client type %T", m)
	}

	name := d.Get("name").(string)
	iconID := d.Get("icon_id").(int)
	isVisible := d.Get("is_visible").(bool)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	i, err := client.OsImages.Update(ctx, id, solus.OsImageRequest{
		Name:      name,
		IconID:    newNullableIntForID(iconID),
		IsVisible: isVisible,
	})
	if err != nil {
		return diag.Errorf("failed to update os image with id %d: %s", id, err)
	}

	d.SetId(strconv.Itoa(i.ID))
	return resourceOsImageRead(ctx, d, m)
}

func resourceOsImageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, ok := m.(*solus.Client)
	if !ok {
		return diag.Errorf("invalid Solus client type %T", m)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.OsImages.Delete(ctx, id)
	if err != nil {
		return diag.Errorf("failed to delete os image by id %d: %s", id, err)
	}
	return nil
}
