package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func dataSourcePlan() *schema.Resource {
	return &schema.Resource{
		ReadContext: adoptRead("Plan", dataSourcePlanRead),
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "ID of the Plan",
				ValidateFunc: validation.IntAtLeast(1),
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Name of the Plan",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
		},
	}
}

func dataSourcePlanRead(ctx context.Context, client *client, d *schema.ResourceData) error {
	var (
		res solus.Plan
		err error
	)

	id, hasID := d.GetOk("id")
	name, hasName := d.GetOk("name")

	switch {
	case hasID:
		res, err = client.Plans.Get(ctx, id.(int))
		err = normalizeAPIError(err)

	case hasName:
		res, err = dataSourcePlanByName(ctx, client, name.(string))
	}

	if err != nil {
		return err
	}

	return newSchemaChainSetter(d).
		SetID(res.ID).
		Set("name", res.Name).
		Error()
}

func dataSourcePlanByName(ctx context.Context, client *client, name string) (solus.Plan, error) {
	res, err := client.Plans.List(ctx, new(solus.FilterPlans).ByName(name))
	if err != nil {
		return solus.Plan{}, err
	}

	if len(res.Data) == 1 {
		return res.Data[0], nil
	}

	err = errResourceNotFound
	if len(res.Data) > 1 {
		err = errTooManyResults
	}
	return solus.Plan{}, err
}
