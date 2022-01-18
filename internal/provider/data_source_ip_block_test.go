package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourceIPBlock(t *testing.T) {
	name := generateResourceName()
	resName := "solus_ip_block." + name

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckIPBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					`
resource "solus_ip_block" "%[1]s" {
	name = "%[1]s"
	ns1 = "192.0.2.1"
	ns2 = "192.0.2.2"
	gateway = "192.0.2.3"
	type = "IPv4"
	netmask = "255.255.255.0"
	from = "192.0.2.10"
	to = "192.0.2.20"
}

data "solus_ip_block" "%[1]s_by_id" {
	id = solus_ip_block.%[1]s.id
}

data "solus_ip_block" "%[1]s_by_name" {
	name = solus_ip_block.%[1]s.name
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
