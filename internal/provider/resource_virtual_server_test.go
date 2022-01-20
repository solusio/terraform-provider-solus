package provider

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/solusio/solus-go-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

func TestAccResourceVirtualServer(t *testing.T) {
	const (
		description = "for acc test"
		hostname    = "vs.example.com"
	)
	name := generateResourceName()
	resName := "solus_virtual_server." + name

	pubKey, err := generateSSHPublicKey()
	require.NoError(t, err)

	var locationID int
	if raw := os.Getenv("SOLUS_TEST_LOCATION_ID"); raw != "" {
		var err error
		locationID, err = strconv.Atoi(raw)
		require.NoError(t, err)
	}

	checker := func(hostname, description string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(resName, "id"),
			resource.TestCheckResourceAttr(resName, "hostname", hostname),
			resource.TestCheckResourceAttr(resName, "description", description),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			assert.NotEmpty(
				t,
				os.Getenv("SOLUS_TEST_LOCATION_ID"),
				`"SOLUS_TEST_LOCATION_ID" environment variable must be set for acceptance tests`,
			)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckVirtualServerDestroy,
		Steps: []resource.TestStep{
			// Create resource.
			{
				// todo remove hardcoded location id
				Config: fmt.Sprintf(`
data "solus_location" "%[1]s" {
	id = %[4]d
}

resource "solus_project" "%[1]s" {
	name = "%[1]s"
}

resource "solus_os_image" "%[1]s" {
	name = "Alpine %[1]s"
}

resource "solus_os_image_version" "%[1]s" {
	os_image_id = solus_os_image.%[1]s.id
	version = "%[1]s"
	url = "https://images.prod.solus.io/solus-alpine-3.15.qcow2"
	cloud_init_version = "v2"
	virtualization_type = "kvm"
}

resource "solus_plan" "%[1]s" {
	name = "%[1]s"
	virtualization_type = "kvm"
	storage_type = "fb"
	image_format = "qcow2"
	params {
		disk = 1
		ram_mb = 1024
		vcpu = 1
	}
	available_locations = [
		data.solus_location.%[1]s.id
	]
	available_os_image_versions = [
		solus_os_image_version.%[1]s.id
	]
}

resource "solus_ssh_key" "%[1]s" {
	name = "%[1]s"
	body = "%[5]s"
}

resource "solus_virtual_server" "%[1]s" {
	hostname = "%[2]s"
	description = "%[3]s"

	location_id = data.solus_location.%[1]s.id

	os_image_version_id = solus_os_image_version.%[1]s.id

	plan_id = solus_plan.%[1]s.id
	project_id = solus_project.%[1]s.id
	ssh_keys = [
		solus_ssh_key.%[1]s.id,
	]
	user_data = <<EOT
#cloud-config
runcmd:
  - echo "echo "Hello from User Data"" > /root/hello-world-user-data.sh
  - [ chmod, +x, "/root/hello-world-user-data.sh" ]
EOT
}
`,
					name,
					hostname,
					description,
					locationID,
					pubKey,
				),
				Check: checker(hostname, description),
			},
		},
	})
}

func testAccCheckVirtualServerDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "solus_virtual_server" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = c.VirtualServers.Get(context.Background(), id)
		if err == nil {
			return fmt.Errorf("location %d still exists", id)
		}

		if !solus.IsNotFound(err) {
			return err
		}
	}

	return nil
}

func generateSSHPublicKey() (string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return "", err
	}

	var privKeyBuf strings.Builder

	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(&privKeyBuf, privateKeyPEM); err != nil {
		return "", err
	}

	// generate and write public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", err
	}

	var pubKeyBuf strings.Builder
	pubKeyBuf.Write(ssh.MarshalAuthorizedKey(pub))

	return strings.TrimSpace(pubKeyBuf.String()), nil
}
