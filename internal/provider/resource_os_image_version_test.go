package provider

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/solusio/solus-go-sdk"
)

func TestAccResourceOsImageVersion(t *testing.T) {
	name := generateResourceName()
	resName := "solus_os_image_version." + name

	checker := func(version, url, cloudInitVersion string, vt solus.VirtualizationType) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(resName, "id"),
			resource.TestCheckResourceAttrSet(resName, "os_image_id"),
			resource.TestCheckResourceAttr(resName, "version", version),
			resource.TestCheckResourceAttr(resName, "url", url),
			resource.TestCheckResourceAttr(resName, "cloud_init_version", cloudInitVersion),
			resource.TestCheckResourceAttr(resName, "virtualization_type", string(vt)),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckOsImageVersionDestroy,
		Steps: []resource.TestStep{
			// Create resource.
			{
				Config: fmt.Sprintf(
					`
resource "solus_os_image" "%[1]s" {
	name = "%[1]s"
}

resource "solus_os_image_version" "%[1]s" {
	os_image_id = solus_os_image.%[1]s.id
	version = "version"
	url = "http://example.com/foo"
	cloud_init_version = "v2"
	virtualization_type = "kvm"
}
`,
					name,
				),
				Check: checker("version", "http://example.com/foo", "v2", solus.VirtualizationTypeKVM),
			},

			// Update created resource.
			{
				Config: fmt.Sprintf(
					`
resource "solus_os_image" "%[1]s" {
	name = "%[1]s"
}

resource "solus_os_image_version" "%[1]s" {
	os_image_id = solus_os_image.%[1]s.id
	version = "version_changed"
	url = "http://example.com/bar"
	cloud_init_version = "v0"
	virtualization_type = "kvm"
}
`,
					name,
				),
				Check: checker("version_changed", "http://example.com/bar", "v0", solus.VirtualizationTypeKVM),
			},

			// Move to another os image.
			{
				Config: fmt.Sprintf(
					`
resource "solus_os_image" "%[1]s" {
	name = "%[1]s"
}

resource "solus_os_image" "%[1]s_new" {
	name = "%[1]s_new"
}

resource "solus_os_image_version" "%[1]s" {
	os_image_id = solus_os_image.%[1]s_new.id
	version = "version_changed"
	url = "http://example.com/bar"
	cloud_init_version = "v0"
	virtualization_type = "kvm"
}
`,
					name,
				),
				Check: checker("version_changed", "http://example.com/bar", "v0", solus.VirtualizationTypeKVM),
			},
		},
	})
}

func testAccCheckOsImageVersionDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*solus.Client)

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "solus_os_image":
			id, err := strconv.Atoi(rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = c.OsImages.Get(context.Background(), id)
			if err == nil {
				return fmt.Errorf("os image %d still exists", id)
			}

			if !solus.IsNotFound(err) {
				return err
			}

		case "solus_os_image_version":
			id, err := strconv.Atoi(rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = c.OsImageVersions.Get(context.Background(), id)
			if err == nil {
				return fmt.Errorf("os image version %d still exists", id)
			}

			if !solus.IsNotFound(err) {
				return err
			}
		}
	}
	return nil
}
