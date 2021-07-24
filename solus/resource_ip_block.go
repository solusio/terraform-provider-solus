package solus

import (
	"context"
	"errors"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func resourceIPBlock() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIPBlockCreate,
		ReadContext:   resourceIPBlockRead,
		UpdateContext: resourceIPBlockUpdate,
		DeleteContext: resourceIPBlockDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"ns1": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.Any(
					validation.IsIPAddress,
					validationIsDomainName,
				),
			},
			"ns2": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.Any(
					validation.IsIPAddress,
					validationIsDomainName,
				),
			},
			"gateway": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(solus.IPv4),
					string(solus.IPv6),
				}, false),
			},

			// Unfortunately terraform plugin SDK don't give to us any tools to
			// validate field value according another field value so we have to
			// make a request and only after that we will get an error.

			// Type - IPv4
			"netmask": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.IsIPv4Address,
				ConflictsWith: []string{"range", "subnet"},
			},
			"from": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.IsIPv4Address,
				ConflictsWith: []string{"range", "subnet"},
			},
			"to": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.IsIPv4Address,
				ConflictsWith: []string{"range", "subnet"},
			},

			// Type - IPv6
			"range": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.IsCIDRNetwork(0, 128),
				ConflictsWith: []string{"netmask", "from", "to"},
			},
			"subnet": {
				Type:          schema.TypeInt,
				Optional:      true,
				ValidateFunc:  validation.IntBetween(83, 128),
				ConflictsWith: []string{"netmask", "from", "to"},
			},
		},
	}
}

func resourceIPBlockCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, ok := m.(*solus.Client)
	if !ok {
		return diag.Errorf("invalid Solus client type %T", m)
	}

	req, err := buildIPBlockRequest(d)
	if err != nil {
		return diag.Errorf("failed to build request: %s", err)
	}

	res, err := client.IPBlocks.Create(ctx, req)
	if err != nil {
		return diag.Errorf("failed to create new ip block: %s", err)
	}

	d.SetId(strconv.Itoa(res.ID))
	return resourceIPBlockRead(ctx, d, m)
}

func resourceIPBlockRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, ok := m.(*solus.Client)
	if !ok {
		return diag.Errorf("invalid Solus client type %T", m)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := client.IPBlocks.Get(ctx, id)
	if err != nil {
		return diag.Errorf("failed to get ip block by id %d: %s", id, err)
	}

	s := (&schemaChainSetter{d: d}).
		SetID(res.ID).
		Set("name", res.Name).
		Set("ns1", res.Ns1).
		Set("ns2", res.Ns2).
		Set("gateway", res.Gateway).
		Set("type", res.Type)

	switch res.Type {
	case solus.IPv4:
		s.
			Set("netmask", res.Netmask).
			Set("from", res.From).
			Set("to", res.To)

	case solus.IPv6:
		s.
			Set("range", res.Range).
			Set("subnet", res.Subnet)
	default:
		return diag.Errorf("unhandled IP type %q", res.Type)
	}

	if err := s.Error(); err != nil {
		return diag.Errorf("failed to map ip block response to resource: %s", err)
	}

	return nil
}

func resourceIPBlockUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, ok := m.(*solus.Client)
	if !ok {
		return diag.Errorf("invalid Solus client type %T", m)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := buildIPBlockRequest(d)
	if err != nil {
		return diag.Errorf("failed to build request: %s", err)
	}

	res, err := client.IPBlocks.Update(ctx, id, req)
	if err != nil {
		return diag.Errorf("failed to update ip block with id %d: %s", id, err)
	}

	d.SetId(strconv.Itoa(res.ID))
	return resourceIPBlockRead(ctx, d, m)
}

func resourceIPBlockDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, ok := m.(*solus.Client)
	if !ok {
		return diag.Errorf("invalid Solus client type %T", m)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.IPBlocks.Delete(ctx, id)
	if err != nil {
		return diag.Errorf("failed to delete ip block by id %d: %s", id, err)
	}
	return nil
}

func buildIPBlockRequest(d *schema.ResourceData) (solus.IPBlockRequest, error) {
	req := solus.IPBlockRequest{
		Name:    d.Get("name").(string),
		Ns1:     d.Get("ns1").(string),
		Ns2:     d.Get("ns2").(string),
		Gateway: d.Get("gateway").(string),
		Type:    solus.IPVersion(d.Get("type").(string)),
	}

	switch req.Type {
	case solus.IPv4:
		req.Netmask = d.Get("netmask").(string)
		req.From = d.Get("from").(string)
		req.To = d.Get("to").(string)

	case solus.IPv6:
		req.Range = d.Get("range").(string)
		req.Subnet = d.Get("subnet").(int)

	default:
		return solus.IPBlockRequest{}, errors.New("unhandled IP version")
	}
	return req, nil
}
