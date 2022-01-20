package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/guregu/null.v4"
)

type schemaChainSetter struct {
	d   *schema.ResourceData
	err error
}

func newSchemaChainSetter(d *schema.ResourceData) *schemaChainSetter {
	return &schemaChainSetter{
		d: d,
	}
}

func (s *schemaChainSetter) SetID(v int) *schemaChainSetter {
	if s.err == nil {
		s.d.SetId(strconv.Itoa(v))
	}
	return s
}

func (s *schemaChainSetter) Set(k string, v interface{}) *schemaChainSetter {
	if s.err != nil {
		return s
	}

	if err := s.d.Set(k, v); err != nil {
		s.err = fmt.Errorf("failed to set value for %q key: %w", k, err)
	}
	return s
}

func (s *schemaChainSetter) Error() error {
	return s.err
}

func newNullableIntForID(i int) null.Int {
	// Because valid ID can't be null.
	return null.NewInt(int64(i), i != 0)
}

func listOfIDs(i interface{}) []int {
	if i == nil {
		return nil
	}

	vv := i.([]interface{}) //nolint:errcheck // We are sure about type.
	if len(vv) == 0 {
		return nil
	}

	res := make([]int, 0, len(vv))
	for _, v := range vv {
		res = append(res, v.(int))
	}

	return res
}

func adoptCreate(resourceName string, fn operationFunc) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		err := fn(ctx, m.(*client), d)
		if err != nil {
			err = fmt.Errorf("failed to create %s: %w", resourceName, err)
		}
		return diag.FromErr(err)
	}
}

func adoptRead(resourceName string, fn operationFunc) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		err := fn(ctx, m.(*client), d)
		if err != nil {
			err = fmt.Errorf("failed to read %s: %w", resourceName, err)
		}
		return diag.FromErr(err)
	}
}

func adoptUpdate(resourceName string, fn operationFunc) schema.UpdateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		err := fn(ctx, m.(*client), d)
		if err != nil {
			err = fmt.Errorf("failed to update %s: %w", resourceName, err)
		}
		return diag.FromErr(err)
	}
}

func adoptDelete(resourceName string, fn operationFunc) schema.DeleteContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		err := fn(ctx, m.(*client), d)
		if err != nil {
			err = fmt.Errorf("failed to delete %s: %w", resourceName, err)
		}
		return diag.FromErr(err)
	}
}

type operationFunc func(ctx context.Context, client *client, d *schema.ResourceData) error
