package solus

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/solusio/solus-go-sdk"
	"strconv"
)

func resourceLocation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLocationCreate,
		ReadContext:   resourceLocationRead,
		UpdateContext: resourceLocationUpdate,
		DeleteContext: resourceLocationDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"is_visible": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceLocationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	isDefault := d.Get("is_default").(bool)
	isVisible := d.Get("is_visible").(bool)

	client := m.(metadata).Client
	timeout := m.(metadata).RequestTimeout

	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	l, err := client.Locations.Create(reqCtx, solus.LocationCreateRequest{
		Name:             name,
		Description:      description,
		ComputeResources: []int{},
		IsDefault:        isDefault,
		IsVisible:        isVisible,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(l.ID))
	return resourceLocationRead(ctx, d, m)
}

func resourceLocationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(metadata).Client
	timeout := m.(metadata).RequestTimeout

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	l, err := client.Locations.Get(reqCtx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := locationToResourceData(l, d); err != nil {
		return diagErr("Failed to map location response to resource", err.Error())
	}

	return nil
}

func resourceLocationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(metadata).Client
	timeout := m.(metadata).RequestTimeout

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	isDefault := d.Get("is_default").(bool)
	isVisible := d.Get("is_visible").(bool)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	l, err := client.Locations.Update(reqCtx, id, solus.LocationCreateRequest{
		Name:             name,
		Description:      description,
		ComputeResources: []int{},
		IsDefault:        isDefault,
		IsVisible:        isVisible,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(l.ID))
	return resourceLocationRead(ctx, d, m)
}

func resourceLocationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(metadata).Client
	timeout := m.(metadata).RequestTimeout

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err = client.Locations.Delete(reqCtx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func locationToResourceData(l solus.Location, d *schema.ResourceData) error {
	return (&schemaChainSetter{d: d}).
		Set("id", l.ID).
		Set("name", l.Name).
		Set("icon", l.Icon).
		Set("description", l.Description).
		Set("is_default", l.IsDefault).
		Set("is_visible", l.IsVisible).
		Set("compute_resources", l.ComputeResources).
		Error()
}
