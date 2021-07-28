package solus

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func dataSourcePlan() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePlanRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "id of the plan",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "name of the plan",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
		},
	}
}

func dataSourcePlanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*solus.Client)

	var (
		res solus.Plan
		err error
	)

	id, hasID := d.GetOk("id")
	name, hasName := d.GetOk("name")

	// These two properties are mutually exclusive.
	switch {
	case hasID:
		res, err = client.Plans.Get(ctx, id.(int))
		if err != nil {
			return diag.Errorf("failed to get plan by id %d: %s", id, err)
		}

	case hasName:
		p, err := client.Plans.List(ctx, new(solus.FilterPlans).ByName(name.(string)))
		if err != nil {
			return diag.Errorf("failed to get plan by name %q: %s", name, err)
		}

		if len(p.Data) == 0 {
			return diag.Errorf("plan not found")
		}

		if len(p.Data) > 1 {
			return diag.Errorf("find more than one plan")
		}

		res = p.Data[0]
	}

	err = (&schemaChainSetter{d: d}).
		SetID(res.ID).
		Set("name", res.Name).
		Error()
	if err != nil {
		return diag.Errorf("failed to map plan response to data: %s", err)
	}

	return nil
}
