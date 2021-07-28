package solus

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
	resName := "solusio_os_image_version." + name

	checker := func(version, url, cloudInitVersion string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(resName, "id"),
			resource.TestCheckResourceAttrSet(resName, "os_image_id"),
			resource.TestCheckResourceAttr(resName, "version", version),
			resource.TestCheckResourceAttr(resName, "url", url),
			resource.TestCheckResourceAttr(resName, "cloud_init_version", cloudInitVersion),
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
resource "solusio_os_image" "%[1]s" {
	name = "%[1]s"
}

resource "solusio_os_image_version" "%[1]s" {
	os_image_id = solusio_os_image.%[1]s.id
	version = "version"
	url = "http://example.com/foo"
	cloud_init_version = "v2"
}
`,
					name,
				),
				Check: checker("version", "http://example.com/foo", "v2"),
			},

			// Update created resource.
			{
				Config: fmt.Sprintf(
					`
resource "solusio_os_image" "%[1]s" {
	name = "%[1]s"
}

resource "solusio_os_image_version" "%[1]s" {
	os_image_id = solusio_os_image.%[1]s.id
	version = "version_changed"
	url = "http://example.com/bar"
	cloud_init_version = "v0"
}
`,
					name,
				),
				Check: checker("version_changed", "http://example.com/bar", "v0"),
			},

			// Move to another os image.
			{
				Config: fmt.Sprintf(
					`
resource "solusio_os_image" "%[1]s" {
	name = "%[1]s"
}

resource "solusio_os_image" "%[1]s_new" {
	name = "%[1]s_new"
}

resource "solusio_os_image_version" "%[1]s" {
	os_image_id = solusio_os_image.%[1]s_new.id
	version = "version_changed"
	url = "http://example.com/bar"
	cloud_init_version = "v0"
}
`,
					name,
				),
				Check: checker("version_changed", "http://example.com/bar", "v0"),
			},
		},
	})
}

func testAccCheckOsImageVersionDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*solus.Client)

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "solusio_os_image":
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

		case "solusio_os_image_version":
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
