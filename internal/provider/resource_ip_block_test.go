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

func TestAccResourceIPBlock(t *testing.T) {
	name := generateResourceName()
	resName := "solus_ip_block." + name

	checkerV4 := func(
		ns1, ns2 string,
		gateway string,
		netmask string,
		from, to string,
	) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(resName+"_ipv4", "id"),
			resource.TestCheckResourceAttr(resName+"_ipv4", "ns1", ns1),
			resource.TestCheckResourceAttr(resName+"_ipv4", "ns2", ns2),
			resource.TestCheckResourceAttr(resName+"_ipv4", "gateway", gateway),
			resource.TestCheckResourceAttr(resName+"_ipv4", "type", string(solus.IPv4)),
			resource.TestCheckResourceAttr(resName+"_ipv4", "netmask", netmask),
			resource.TestCheckResourceAttr(resName+"_ipv4", "from", from),
			resource.TestCheckResourceAttr(resName+"_ipv4", "to", to),
		)
	}

	checkerV6 := func(
		ns1, ns2 string,
		gateway string,
		r string,
		subnet int,
	) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(resName+"_ipv6", "id"),
			resource.TestCheckResourceAttr(resName+"_ipv6", "ns1", ns1),
			resource.TestCheckResourceAttr(resName+"_ipv6", "ns2", ns2),
			resource.TestCheckResourceAttr(resName+"_ipv6", "gateway", gateway),
			resource.TestCheckResourceAttr(resName+"_ipv6", "type", string(solus.IPv6)),
			resource.TestCheckResourceAttr(resName+"_ipv6", "range", r),
			resource.TestCheckResourceAttr(resName+"_ipv6", "subnet", strconv.Itoa(subnet)),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckIPBlockDestroy,
		Steps: []resource.TestStep{
			// Create resource.
			{
				Config: fmt.Sprintf(
					`
resource "solus_ip_block" "%[1]s_ipv4" {
	name = "%[1]s_ipv4"
	ns1 = "192.0.2.1"
	ns2 = "192.0.2.2"
	gateway = "192.0.2.3"
	type = "IPv4"
	netmask = "255.255.255.0"
	from = "192.0.2.10"
	to = "192.0.2.20"
}

resource "solus_ip_block" "%[1]s_ipv6" {
	name = "%[1]s_ipv6"
	ns1 = "2001:db8::1"
	ns2 = "2001:db8::2"
	gateway = "2001:db8::3"
	type = "IPv6"
	range = "2001:db8:10::/64"
	subnet = 120
}
`,
					name,
				),
				Check: resource.ComposeTestCheckFunc(
					checkerV4(
						"192.0.2.1",
						"192.0.2.2",
						"192.0.2.3",
						"255.255.255.0",
						"192.0.2.10",
						"192.0.2.20",
					),
					checkerV6(
						"2001:db8::1",
						"2001:db8::2",
						"2001:db8::3",
						"2001:db8:10::/64",
						120,
					),
				),
			},

			// Update created resource.
			{
				Config: fmt.Sprintf(
					`
resource "solus_ip_block" "%[1]s_ipv4" {
	name = "%[1]s_ipv4"
	ns1 = "192.0.2.15"
	ns2 = "192.0.2.25"
	gateway = "192.0.2.35"
	type = "IPv4"
	netmask = "255.255.255.128"
	from = "192.0.2.105"
	to = "192.0.2.205"
}

resource "solus_ip_block" "%[1]s_ipv6" {
	name = "%[1]s_ipv6"
	ns1 = "2001:db8::10"
	ns2 = "2001:db8::20"
	gateway = "2001:db8::30"
	type = "IPv6"
	range = "2001:db8:15::/67"
	subnet = 116
}
`,
					name,
				),
				Check: resource.ComposeTestCheckFunc(
					checkerV4(
						"192.0.2.15",
						"192.0.2.25",
						"192.0.2.35",
						"255.255.255.128",
						"192.0.2.105",
						"192.0.2.205",
					),
					checkerV6(
						"2001:db8::10",
						"2001:db8::20",
						"2001:db8::30",
						"2001:db8:15::/67",
						116,
					),
				),
			},
		},
	})
}

func testAccCheckIPBlockDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*solus.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "solus_ip_block" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = c.IPBlocks.Get(context.Background(), id)
		if err == nil {
			return fmt.Errorf("location %d still exists", id)
		}

		if !solus.IsNotFound(err) {
			return err
		}
	}

	return nil
}
