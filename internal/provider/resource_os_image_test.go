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

func TestAccResourceOsImage(t *testing.T) {
	name := generateResourceName()
	resName := "solus_os_image." + name

	checker := func(name string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(resName, "id"),
			resource.TestCheckResourceAttr(resName, "name", name),
			resource.TestCheckResourceAttrSet(resName, "icon_id"),
			resource.TestCheckResourceAttrSet(resName, "is_visible"),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckOsImageDestroy,
		Steps: []resource.TestStep{
			// Create resource.
			{
				Config: fmt.Sprintf(
					`
resource "solus_os_image" "%[1]s" {
	name = "%[1]s"
}
`,
					name,
				),
				Check: checker(name),
			},

			// Update created resource.
			{
				Config: fmt.Sprintf(
					`
resource "solus_os_image" "%s" {
	name = "%s"
}
`,
					name,
					name+"-changed",
				),
				Check: checker(name + "-changed"),
			},
		},
	})
}

func testAccCheckOsImageDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "solus_os_image" {
			continue
		}

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
	}

	return nil
}
