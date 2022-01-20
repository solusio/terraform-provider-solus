package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: adoptCreate("Project", resourceProjectCreate),
		ReadContext:   adoptRead("Project", resourceProjectRead),
		UpdateContext: adoptUpdate("Project", resourceProjectUpdate),
		DeleteContext: adoptDelete("Project", resourceProjectDelete),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func resourceProjectCreate(ctx context.Context, client *client, d *schema.ResourceData) error {
	res, err := client.Projects.Create(ctx, buildProjectRequest(d))
	if err != nil {
		return normalizeAPIError(err)
	}

	d.SetId(strconv.Itoa(res.ID))
	return resourceProjectRead(ctx, client, d)
}

func resourceProjectRead(ctx context.Context, client *client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	res, err := client.Projects.Get(ctx, id)
	if err != nil {
		return normalizeAPIError(err)
	}

	return newSchemaChainSetter(d).
		Set("name", res.Name).
		Set("description", res.Description).
		Error()
}

func resourceProjectUpdate(ctx context.Context, client *client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	res, err := client.Projects.Update(ctx, id, buildProjectRequest(d))
	if err != nil {
		return normalizeAPIError(err)
	}

	d.SetId(strconv.Itoa(res.ID))
	return resourceProjectRead(ctx, client, d)
}

func resourceProjectDelete(ctx context.Context, client *client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	return normalizeAPIError(client.Projects.Delete(ctx, id))
}

func buildProjectRequest(d *schema.ResourceData) solus.ProjectRequest {
	return solus.ProjectRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
}
