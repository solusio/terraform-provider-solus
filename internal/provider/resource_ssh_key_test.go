package provider

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/solusio/solus-go-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccResourceSSHKey(t *testing.T) {
	const (
		body1 = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCvFALx+bFIGwzLTt8wFLgLgTpZaQLLndGL0Fv/5EAUSu6dGm7pFADDhQko7IUepziHiyRr8Y1AExg1z4QHW81JSH0XUxE3CTf0q91L1Ax2fBbYxLyi5jwt/OfqRQymwz6aRrEZb7BURkTI/VODDHstBAHtd44zQAbnC9QinkxRg+66gzfQlnu/e+sU2lm+1j9rrzIWUiYW9sae8XgDGmAB38gooakrqN9zlT3tlgsOWuhrgbTHLGg6iAHlfws6G820+ijTL8Z+9b6I3lOiWCPp1IPGtRoFuGhf8trN5vSCnxoXWC2iDytVUptPSRu/uR18ibEFYDx6G9+toNDtieAJ" //nolint:lll
		body2 = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFbNUAbScm4GGCTjKwgC4T/zitU9kdHFKvOp3U//bVFQ"
	)

	name := generateResourceName()
	resName := "solus_ssh_key." + name

	checker := func(name, body string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(resName, "id"),
			resource.TestCheckResourceAttr(resName, "name", name),
			resource.TestCheckResourceAttr(resName, "body", body),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSSHKeyDestroy,
		Steps: []resource.TestStep{
			// Create resource.
			{
				Config: fmt.Sprintf(
					`
resource "solus_ssh_key" "%[1]s" {
	name = "%[1]s"
	body = "%[2]s"
}
`,
					name,
					body1,
				),
				Check: checker(name, body1),
			},

			// Update created resource.
			{
				Config: fmt.Sprintf(
					`
resource "solus_ssh_key" "%s" {
	name = "%s"
	body = "%s"
}
`,
					name,
					name+"-changed",
					body2,
				),
				Check: checker(name+"-changed", body2),
			},
		},
	})
}

func testAccCheckSSHKeyDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "solus_SSHKey" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = c.SSHKeys.Get(context.Background(), id)
		if err == nil {
			return fmt.Errorf("SSH key %d still exists", id)
		}

		if !solus.IsNotFound(err) {
			return err
		}
	}

	return nil
}

func Test_checkPublicSSHKeyBody(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		ss := []string{
			"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCvFALx+bFIGwzLTt8wFLgLgTpZaQLLndGL0Fv/5EAUSu6dGm7pFADDhQko7IUepziHiyRr8Y1AExg1z4QHW81JSH0XUxE3CTf0q91L1Ax2fBbYxLyi5jwt/OfqRQymwz6aRrEZb7BURkTI/VODDHstBAHtd44zQAbnC9QinkxRg+66gzfQlnu/e+sU2lm+1j9rrzIWUiYW9sae8XgDGmAB38gooakrqN9zlT3tlgsOWuhrgbTHLGg6iAHlfws6G820+ijTL8Z+9b6I3lOiWCPp1IPGtRoFuGhf8trN5vSCnxoXWC2iDytVUptPSRu/uR18ibEFYDx6G9+toNDtieAJ user@example.com", //nolint:lll
			"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCvFALx+bFIGwzLTt8wFLgLgTpZaQLLndGL0Fv/5EAUSu6dGm7pFADDhQko7IUepziHiyRr8Y1AExg1z4QHW81JSH0XUxE3CTf0q91L1Ax2fBbYxLyi5jwt/OfqRQymwz6aRrEZb7BURkTI/VODDHstBAHtd44zQAbnC9QinkxRg+66gzfQlnu/e+sU2lm+1j9rrzIWUiYW9sae8XgDGmAB38gooakrqN9zlT3tlgsOWuhrgbTHLGg6iAHlfws6G820+ijTL8Z+9b6I3lOiWCPp1IPGtRoFuGhf8trN5vSCnxoXWC2iDytVUptPSRu/uR18ibEFYDx6G9+toNDtieAJ",                  //nolint:lll
			"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFbNUAbScm4GGCTjKwgC4T/zitU9kdHFKvOp3U//bVFQ root@example.com",                                                                                                                                                                             //nolint:lll
			"ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBKtLaqXUkYCPq5AgJuppFWTRjc8IU6OqfJtXwyoxqMt+GcDyIAyzD5LAcEapzurTzKXAiZrXM27vzJqKpmMV+90= root@example.com",                                                                                             //nolint:lll
			"ecdsa-sha2-nistp384 AAAAE2VjZHNhLXNoYTItbmlzdHAzODQAAAAIbmlzdHAzODQAAABhBFO7JVBzXZxXIs0noOjCNMZkdNTug/S6B4QNlcCnXJjNN3aKjDYcs23wUbTqGIrzgMjjud/ZEJMxN5U2OfwHcn+NXfHYAvdBq1/0Hom40AOKv1FYByPIDUUA2lFcZ/S9fQ== root@example.com",                                                 //nolint:lll
			"ecdsa-sha2-nistp521 AAAAE2VjZHNhLXNoYTItbmlzdHA1MjEAAAAIbmlzdHA1MjEAAACFBAFA6HcvQaG9hJS9PDuQ4nBMqcw5WE+hd8EFSJM37tVkKwkHpJAEyWZXzDygnAorJkmfaiEzW3eDqzLoksiJZGDt9QD42J4+9+2gVzUMZrnsG4ldQ+v+RD58vHgAL7Xw/cnWGGGS0x14gkK51rPNkctwYCFy+pa1s0CPHhJFqSN8tpSIcA== root@example.com", //nolint:lll
		}

		for _, s := range ss {
			t.Run(s, func(t *testing.T) {
				err := checkPublicSSHKeyBody(s)
				require.NoError(t, err)
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		cc := map[string]string{
			"sadfsd":                   "invalid format",
			"dsfds sadfsd":             `unsupported algorithm "dsfds"`,
			"ssh-rsa sadfsdasdf":       "decode body: illegal base64 data at input byte 8",
			"ssh-rsa sadfsdasdf a@b.c": "decode body: illegal base64 data at input byte 8",
			"ssh-rsa AAAAC3NzaC1lZDI1NTE5AAAAIFbNUAbScm4GGCTjKwgC4T/zitU9kdHFKvOp3U//bVFQ": "body have different algorithm", //nolint:lll
		}

		for given, expected := range cc {
			t.Run(given, func(t *testing.T) {
				err := checkPublicSSHKeyBody(given)
				assert.EqualError(t, err, expected)
			})
		}
	})
}
