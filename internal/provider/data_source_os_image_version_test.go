package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourceOsImageVersion(t *testing.T) {
	name := generateResourceName()
	resName := "solus_os_image_version." + name

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckOsImageVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					`
resource "solus_os_image" "%[1]s" {
	name = "%[1]s"
}

resource "solus_os_image_version" "%[1]s" {
	os_image_id = solus_os_image.%[1]s.id
	version = "%[1]s"
	url = "http://example.com/foo"
	cloud_init_version = "v2"
	virtualization_type = "kvm"
}

data "solus_os_image_version" "%[1]s_by_id" {
	id = solus_os_image_version.%[1]s.id
}

data "solus_os_image_version" "%[1]s_by_version" {
	os_image_id = solus_os_image.%[1]s.id
	version = solus_os_image_version.%[1]s.version
}
`,
					name,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data."+resName+"_by_id", "id"),
					resource.TestCheckResourceAttr("data."+resName+"_by_id", "version", name),

					resource.TestCheckResourceAttrSet("data."+resName+"_by_version", "id"),
					resource.TestCheckResourceAttr("data."+resName+"_by_version", "version", name),
				),
			},
		},
	})
}
