package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourcePlan(t *testing.T) {
	name := generateResourceName()
	resName := "solus_plan." + name

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					`
resource "solus_plan" "%[1]s" {
	name = "%[1]s"
	virtualization_type = "kvm"
	storage_type = "fb"
	image_format = "raw"
	params {
		disk = 1
		ram_mb = 128
		vcpu = 3
	}
}

data "solus_plan" "%[1]s_by_id" {
	id = solus_plan.%[1]s.id
}

data "solus_plan" "%[1]s_by_name" {
	name = solus_plan.%[1]s.name
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
