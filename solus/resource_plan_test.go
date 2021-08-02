package solus

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/solusio/solus-go-sdk"
)

func TestAccResourcePlan(t *testing.T) {
	name := generateResourceName()

	checker := func(
		resName string,
		name string,
		storageType solus.StorageTypeName,
		imageFormat solus.ImageFormat,
		params solus.PlanParams,
		tokenPerHour, tokenPerMonth int,
		ipTokenPerHour, ipTokenPerMonth int,
		backupPrice float32,
		resetLimitPolicy solus.PlanResetLimitPolicy,
		networkTrafficLimitType solus.PlanNetworkTotalTrafficType,
	) resource.TestCheckFunc {
		convBackupPrice := func(f float32) string {
			return fmt.Sprintf("%g", f)
		}

		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(resName, "id"),
			resource.TestCheckResourceAttr(resName, "name", name),
			resource.TestCheckResourceAttr(resName, "storage_type", string(storageType)),
			resource.TestCheckResourceAttr(resName, "image_format", string(imageFormat)),
			resource.TestCheckResourceAttr(resName, "params.0.disk", strconv.Itoa(params.Disk)),
			resource.TestCheckResourceAttr(resName, "params.0.ram", strconv.Itoa(params.RAM)),
			resource.TestCheckResourceAttr(resName, "params.0.vcpu", strconv.Itoa(params.VCPU)),
			resource.TestCheckResourceAttr(resName, "tokens_per_hour", strconv.Itoa(tokenPerHour)),
			resource.TestCheckResourceAttr(resName, "tokens_per_month", strconv.Itoa(tokenPerMonth)),
			resource.TestCheckResourceAttr(resName, "ip_tokens_per_hour", strconv.Itoa(ipTokenPerHour)),
			resource.TestCheckResourceAttr(resName, "ip_tokens_per_month", strconv.Itoa(ipTokenPerMonth)),
			resource.TestCheckResourceAttr(resName, "backup_price", convBackupPrice(backupPrice)),
			resource.TestCheckResourceAttr(resName, "reset_limit_policy", string(resetLimitPolicy)),
			resource.TestCheckResourceAttr(resName, "network_traffic_limit_type", string(networkTrafficLimitType)),
			resource.TestCheckResourceAttrSet(resName, "is_default"),
			resource.TestCheckResourceAttrSet(resName, "is_visible"),
			resource.TestCheckResourceAttrSet(resName, "is_additional_ips_available"),
			resource.TestCheckResourceAttrSet(resName, "is_snapshots_enabled"),
			resource.TestCheckResourceAttrSet(resName, "is_backup_available"),
		)
	}

	limitsChecker := func(
		resName string,
		limits solus.PlanLimits,
	) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.disk_iops.0.is_enabled",
				fmt.Sprintf("%t", limits.DiskIOPS.IsEnabled),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.disk_iops.0.limit",
				strconv.Itoa(limits.DiskIOPS.Limit),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.disk_iops.0.unit",
				string(limits.DiskIOPS.Unit),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.disk_bandwidth.0.is_enabled",
				fmt.Sprintf("%t", limits.DiskBandwidth.IsEnabled),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.disk_bandwidth.0.limit",
				strconv.Itoa(limits.DiskBandwidth.Limit),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.disk_bandwidth.0.unit",
				string(limits.DiskBandwidth.Unit),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.network_incoming_bandwidth.0.is_enabled",
				fmt.Sprintf("%t", limits.NetworkIncomingBandwidth.IsEnabled),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.network_incoming_bandwidth.0.limit",
				strconv.Itoa(limits.NetworkIncomingBandwidth.Limit),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.network_incoming_bandwidth.0.unit",
				string(limits.NetworkIncomingBandwidth.Unit),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.network_outgoing_bandwidth.0.is_enabled",
				fmt.Sprintf("%t", limits.NetworkOutgoingBandwidth.IsEnabled),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.network_outgoing_bandwidth.0.limit",
				strconv.Itoa(limits.NetworkOutgoingBandwidth.Limit),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.network_outgoing_bandwidth.0.unit",
				string(limits.NetworkOutgoingBandwidth.Unit),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.network_incoming_traffic.0.is_enabled",
				fmt.Sprintf("%t", limits.NetworkIncomingTraffic.IsEnabled),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.network_incoming_traffic.0.limit",
				strconv.Itoa(limits.NetworkIncomingTraffic.Limit),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.network_incoming_traffic.0.unit",
				string(limits.NetworkIncomingTraffic.Unit),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.network_outgoing_traffic.0.is_enabled",
				fmt.Sprintf("%t", limits.NetworkOutgoingTraffic.IsEnabled),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.network_outgoing_traffic.0.limit",
				strconv.Itoa(limits.NetworkOutgoingTraffic.Limit),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.network_outgoing_traffic.0.unit",
				string(limits.NetworkOutgoingTraffic.Unit),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.network_total_traffic.0.is_enabled",
				fmt.Sprintf("%t", limits.NetworkTotalTraffic.IsEnabled),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.network_total_traffic.0.limit",
				strconv.Itoa(limits.NetworkTotalTraffic.Limit),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.network_total_traffic.0.unit",
				string(limits.NetworkTotalTraffic.Unit),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.network_reduce_bandwidth.0.is_enabled",
				fmt.Sprintf("%t", limits.NetworkReduceBandwidth.IsEnabled),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.network_reduce_bandwidth.0.limit",
				strconv.Itoa(limits.NetworkReduceBandwidth.Limit),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.network_reduce_bandwidth.0.unit",
				string(limits.NetworkReduceBandwidth.Unit),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.backups_number.0.is_enabled",
				fmt.Sprintf("%t", limits.BackupsNumber.IsEnabled),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.backups_number.0.limit",
				strconv.Itoa(limits.BackupsNumber.Limit),
			),
			resource.TestCheckResourceAttr(
				resName,
				"limits.0.backups_number.0.unit",
				string(limits.BackupsNumber.Unit),
			),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckPlanDestroy,
		Steps: []resource.TestStep{
			// Create resource with all default values.
			{
				Config: fmt.Sprintf(
					`
resource "solusio_plan" "%[1]s_default" {
	name = "%[1]s_default"
	storage_type = "fb"
	image_format = "raw"
	params {
		disk = 1
		ram = 2
		vcpu = 3
	}
}
`,
					name,
				),
				Check: checker(
					fmt.Sprintf("solusio_plan.%s_default", name),
					name+"_default",
					solus.StorageTypeNameFB,
					solus.ImageFormatRaw,
					solus.PlanParams{
						Disk: 1,
						RAM:  2,
						VCPU: 3,
					},
					0,
					0,
					0,
					0,
					0,
					solus.PlanResetLimitPolicyNever,
					solus.PlanNetworkTotalTrafficTypeSeparate,
				),
			},

			// Create resource with some properties.
			{
				Config: fmt.Sprintf(
					`
resource "solusio_plan" "%[1]s_some" {
	name = "%[1]s_some"
	storage_type = "fb"
	image_format = "raw"
	params {
		disk = 1
		ram = 2
		vcpu = 3
	}
	tokens_per_hour = 4
	tokens_per_month = 5
	ip_tokens_per_hour = 6
	ip_tokens_per_month = 7

	limits {
		network_outgoing_bandwidth {
			is_enabled = true
			limit = 50
			unit = "Mbps"
		}
		backups_number {
			is_enabled = true
			limit = 100
		}
	}

	backup_price = 3.14
	reset_limit_policy = "vm_created_day"
	network_traffic_limit_type = "total"
}
`,
					name,
				),
				Check: resource.ComposeTestCheckFunc(
					checker(
						fmt.Sprintf("solusio_plan.%s_some", name),
						name+"_some",
						solus.StorageTypeNameFB,
						solus.ImageFormatRaw,
						solus.PlanParams{
							Disk: 1,
							RAM:  2,
							VCPU: 3,
						},
						4,
						5,
						6,
						7,
						3.14,
						solus.PlanResetLimitPolicyVMCreatedDay,
						solus.PlanNetworkTotalTrafficTypeTotal,
					),
					resource.TestCheckResourceAttr(
						fmt.Sprintf("solusio_plan.%s_some", name),
						"limits.0.network_outgoing_bandwidth.0.is_enabled",
						"true",
					),
					resource.TestCheckResourceAttr(
						fmt.Sprintf("solusio_plan.%s_some", name),
						"limits.0.network_outgoing_bandwidth.0.limit",
						"50",
					),
					resource.TestCheckResourceAttr(
						fmt.Sprintf("solusio_plan.%s_some", name),
						"limits.0.network_outgoing_bandwidth.0.unit",
						"Mbps",
					),
					resource.TestCheckResourceAttr(
						fmt.Sprintf("solusio_plan.%s_some", name),
						"limits.0.backups_number.0.is_enabled",
						"true",
					),
					resource.TestCheckResourceAttr(
						fmt.Sprintf("solusio_plan.%s_some", name),
						"limits.0.backups_number.0.limit",
						"100",
					),
					resource.TestCheckResourceAttr(
						fmt.Sprintf("solusio_plan.%s_some", name),
						"limits.0.backups_number.0.unit",
						"units",
					),
				),
			},

			// Create resource with all properties.
			{
				Config: fmt.Sprintf(
					`
resource "solusio_plan" "%[1]s_full" {
	name = "%[1]s_full"
	storage_type = "fb"
	image_format = "raw"
	params {
		disk = 1
		ram = 2
		vcpu = 3
	}
	tokens_per_hour = 4
	tokens_per_month = 5
	ip_tokens_per_hour = 6
	ip_tokens_per_month = 7

	limits {
		disk_iops {
			is_enabled = true
			limit = 20
		}
		disk_bandwidth {
			is_enabled = true
			limit = 30
		}
		network_incoming_bandwidth {
			is_enabled = true
			limit = 40
			unit = "Kbps"
		}
		network_outgoing_bandwidth {
			is_enabled = true
			limit = 50
			unit = "Mbps"
		}
		network_incoming_traffic {
			is_enabled = true
			limit = 60
			unit = "KB"
		}
		network_outgoing_traffic {
			is_enabled = true
			limit = 70
			unit = "MB"
		}
		network_total_traffic {
			is_enabled = true
			limit = 80
			unit = "TB"
		}
		network_reduce_bandwidth {
			is_enabled = true
			limit = 90
			unit = "Gbps"
		}
		backups_number {
			is_enabled = true
			limit = 100
		}
	}

	backup_price = 3.14
	reset_limit_policy = "vm_created_day"
	network_traffic_limit_type = "total"
}
`,
					name,
				),
				Check: resource.ComposeTestCheckFunc(
					checker(
						fmt.Sprintf("solusio_plan.%s_full", name),
						name+"_full",
						solus.StorageTypeNameFB,
						solus.ImageFormatRaw,
						solus.PlanParams{
							Disk: 1,
							RAM:  2,
							VCPU: 3,
						},
						4,
						5,
						6,
						7,
						3.14,
						solus.PlanResetLimitPolicyVMCreatedDay,
						solus.PlanNetworkTotalTrafficTypeTotal,
					),
					limitsChecker(
						fmt.Sprintf("solusio_plan.%s_full", name),
						solus.PlanLimits{
							DiskIOPS: solus.DiskIOPSPlanLimit{
								IsEnabled: true,
								Limit:     20,
								Unit:      solus.DiskIOPSPlanLimitUnitIOPS,
							},
							DiskBandwidth: solus.DiskBandwidthPlanLimit{
								IsEnabled: true,
								Limit:     30,
								Unit:      solus.DiskBandwidthPlanLimitUnitBps,
							},
							NetworkIncomingBandwidth: solus.BandwidthPlanLimit{
								IsEnabled: true,
								Limit:     40,
								Unit:      solus.BandwidthPlanLimitUnitKbps,
							},
							NetworkOutgoingBandwidth: solus.BandwidthPlanLimit{
								IsEnabled: true,
								Limit:     50,
								Unit:      solus.BandwidthPlanLimitUnitMbps,
							},
							NetworkIncomingTraffic: solus.TrafficPlanLimit{
								IsEnabled: true,
								Limit:     60,
								Unit:      solus.TrafficPlanLimitUnitKB,
							},
							NetworkOutgoingTraffic: solus.TrafficPlanLimit{
								IsEnabled: true,
								Limit:     70,
								Unit:      solus.TrafficPlanLimitUnitMB,
							},
							NetworkTotalTraffic: solus.TrafficPlanLimit{
								IsEnabled: true,
								Limit:     80,
								Unit:      solus.TrafficPlanLimitUnitTB,
							},
							NetworkReduceBandwidth: solus.BandwidthPlanLimit{
								IsEnabled: true,
								Limit:     90,
								Unit:      solus.BandwidthPlanLimitUnitGbps,
							},
							BackupsNumber: solus.UnitPlanLimit{
								IsEnabled: true,
								Limit:     100,
								Unit:      solus.PlanLimitUnits,
							},
						},
					),
				),
			},

			// Update plan but not touch fields which may lead to creating new
			// resource.
			{
				Config: fmt.Sprintf(
					`
resource "solusio_plan" "%[1]s_full" {
	name = "%[1]s_full_new"
	storage_type = "fb"
	image_format = "raw"
	params {
		disk = 1
		ram = 2
		vcpu = 3
	}
	tokens_per_hour = 40
	tokens_per_month = 50
	ip_tokens_per_hour = 60
	ip_tokens_per_month = 70

	limits {
		disk_iops {
			is_enabled = true
			limit = 20
		}
		disk_bandwidth {
			is_enabled = true
			limit = 30
		}
		network_incoming_bandwidth {
			is_enabled = true
			limit = 4
			unit = "Gbps"
		}
		network_outgoing_bandwidth {
			is_enabled = true
			limit = 5
			unit = "Kbps"
		}
		network_incoming_traffic {
			is_enabled = true
			limit = 6
			unit = "PB"
		}
		network_outgoing_traffic {
			is_enabled = true
			limit = 7
			unit = "KB"
		}
		network_total_traffic {
			is_enabled = true
			limit = 8
			unit = "MB"
		}
		network_reduce_bandwidth {
			is_enabled = true
			limit = 9
			unit = "Mbps"
		}
		backups_number {
			is_enabled = false
		}
	}

	backup_price = 30
	reset_limit_policy = "never"
	network_traffic_limit_type = "separate"
}
`,
					name,
				),
				Check: resource.ComposeTestCheckFunc(
					checker(
						fmt.Sprintf("solusio_plan.%s_full", name),
						name+"_full_new",
						solus.StorageTypeNameFB,
						solus.ImageFormatRaw,
						solus.PlanParams{
							Disk: 1,
							RAM:  2,
							VCPU: 3,
						},
						40,
						50,
						60,
						70,
						30,
						solus.PlanResetLimitPolicyNever,
						solus.PlanNetworkTotalTrafficTypeSeparate,
					),
					limitsChecker(
						fmt.Sprintf("solusio_plan.%s_full", name),
						solus.PlanLimits{
							DiskIOPS: solus.DiskIOPSPlanLimit{
								IsEnabled: true,
								Limit:     20,
								Unit:      solus.DiskIOPSPlanLimitUnitIOPS,
							},
							DiskBandwidth: solus.DiskBandwidthPlanLimit{
								IsEnabled: true,
								Limit:     30,
								Unit:      solus.DiskBandwidthPlanLimitUnitBps,
							},
							NetworkIncomingBandwidth: solus.BandwidthPlanLimit{
								IsEnabled: true,
								Limit:     4,
								Unit:      solus.BandwidthPlanLimitUnitGbps,
							},
							NetworkOutgoingBandwidth: solus.BandwidthPlanLimit{
								IsEnabled: true,
								Limit:     5,
								Unit:      solus.BandwidthPlanLimitUnitKbps,
							},
							NetworkIncomingTraffic: solus.TrafficPlanLimit{
								IsEnabled: true,
								Limit:     6,
								Unit:      solus.TrafficPlanLimitUnitPB,
							},
							NetworkOutgoingTraffic: solus.TrafficPlanLimit{
								IsEnabled: true,
								Limit:     7,
								Unit:      solus.TrafficPlanLimitUnitKB,
							},
							NetworkTotalTraffic: solus.TrafficPlanLimit{
								IsEnabled: true,
								Limit:     8,
								Unit:      solus.TrafficPlanLimitUnitMB,
							},
							NetworkReduceBandwidth: solus.BandwidthPlanLimit{
								IsEnabled: true,
								Limit:     9,
								Unit:      solus.BandwidthPlanLimitUnitMbps,
							},
							BackupsNumber: solus.UnitPlanLimit{
								IsEnabled: false,
								Limit:     0,
								Unit:      solus.PlanLimitUnits,
							},
						},
					),
				),
			},

			// Update plan's field which should force to create new one.
			{
				Config: fmt.Sprintf(
					`
resource "solusio_plan" "%[1]s_full" {
	name = "%[1]s_full_new"
	storage_type = "fb"
	image_format = "raw"
	params {
		disk = 10
		ram = 2
		vcpu = 3
	}
	tokens_per_hour = 40
	tokens_per_month = 50
	ip_tokens_per_hour = 60
	ip_tokens_per_month = 70

	limits {
		disk_iops {
			is_enabled = true
			limit = 20
		}
		disk_bandwidth {
			is_enabled = true
			limit = 30
		}
		network_incoming_bandwidth {
			is_enabled = true
			limit = 4
			unit = "Gbps"
		}
		network_outgoing_bandwidth {
			is_enabled = true
			limit = 5
			unit = "Kbps"
		}
		network_incoming_traffic {
			is_enabled = true
			limit = 6
			unit = "PB"
		}
		network_outgoing_traffic {
			is_enabled = true
			limit = 7
			unit = "KB"
		}
		network_total_traffic {
			is_enabled = true
			limit = 8
			unit = "MB"
		}
		network_reduce_bandwidth {
			is_enabled = true
			limit = 9
			unit = "Mbps"
		}
		backups_number {
			is_enabled = false
		}
	}

	backup_price = 30
	reset_limit_policy = "never"
	network_traffic_limit_type = "separate"
}
`,
					name,
				),
				Check: resource.ComposeTestCheckFunc(
					checker(
						fmt.Sprintf("solusio_plan.%s_full", name),
						name+"_full_new",
						solus.StorageTypeNameFB,
						solus.ImageFormatRaw,
						solus.PlanParams{
							Disk: 10,
							RAM:  2,
							VCPU: 3,
						},
						40,
						50,
						60,
						70,
						30,
						solus.PlanResetLimitPolicyNever,
						solus.PlanNetworkTotalTrafficTypeSeparate,
					),
					limitsChecker(
						fmt.Sprintf("solusio_plan.%s_full", name),
						solus.PlanLimits{
							DiskIOPS: solus.DiskIOPSPlanLimit{
								IsEnabled: true,
								Limit:     20,
								Unit:      solus.DiskIOPSPlanLimitUnitIOPS,
							},
							DiskBandwidth: solus.DiskBandwidthPlanLimit{
								IsEnabled: true,
								Limit:     30,
								Unit:      solus.DiskBandwidthPlanLimitUnitBps,
							},
							NetworkIncomingBandwidth: solus.BandwidthPlanLimit{
								IsEnabled: true,
								Limit:     4,
								Unit:      solus.BandwidthPlanLimitUnitGbps,
							},
							NetworkOutgoingBandwidth: solus.BandwidthPlanLimit{
								IsEnabled: true,
								Limit:     5,
								Unit:      solus.BandwidthPlanLimitUnitKbps,
							},
							NetworkIncomingTraffic: solus.TrafficPlanLimit{
								IsEnabled: true,
								Limit:     6,
								Unit:      solus.TrafficPlanLimitUnitPB,
							},
							NetworkOutgoingTraffic: solus.TrafficPlanLimit{
								IsEnabled: true,
								Limit:     7,
								Unit:      solus.TrafficPlanLimitUnitKB,
							},
							NetworkTotalTraffic: solus.TrafficPlanLimit{
								IsEnabled: true,
								Limit:     8,
								Unit:      solus.TrafficPlanLimitUnitMB,
							},
							NetworkReduceBandwidth: solus.BandwidthPlanLimit{
								IsEnabled: true,
								Limit:     9,
								Unit:      solus.BandwidthPlanLimitUnitMbps,
							},
							BackupsNumber: solus.UnitPlanLimit{
								IsEnabled: false,
								Limit:     0,
								Unit:      solus.PlanLimitUnits,
							},
						},
					),
				),
			},
		},
	})
}

func testAccCheckPlanDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*solus.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "solusio_plan" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = c.Plans.Get(context.Background(), id)
		if err == nil {
			return fmt.Errorf("plan %d still exists", id)
		}

		if !solus.IsNotFound(err) {
			return err
		}
	}

	return nil
}
