package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func resourceOSImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: adoptCreate("OS Image", resourceOSImageCreate),
		ReadContext:   adoptRead("OS Image", resourceOSImageRead),
		UpdateContext: adoptUpdate("OS Image", resourceOSImageUpdate),
		DeleteContext: adoptDelete("OS Image", resourceOSImageDelete),

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

func resourceOSImageCreate(ctx context.Context, client *solus.Client, d *schema.ResourceData) error {
	res, err := client.OsImages.Create(ctx, buildOSImageRequest(d))
	if err != nil {
		return normalizeAPIError(err)
	}

	d.SetId(strconv.Itoa(res.ID))
	return resourceOSImageRead(ctx, client, d)
}

func resourceOSImageRead(ctx context.Context, client *solus.Client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	res, err := client.OsImages.Get(ctx, id)
	if err != nil {
		return normalizeAPIError(err)
	}

	return newSchemaChainSetter(d).
		SetID(res.ID).
		Set("name", res.Name).
		Set("icon_id", res.Icon.ID).
		Set("is_visible", res.IsVisible).
		Error()
}

func resourceOSImageUpdate(ctx context.Context, client *solus.Client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	res, err := client.OsImages.Update(ctx, id, buildOSImageRequest(d))
	if err != nil {
		return normalizeAPIError(err)
	}

	d.SetId(strconv.Itoa(res.ID))
	return resourceOSImageRead(ctx, client, d)
}

func resourceOSImageDelete(ctx context.Context, client *solus.Client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	return normalizeAPIError(client.OsImages.Delete(ctx, id))
}

func buildOSImageRequest(d *schema.ResourceData) solus.OsImageRequest {
	return solus.OsImageRequest{
		Name:      d.Get("name").(string),
		IconID:    newNullableIntForID(d.Get("icon_id").(int)),
		IsVisible: d.Get("is_visible").(bool),
	}
}
