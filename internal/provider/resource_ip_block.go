package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func resourceIPBlock() *schema.Resource {
	return &schema.Resource{
		CreateContext: adoptCreate("IP Block", resourceIPBlockCreate),
		ReadContext:   adoptRead("IP Block", resourceIPBlockRead),
		UpdateContext: adoptUpdate("IP Block", resourceIPBlockUpdate),
		DeleteContext: adoptDelete("IP Block", resourceIPBlockDelete),

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

func resourceIPBlockCreate(ctx context.Context, client *client, d *schema.ResourceData) error {
	req, err := buildIPBlockRequest(d)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	res, err := client.IPBlocks.Create(ctx, req)
	if err != nil {
		return normalizeAPIError(err)
	}

	d.SetId(strconv.Itoa(res.ID))
	return resourceIPBlockRead(ctx, client, d)
}

func resourceIPBlockRead(ctx context.Context, client *client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	res, err := client.IPBlocks.Get(ctx, id)
	if err != nil {
		return normalizeAPIError(err)
	}

	s := newSchemaChainSetter(d).
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
		return fmt.Errorf("unhandled IP type %q", res.Type)
	}

	return s.Error()
}

func resourceIPBlockUpdate(ctx context.Context, client *client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	req, err := buildIPBlockRequest(d)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	res, err := client.IPBlocks.Update(ctx, id, req)
	if err != nil {
		return normalizeAPIError(err)
	}

	d.SetId(strconv.Itoa(res.ID))
	return resourceIPBlockRead(ctx, client, d)
}

func resourceIPBlockDelete(ctx context.Context, client *client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	return normalizeAPIError(client.IPBlocks.Delete(ctx, id))
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
