package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAccDatasourceSSHKey(t *testing.T) {
	name := generateResourceName()
	resName := "solus_ssh_key." + name

	pubKey, err := generateSSHPublicKey()
	require.NoError(t, err)

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckLocationDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					`
resource "solus_ssh_key" "%[1]s" {
	name = "%[1]s"
	body = "%[2]s"
}

data "solus_ssh_key" "%[1]s_by_id" {
	id = solus_ssh_key.%[1]s.id
}

data "solus_ssh_key" "%[1]s_by_name" {
	name = solus_ssh_key.%[1]s.name
}
`,
					name,
					pubKey,
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
