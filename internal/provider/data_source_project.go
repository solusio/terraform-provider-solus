package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: adoptRead("Project", dataSourceProjectRead),
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "ID of the Project",
				ValidateFunc: validation.IntAtLeast(1),
			},
		},
	}
}

func dataSourceProjectRead(ctx context.Context, client *client, d *schema.ResourceData) error {
	var (
		res solus.Project
		err error
	)

	if id, hasID := d.GetOk("id"); hasID {
		res, err = client.Projects.Get(ctx, id.(int))
		err = normalizeAPIError(err)
	}

	if err != nil {
		return err
	}

	return newSchemaChainSetter(d).
		SetID(res.ID).
		Error()
}
