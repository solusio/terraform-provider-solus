package solus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourceOsImageVersion(t *testing.T) {
	name := generateResourceName()
	resName := "solusio_os_image_version." + name

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckOsImageVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					`
resource "solusio_os_image" "%[1]s" {
	name = "%[1]s"
}

resource "solusio_os_image_version" "%[1]s" {
	os_image_id = solusio_os_image.%[1]s.id
	version = "%[1]s"
	url = "http://example.com/foo"
	cloud_init_version = "v2"
}

data "solusio_os_image_version" "%[1]s_by_id" {
	id = solusio_os_image_version.%[1]s.id
}

data "solusio_os_image_version" "%[1]s_by_version" {
	os_image_id = solusio_os_image.%[1]s.id
	version = solusio_os_image_version.%[1]s.version
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
