package solus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourcePlan(t *testing.T) {
	name := generateResourceName()
	resName := "solusio_plan." + name

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					`
resource "solusio_plan" "%[1]s" {
	name = "%[1]s"
	storage_type = "fb"
	image_format = "raw"
	params {
		disk = 1
		ram = 2
		vcpu = 3
	}
}

data "solusio_plan" "%[1]s_by_id" {
	id = solusio_plan.%[1]s.id
}

data "solusio_plan" "%[1]s_by_name" {
	name = solusio_plan.%[1]s.name
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
