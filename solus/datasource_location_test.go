package solus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourceLocation(t *testing.T) {
	name := generateResourceName()
	resName := "solusio_location." + name

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckLocationDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					`
resource "solusio_location" "%[1]s" {
	name = "%[1]s"
	description = "for acc test"
}

data "solusio_location" "%[1]s_by_id" {
	id = solusio_location.%[1]s.id
}

data "solusio_location" "%[1]s_by_name" {
	name = solusio_location.%[1]s.name
}
`,
					name,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data."+resName+"_by_id", "id"),
					resource.TestCheckResourceAttr("data."+resName+"_by_id", "name", name),

					resource.TestCheckResourceAttrSet("data."+resName+"_by_name", "id"),
					resource.TestCheckResourceAttr("data."+resName+"_by_name", "name", name),
				),
			},
		},
	})
}
