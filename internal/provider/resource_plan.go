package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func resourcePlan() *schema.Resource {
	createPlanLimitResource := func(t interface{}, forceNew bool) *schema.Schema {
		ss := map[interface{}][]string{
			solus.DiskBandwidthPlanLimit{}: {
				string(solus.DiskBandwidthPlanLimitUnitBps),
			},
			solus.DiskIOPSPlanLimit{}: {
				string(solus.DiskIOPSPlanLimitUnitIOPS),
			},
			solus.BandwidthPlanLimit{}: {
				string(solus.BandwidthPlanLimitUnitKbps),
				string(solus.BandwidthPlanLimitUnitMbps),
				string(solus.BandwidthPlanLimitUnitGbps),
			},
			solus.TrafficPlanLimit{}: {
				string(solus.TrafficPlanLimitUnitKB),
				string(solus.TrafficPlanLimitUnitMB),
				string(solus.TrafficPlanLimitUnitGB),
				string(solus.TrafficPlanLimitUnitTB),
				string(solus.TrafficPlanLimitUnitPB),
			},
			solus.UnitPlanLimit{}: {
				string(solus.PlanLimitUnits),
			},
		}

		s, ok := ss[t]
		if !ok {
			panic(fmt.Sprintf("unhandled limit type %T", t))
		}

		return &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			ForceNew: forceNew,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"is_enabled": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},

					"limit": {
						Type:     schema.TypeInt,
						Optional: true,
						Default:  0,
					},

					"unit": {
						Type:         schema.TypeString,
						Optional:     true,
						Default:      defaultLimitUnit(t),
						ValidateFunc: validation.StringInSlice(s, false),
					},
				},
			},
		}
	}

	return &schema.Resource{
		CreateContext: adoptCreate("Plan", resourcePlanCreate),
		ReadContext:   adoptRead("Plan", resourcePlanRead),
		UpdateContext: adoptUpdate("Plan", resourcePlanUpdate),
		DeleteContext: adoptDelete("Plan", resourcePlanDelete),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"virtualization_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validationIsVirtualizationType,
			},

			"storage_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(solus.StorageTypeNameFB),
					string(solus.StorageTypeNameLVM),
					string(solus.StorageTypeNameThinLVM),
					string(solus.StorageTypeNameNFS),
					string(solus.StorageTypeNameVZ),
				}, false),
			},

			"image_format": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// Actually list of valid types for this property depends on
				// `storage_type` value, but we didn't have an ability to make
				// proper validation here.
				ValidateFunc: validation.StringInSlice([]string{
					string(solus.ImageFormatRaw),
					string(solus.ImageFormatQCOW2),
					string(solus.ImageFormatPLOOP),
				}, false),
			},

			"params": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk": {
							Type:         schema.TypeInt,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
						"ram_mb": {
							Type:         schema.TypeInt,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntAtLeast(128),
						},
						"vcpu": {
							Type:         schema.TypeInt,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntBetween(1, 64),
						},

						// Required for VZ plans.
						"vcpu_units": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntBetween(8, 500000),
						},
						"vcpu_limit": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntBetween(10, 6400),
						},
						"io_priority": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntBetween(0, 7),
						},
					},
				},
			},

			"is_default": {
				Type:         schema.TypeBool,
				Optional:     true,
				Default:      false,
				ValidateFunc: validation.NoZeroValues,
			},

			"tokens_per_hour": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 1000),
			},

			"tokens_per_month": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 1000),
			},

			"ip_tokens_per_hour": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 1000),
			},

			"ip_tokens_per_month": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 1000),
			},

			"is_additional_ips_available": {
				Type:         schema.TypeBool,
				Optional:     true,
				Default:      false,
				ValidateFunc: validation.NoZeroValues,
			},

			"limits": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_iops":                  createPlanLimitResource(solus.DiskIOPSPlanLimit{}, true),
						"disk_bandwidth":             createPlanLimitResource(solus.DiskBandwidthPlanLimit{}, true),
						"network_incoming_bandwidth": createPlanLimitResource(solus.BandwidthPlanLimit{}, false),
						"network_outgoing_bandwidth": createPlanLimitResource(solus.BandwidthPlanLimit{}, false),
						"network_incoming_traffic":   createPlanLimitResource(solus.TrafficPlanLimit{}, false),
						"network_outgoing_traffic":   createPlanLimitResource(solus.TrafficPlanLimit{}, false),
						"network_total_traffic":      createPlanLimitResource(solus.TrafficPlanLimit{}, false),
						"network_reduce_bandwidth":   createPlanLimitResource(solus.BandwidthPlanLimit{}, false),
						"backups_number":             createPlanLimitResource(solus.UnitPlanLimit{}, false),
					},
				},
			},

			"is_snapshots_enabled": {
				Type:         schema.TypeBool,
				Optional:     true,
				Default:      false,
				ValidateFunc: validation.NoZeroValues,
			},

			"is_visible": {
				Type:         schema.TypeBool,
				Optional:     true,
				Default:      true,
				ValidateFunc: validation.NoZeroValues,
			},

			"is_backup_available": {
				Type:         schema.TypeBool,
				Optional:     true,
				Default:      false,
				ValidateFunc: validation.NoZeroValues,
			},

			"backup_price": {
				Type:         schema.TypeFloat,
				Optional:     true,
				Default:      0.0,
				ValidateFunc: validation.FloatBetween(0, 100),
			},

			"reset_limit_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(solus.PlanResetLimitPolicyNever),
				ValidateFunc: validation.StringInSlice([]string{
					string(solus.PlanResetLimitPolicyNever),
					string(solus.PlanResetLimitPolicyFirstDayOfMonth),
					string(solus.PlanResetLimitPolicyVMCreatedDay),
				}, false),
			},

			"network_traffic_limit_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(solus.PlanNetworkTotalTrafficTypeSeparate),
				ValidateFunc: validation.StringInSlice([]string{
					string(solus.PlanNetworkTotalTrafficTypeSeparate),
					string(solus.PlanNetworkTotalTrafficTypeTotal),
				}, false),
			},
		},
	}
}

func resourcePlanCreate(ctx context.Context, client *solus.Client, d *schema.ResourceData) error {
	res, err := client.Plans.Create(ctx, solus.PlanCreateRequest{
		Name:                     d.Get("name").(string),
		VirtualizationType:       solus.VirtualizationType(d.Get("virtualization_type").(string)),
		Params:                   resourceToPlanParams(d.Get("params")),
		StorageType:              solus.StorageTypeName(d.Get("storage_type").(string)),
		ImageFormat:              solus.ImageFormat(d.Get("image_format").(string)),
		Limits:                   resourceToPlanLimits(d.Get("limits")),
		TokensPerHour:            d.Get("tokens_per_hour").(int),
		TokensPerMonth:           d.Get("tokens_per_month").(int),
		IPTokensPerHour:          d.Get("ip_tokens_per_hour").(int),
		IPTokensPerMonth:         d.Get("ip_tokens_per_month").(int),
		IsVisible:                d.Get("is_visible").(bool),
		IsDefault:                d.Get("is_default").(bool),
		IsSnapshotsEnabled:       d.Get("is_snapshots_enabled").(bool),
		IsBackupAvailable:        d.Get("is_backup_available").(bool),
		IsAdditionalIPsAvailable: d.Get("is_additional_ips_available").(bool),
		BackupPrice:              d.Get("backup_price").(float64),
		ResetLimitPolicy:         solus.PlanResetLimitPolicy(d.Get("reset_limit_policy").(string)),
		NetworkTotalTrafficType:  solus.PlanNetworkTotalTrafficType(d.Get("network_traffic_limit_type").(string)),
		BackupSettings: solus.PlanBackupSettings{
			IsIncrementalBackupEnabled: false,
			IncrementalBackupsLimit:    1,
		},
	})
	if err != nil {
		return normalizeAPIError(err)
	}

	d.SetId(strconv.Itoa(res.ID))
	return resourcePlanRead(ctx, client, d)
}

func resourceToPlanLimits(i interface{}) solus.PlanLimits {
	getDefaultUnit := func(k string) string {
		switch k {
		case "disk_bandwidth":
			return string(solus.DiskBandwidthPlanLimitUnitBps)

		case "disk_iops":
			return string(solus.DiskIOPSPlanLimitUnitIOPS)

		case "network_incoming_bandwidth", "network_outgoing_bandwidth", "network_reduce_bandwidth":
			return string(solus.BandwidthPlanLimitUnitMbps)

		case "network_incoming_traffic", "network_outgoing_traffic", "network_total_traffic":
			return string(solus.TrafficPlanLimitUnitGB)

		case "backups_number":
			return string(solus.PlanLimitUnits)
		}
		panic(fmt.Sprintf("unknown key %q", k))
	}

	getKeyValue := func(m map[string]interface{}, k string) map[string]interface{} {
		defaultUnit := getDefaultUnit(k)
		v, ok := m[k]
		if !ok {
			return map[string]interface{}{
				"is_enabled": false,
				"limit":      0,
				"unit":       defaultUnit,
			}
		}

		iv, ok := v.([]interface{})
		if !ok || len(iv) == 0 {
			return map[string]interface{}{
				"is_enabled": false,
				"limit":      0,
				"unit":       defaultUnit,
			}
		}

		return iv[0].(map[string]interface{})
	}

	mm := i.([]interface{}) //nolint:errcheck // Not necessary.
	if len(mm) == 0 {
		mm = []interface{}{
			map[string]interface{}{},
		}
	}

	m := mm[0].(map[string]interface{}) //nolint:errcheck // Not necessary.

	dband := getKeyValue(m, "disk_bandwidth")
	diops := getKeyValue(m, "disk_iops")
	niband := getKeyValue(m, "network_incoming_bandwidth")
	noband := getKeyValue(m, "network_outgoing_bandwidth")
	nitraff := getKeyValue(m, "network_incoming_traffic")
	notraff := getKeyValue(m, "network_outgoing_traffic")
	nttraff := getKeyValue(m, "network_total_traffic")
	nrband := getKeyValue(m, "network_reduce_bandwidth")
	bnumb := getKeyValue(m, "backups_number")

	return solus.PlanLimits{
		DiskBandwidth: solus.DiskBandwidthPlanLimit{
			IsEnabled: dband["is_enabled"].(bool),
			Limit:     dband["limit"].(int),
			Unit:      solus.DiskBandwidthPlanLimitUnit(dband["unit"].(string)),
		},
		DiskIOPS: solus.DiskIOPSPlanLimit{
			IsEnabled: diops["is_enabled"].(bool),
			Limit:     diops["limit"].(int),
			Unit:      solus.DiskIOPSPlanLimitUnit(diops["unit"].(string)),
		},
		NetworkIncomingBandwidth: solus.BandwidthPlanLimit{
			IsEnabled: niband["is_enabled"].(bool),
			Limit:     niband["limit"].(int),
			Unit:      solus.BandwidthPlanLimitUnit(niband["unit"].(string)),
		},
		NetworkOutgoingBandwidth: solus.BandwidthPlanLimit{
			IsEnabled: noband["is_enabled"].(bool),
			Limit:     noband["limit"].(int),
			Unit:      solus.BandwidthPlanLimitUnit(noband["unit"].(string)),
		},
		NetworkIncomingTraffic: solus.TrafficPlanLimit{
			IsEnabled: nitraff["is_enabled"].(bool),
			Limit:     nitraff["limit"].(int),
			Unit:      solus.TrafficPlanLimitUnit(nitraff["unit"].(string)),
		},
		NetworkOutgoingTraffic: solus.TrafficPlanLimit{
			IsEnabled: notraff["is_enabled"].(bool),
			Limit:     notraff["limit"].(int),
			Unit:      solus.TrafficPlanLimitUnit(notraff["unit"].(string)),
		},
		NetworkTotalTraffic: solus.TrafficPlanLimit{
			IsEnabled: nttraff["is_enabled"].(bool),
			Limit:     nttraff["limit"].(int),
			Unit:      solus.TrafficPlanLimitUnit(nttraff["unit"].(string)),
		},
		NetworkReduceBandwidth: solus.BandwidthPlanLimit{
			IsEnabled: nrband["is_enabled"].(bool),
			Limit:     nrband["limit"].(int),
			Unit:      solus.BandwidthPlanLimitUnit(nrband["unit"].(string)),
		},
		BackupsNumber: solus.UnitPlanLimit{
			IsEnabled: bnumb["is_enabled"].(bool),
			Limit:     bnumb["limit"].(int),
			Unit:      solus.PlanLimitUnit(bnumb["unit"].(string)),
		},
	}
}

func resourcePlanRead(ctx context.Context, client *solus.Client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	res, err := client.Plans.Get(ctx, id)
	if err != nil {
		return normalizeAPIError(err)
	}

	return newSchemaChainSetter(d).
		SetID(res.ID).
		Set("name", res.Name).
		Set("virtualization_type", res.VirtualizationType).
		Set("params", planParamsToResource(res.Params)).
		Set("storage_type", res.StorageType).
		Set("image_format", res.ImageFormat).
		Set("tokens_per_hour", res.TokensPerHour).
		Set("tokens_per_month", res.TokensPerMonth).
		Set("ip_tokens_per_hour", res.IPTokensPerHour).
		Set("ip_tokens_per_month", res.IPTokensPerMonth).
		Set("is_visible", res.IsVisible).
		Set("is_default", res.IsDefault).
		Set("is_snapshots_enabled", res.IsSnapshotsEnabled).
		Set("is_backup_available", res.IsBackupAvailable).
		Set("is_additional_ips_available", res.IsAdditionalIPsAvailable).
		Set("backup_price", res.BackupPrice).
		Set("reset_limit_policy", res.ResetLimitPolicy).
		Set("network_traffic_limit_type", res.NetworkTotalTrafficType).
		Error()
}

func resourcePlanUpdate(ctx context.Context, client *solus.Client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	res, err := client.Plans.Update(ctx, id, solus.PlanUpdateRequest{
		Name:                     d.Get("name").(string),
		Limits:                   resourceToPlanUpdateLimits(d.Get("limits")),
		TokensPerHour:            d.Get("tokens_per_hour").(int),
		TokensPerMonth:           d.Get("tokens_per_month").(int),
		IPTokensPerHour:          d.Get("ip_tokens_per_hour").(int),
		IPTokensPerMonth:         d.Get("ip_tokens_per_month").(int),
		IsVisible:                d.Get("is_visible").(bool),
		IsDefault:                d.Get("is_default").(bool),
		IsSnapshotsEnabled:       d.Get("is_snapshots_enabled").(bool),
		IsBackupAvailable:        d.Get("is_backup_available").(bool),
		IsAdditionalIPsAvailable: d.Get("is_additional_ips_available").(bool),
		BackupPrice:              d.Get("backup_price").(float64),
		ResetLimitPolicy:         solus.PlanResetLimitPolicy(d.Get("reset_limit_policy").(string)),
		NetworkTotalTrafficType:  solus.PlanNetworkTotalTrafficType(d.Get("network_traffic_limit_type").(string)),
		BackupSettings: solus.PlanBackupSettings{
			IsIncrementalBackupEnabled: false,
			IncrementalBackupsLimit:    1,
		},
	})
	if err != nil {
		return normalizeAPIError(err)
	}

	d.SetId(strconv.Itoa(res.ID))
	return resourcePlanRead(ctx, client, d)
}

func resourceToPlanUpdateLimits(i interface{}) solus.PlanUpdateLimits {
	mm := i.([]interface{})             //nolint:errcheck // Not necessary.
	m := mm[0].(map[string]interface{}) //nolint:errcheck // Not necessary.

	getValue := func(k string) map[string]interface{} {
		return (m[k].([]interface{}))[0].(map[string]interface{})
	}

	niband := getValue("network_incoming_bandwidth")
	noband := getValue("network_outgoing_bandwidth")
	nitraff := getValue("network_incoming_traffic")
	notraff := getValue("network_outgoing_traffic")
	nttraff := getValue("network_total_traffic")
	nrband := getValue("network_reduce_bandwidth")
	bnumb := getValue("backups_number")

	return solus.PlanUpdateLimits{
		NetworkIncomingBandwidth: solus.BandwidthPlanLimit{
			IsEnabled: niband["is_enabled"].(bool),
			Limit:     niband["limit"].(int),
			Unit:      solus.BandwidthPlanLimitUnit(niband["unit"].(string)),
		},
		NetworkOutgoingBandwidth: solus.BandwidthPlanLimit{
			IsEnabled: noband["is_enabled"].(bool),
			Limit:     noband["limit"].(int),
			Unit:      solus.BandwidthPlanLimitUnit(noband["unit"].(string)),
		},
		NetworkIncomingTraffic: solus.TrafficPlanLimit{
			IsEnabled: nitraff["is_enabled"].(bool),
			Limit:     nitraff["limit"].(int),
			Unit:      solus.TrafficPlanLimitUnit(nitraff["unit"].(string)),
		},
		NetworkOutgoingTraffic: solus.TrafficPlanLimit{
			IsEnabled: notraff["is_enabled"].(bool),
			Limit:     notraff["limit"].(int),
			Unit:      solus.TrafficPlanLimitUnit(notraff["unit"].(string)),
		},
		NetworkTotalTraffic: solus.TrafficPlanLimit{
			IsEnabled: nttraff["is_enabled"].(bool),
			Limit:     nttraff["limit"].(int),
			Unit:      solus.TrafficPlanLimitUnit(nttraff["unit"].(string)),
		},
		NetworkReduceBandwidth: solus.BandwidthPlanLimit{
			IsEnabled: nrband["is_enabled"].(bool),
			Limit:     nrband["limit"].(int),
			Unit:      solus.BandwidthPlanLimitUnit(nrband["unit"].(string)),
		},
		BackupsNumber: solus.UnitPlanLimit{
			IsEnabled: bnumb["is_enabled"].(bool),
			Limit:     bnumb["limit"].(int),
			Unit:      solus.PlanLimitUnit(bnumb["unit"].(string)),
		},
	}
}

func resourcePlanDelete(ctx context.Context, client *solus.Client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	return normalizeAPIError(client.Plans.Delete(ctx, id))
}

func resourceToPlanParams(i interface{}) solus.PlanParams {
	mm := i.([]interface{})             //nolint:errcheck // Not necessary.
	m := mm[0].(map[string]interface{}) //nolint:errcheck // Not necessary.

	return solus.PlanParams{
		Disk: m["disk"].(int),
		// Convert megabytes to bytes.
		RAM:        m["ram_mb"].(int) * Mb,
		VCPU:       m["vcpu"].(int),
		VCPUUnits:  m["vcpu_units"].(int),
		VCPULimit:  m["vcpu_limit"].(int),
		IOPriority: m["io_priority"].(int),
	}
}

func planParamsToResource(p solus.PlanParams) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"disk": p.Disk,
			// Convert bytes to megabytes.
			"ram_mb":      p.RAM / Mb,
			"vcpu":        p.VCPU,
			"vcpu_units":  p.VCPUUnits,
			"vcpu_limit":  p.VCPULimit,
			"io_priority": p.IOPriority,
		},
	}
}

func defaultLimitUnit(t interface{}) string {
	dd := map[interface{}]string{
		solus.DiskBandwidthPlanLimit{}: string(solus.DiskBandwidthPlanLimitUnitBps),
		solus.DiskIOPSPlanLimit{}:      string(solus.DiskIOPSPlanLimitUnitIOPS),
		solus.BandwidthPlanLimit{}:     string(solus.BandwidthPlanLimitUnitMbps),
		solus.TrafficPlanLimit{}:       string(solus.TrafficPlanLimitUnitGB),
		solus.UnitPlanLimit{}:          string(solus.PlanLimitUnits),
	}

	d, ok := dd[t]
	if !ok {
		panic(fmt.Sprintf("unhandled limit type %T", t))
	}

	return d
}
