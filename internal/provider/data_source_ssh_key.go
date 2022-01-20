package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func dataSourceSSHKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: adoptRead("SSH Key", dataSourceSSHKeyRead),
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "ID of the SSH Key",
				ValidateFunc: validation.IntAtLeast(1),
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Name of the SSH Key",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
		},
	}
}

func dataSourceSSHKeyRead(ctx context.Context, client *client, d *schema.ResourceData) error {
	var (
		res solus.SSHKey
		err error
	)

	id, hasID := d.GetOk("id")
	name, hasName := d.GetOk("name")

	switch {
	case hasID:
		res, err = client.SSHKeys.Get(ctx, id.(int))
		err = normalizeAPIError(err)

	case hasName:
		res, err = dataSourceSSHKeyByName(ctx, client, name.(string))
	}

	if err != nil {
		return err
	}

	return newSchemaChainSetter(d).
		SetID(res.ID).
		Set("name", res.Name).
		Error()
}

func dataSourceSSHKeyByName(ctx context.Context, client *client, name string) (solus.SSHKey, error) {
	res, err := client.SSHKeys.List(ctx, new(solus.FilterSSHKeys).ByName(name))
	if err != nil {
		return solus.SSHKey{}, err
	}

	if len(res.Data) == 1 {
		return res.Data[0], nil
	}

	err = errResourceNotFound
	if len(res.Data) > 1 {
		err = errTooManyResults
	}
	return solus.SSHKey{}, err
}
