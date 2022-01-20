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

func TestAccResourceProject(t *testing.T) {
	const description = "for acc test"
	name := generateResourceName()
	resName := "solus_project." + name

	checker := func(name, description string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(resName, "id"),
			resource.TestCheckResourceAttr(resName, "name", name),
			resource.TestCheckResourceAttr(resName, "description", description),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckProjectDestroy,
		Steps: []resource.TestStep{
			// Create resource.
			{
				Config: fmt.Sprintf(
					`
resource "solus_project" "%[1]s" {
	name = "%[1]s"
	description = "%[2]s"
}
`,
					name,
					description,
				),
				Check: checker(name, description),
			},

			// Update created resource.
			{
				Config: fmt.Sprintf(
					`
resource "solus_project" "%s" {
	name = "%s"
	description = "%s"
}
`,
					name,
					name+"-changed",
					description+"-changed",
				),
				Check: checker(name+"-changed", description+"-changed"),
			},
		},
	})
}

func testAccCheckProjectDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "solus_project" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = c.Projects.Get(context.Background(), id)
		if err == nil {
			return fmt.Errorf("location %d still exists", id)
		}

		if !solus.IsNotFound(err) {
			return err
		}
	}

	return nil
}
