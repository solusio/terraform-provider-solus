package provider

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func resourceSSHKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: adoptCreate("SSH Key", resourceSSHKeyCreate),
		ReadContext:   adoptRead("SSH Key", resourceSSHKeyRead),
		DeleteContext: adoptDelete("SSH Key", resourceSSHKeyDelete),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"body": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(i interface{}, k string) ([]string, []error) {
					v, ok := i.(string)
					if !ok {
						return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
					}

					err := checkPublicSSHKeyBody(v)
					if err != nil {
						return nil, []error{err}
					}
					return nil, nil
				},
			},
		},
	}
}

func resourceSSHKeyCreate(ctx context.Context, client *client, d *schema.ResourceData) error {
	u, err := client.CurrentUser(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	res, err := client.SSHKeys.Create(ctx, solus.SSHKeyCreateRequest{
		Name:   d.Get("name").(string),
		Body:   d.Get("body").(string),
		UserID: u.ID,
	})
	if err != nil {
		return normalizeAPIError(err)
	}

	d.SetId(strconv.Itoa(res.ID))
	return resourceSSHKeyRead(ctx, client, d)
}

func resourceSSHKeyRead(ctx context.Context, client *client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	res, err := client.SSHKeys.Get(ctx, id)
	if err != nil {
		return normalizeAPIError(err)
	}

	return newSchemaChainSetter(d).
		Set("name", res.Name).
		Set("body", res.Body).
		Error()
}

func resourceSSHKeyDelete(ctx context.Context, client *client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	return normalizeAPIError(client.SSHKeys.Delete(ctx, id))
}

func checkPublicSSHKeyBody(s string) error {
	validAlgs := map[string]struct{}{
		"ssh-rsa":             {},
		"ssh-dss":             {},
		"ssh-ed25519":         {},
		"ecdsa-sha2-nistp256": {}, // cspell:disable-line
		"ecdsa-sha2-nistp384": {}, // cspell:disable-line
		"ecdsa-sha2-nistp521": {}, // cspell:disable-line
	}

	parts := strings.Split(s, " ")
	if len(parts) != 2 && len(parts) != 3 {
		return errors.New("invalid format")
	}

	alg := parts[0]
	body := parts[1]

	if _, ok := validAlgs[alg]; !ok {
		return fmt.Errorf("unsupported algorithm %q", alg)
	}

	buf := make([]byte, base64.StdEncoding.DecodedLen(len(body)))
	n, err := base64.StdEncoding.Decode(buf, []byte(body))
	if err != nil {
		return fmt.Errorf("decode body: %w", err)
	}
	buf = buf[:n]

	if !bytes.Contains(buf, []byte(alg)) {
		return errors.New("body have different algorithm")
	}
	return nil
}
