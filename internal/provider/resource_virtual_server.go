package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/solusio/terraform-provider-solus/internal/timer"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/solusio/solus-go-sdk"
)

func resourceVirtualServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: adoptCreate("Virtual VirtualServer", resourceVirtualServerCreate),
		ReadContext:   adoptRead("Virtual VirtualServer", resourceVirtualServerRead),
		UpdateContext: adoptUpdate("Virtual VirtualServer", resourceVirtualServerUpdate),
		DeleteContext: adoptDelete("Virtual VirtualServer", resourceVirtualServerDelete),

		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validationIsDomainName,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"ssh_keys": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntAtLeast(1),
				},
			},
			"plan_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"project_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"location_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"user_data": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues, // todo validate user-data
			},

			// For creating virtual server from OS Image Version.
			"os_image_version_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
				ExactlyOneOf: []string{"os_image_version_id", "application_id"},
			},

			// For creating virtual server from application.
			"application_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
				ExactlyOneOf: []string{"os_image_version_id", "application_id"},
			},
			"application_data": {
				Type:         schema.TypeMap,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
				RequiredWith: []string{"application_id"},
			},
		},
	}
}

func resourceVirtualServerCreate(ctx context.Context, client *client, d *schema.ResourceData) error {
	res, err := client.VirtualServers.Create(ctx, solus.VirtualServerCreateRequest{
		Name:             d.Get("hostname").(string),
		Description:      d.Get("description").(string),
		UserData:         d.Get("user_data").(string),
		SSHKeys:          listOfIDs(d.Get("ssh_keys")),
		PlanID:           d.Get("plan_id").(int),
		ProjectID:        d.Get("project_id").(int),
		LocationID:       d.Get("location_id").(int),
		OSImageVersionID: d.Get("os_image_version_id").(int),
		ApplicationID:    d.Get("application_id").(int),
		ApplicationData:  d.Get("application_data").(map[string]interface{}),
	})
	if err != nil {
		return normalizeAPIError(err)
	}

	if err := resourceVirtualServerWaitFor(ctx, client, res.ID); err != nil {
		return err
	}

	d.SetId(strconv.Itoa(res.ID))
	return resourceVirtualServerRead(ctx, client, d)
}

func resourceVirtualServerRead(ctx context.Context, client *client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	res, err := client.VirtualServers.Get(ctx, id)
	if err != nil {
		return normalizeAPIError(err)
	}

	return newSchemaChainSetter(d).
		SetID(res.ID).
		Set("hostname", res.Name).
		Set("description", res.Description).
		Error()
}

func resourceVirtualServerUpdate(ctx context.Context, client *client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	res, err := client.VirtualServers.Patch(ctx, id, solus.VirtualServerUpdateRequest{
		Name:        d.Get("hostname").(string),
		Description: d.Get("description").(string),
	})
	if err != nil {
		return normalizeAPIError(err)
	}

	d.SetId(strconv.Itoa(res.ID))
	return resourceVirtualServerRead(ctx, client, d)
}

func resourceVirtualServerDelete(ctx context.Context, client *client, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	task, err := client.VirtualServers.Delete(ctx, id)
	if err != nil {
		return normalizeAPIError(err)
	}

	return timer.WaitFor(ctx, 5*time.Second, func() (bool, error) {
		tflog.Trace(ctx, "Wait for Virtual Server %d will be deleted", id)
		t, err := client.Tasks.Get(ctx, task.ID)
		if err != nil {
			return false, normalizeAPIError(err)
		}

		if !t.IsFinished() {
			return false, nil
		}

		if t.Status != solus.TaskStatusDone {
			return false, fmt.Errorf("failed to delete virtual server: %s", t.Output)
		}
		return true, nil
	})
}

func resourceVirtualServerWaitFor(ctx context.Context, client *client, id int) error {
	return timer.WaitFor(ctx, 5*time.Second, func() (bool, error) {
		tflog.Trace(ctx, "Wait for Virtual Server %d will start", id)
		resp, err := client.VirtualServers.Get(ctx, id)
		if err != nil {
			return false, normalizeAPIError(err)
		}

		if resp.IsProcessing {
			return false, nil
		}

		if resp.Status != solus.VirtualServerStatusStarted {
			return false, fmt.Errorf("virtual server didn't started, actual status %q", resp.Status)
		}
		return true, nil
	})
}
