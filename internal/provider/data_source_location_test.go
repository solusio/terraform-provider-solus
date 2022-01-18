package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourceLocation(t *testing.T) {
	name := generateResourceName()
	resName := "solus_location." + name

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckLocationDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					`
resource "solus_location" "%[1]s" {
	name = "%[1]s"
	description = "for acc test"
}

data "solus_location" "%[1]s_by_id" {
	id = solus_location.%[1]s.id
}

data "solus_location" "%[1]s_by_name" {
	name = solus_location.%[1]s.name
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
