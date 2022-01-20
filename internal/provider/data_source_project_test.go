package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourceProject(t *testing.T) {
	name := generateResourceName()
	resName := "solus_project." + name

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					`
resource "solus_project" "%[1]s" {
	name = "%[1]s"
	description = "for acc test"
}

data "solus_project" "%[1]s_by_id" {
	id = solus_project.%[1]s.id
}
`,
					name,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data."+resName+"_by_id", "id"),
				),
			},
		},
	})
}
