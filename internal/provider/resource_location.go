package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func resourceLocation() *schema.Resource {
	return &schema.Resource{
		CreateContext: adoptCreate("Location", resourceLocationCreate),
		ReadContext:   adoptRead("Location", resourceLocationRead),
		UpdateContext: adoptUpdate("Location", resourceLocationUpdate),
		DeleteContext: adoptDelete("Location", resourceLocationDelete),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"icon_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"is_default": {
				Type:         schema.TypeBool,
				Optional:     true,
				Default:      false,
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

func resourceLocationCreate(ctx context.Context, client *solus.Client, d *schema.ResourceData) error {
	res, err := client.Locations.Create(ctx, buildLocationRequest(d))
	if err != nil {
		return normalizeAPIError(err)
	}

	d.SetId(strconv.Itoa(res.ID))
	return resourceLocationRead(ctx, client, d)
}

func resourceLocationRead(ctx context.Context, client *solus.Client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	res, err := client.Locations.Get(ctx, id)
	if err != nil {
		return normalizeAPIError(err)
	}

	return newSchemaChainSetter(d).
		Set("name", res.Name).
		Set("description", res.Description).
		Set("icon_id", res.Icon.ID).
		Set("is_default", res.IsDefault).
		Set("is_visible", res.IsVisible).
		Error()
}

func resourceLocationUpdate(ctx context.Context, client *solus.Client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	res, err := client.Locations.Update(ctx, id, buildLocationRequest(d))
	if err != nil {
		return normalizeAPIError(err)
	}

	d.SetId(strconv.Itoa(res.ID))
	return resourceLocationRead(ctx, client, d)
}

func resourceLocationDelete(ctx context.Context, client *solus.Client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	return normalizeAPIError(client.Locations.Delete(ctx, id))
}

func buildLocationRequest(d *schema.ResourceData) solus.LocationCreateRequest {
	return solus.LocationCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		IconID:      newNullableIntForID(d.Get("icon_id").(int)),
		IsDefault:   d.Get("is_default").(bool),
		IsVisible:   d.Get("is_visible").(bool),
	}
}
