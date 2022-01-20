package solus

import "context"

// ActivityLogsEvent represents an Activity Logs event.
type ActivityLogsEvent string

const (
	ActivityLogsEventAdditionalIPCreateRequested ActivityLogsEvent = "additional_ip_create_requested"
	ActivityLogsEventAdditionalIPCreateFailed    ActivityLogsEvent = "additional_ip_create_failed"
	ActivityLogsEventAdditionalIPCreated         ActivityLogsEvent = "additional_ip_created"

	ActivityLogsEventAdditionalIPDeleteRequested ActivityLogsEvent = "additional_ip_delete_requested"
	ActivityLogsEventAdditionalIPDeleteFailed    ActivityLogsEvent = "additional_ip_delete_failed"
	ActivityLogsEventAdditionalIPDeleted         ActivityLogsEvent = "additional_ip_deleted"

	ActivityLogsEventBackupsDeleted ActivityLogsEvent = "backups_deleted"

	ActivityLogsEventComputeResourceCreateRequested ActivityLogsEvent = "compute_resource_create_requested"
	ActivityLogsEventComputeResourceDeleteRequested ActivityLogsEvent = "compute_resource_delete_requested"

	ActivityLogsEventComputeResourceVMCreateRequested ActivityLogsEvent = "compute_resource_vm_create_requested"
	ActivityLogsEventComputeResourceVMCreateFailed    ActivityLogsEvent = "compute_resource_vm_create_failed"
	ActivityLogsEventComputeResourceVMCreated         ActivityLogsEvent = "compute_resource_vm_created"

	ActivityLogsEventComputeResourceVMDeleteRequested ActivityLogsEvent = "compute_resource_vm_delete_requested"

	ActivityLogsEventComputeResourceVMBatchDeleted ActivityLogsEvent = "compute_resource_vm_batch_deleted"

	ActivityLogsEventExternalIntegrationRequested ActivityLogsEvent = "external_integration_requested"

	ActivityLogsEventFailedToSendEmailNotification ActivityLogsEvent = "failed_to_send_email_notification"

	ActivityLogsEventIPBlockChanged ActivityLogsEvent = "ip_block_changed"

	ActivityLogsEventLocationCreated ActivityLogsEvent = "location_created"
	ActivityLogsEventLocationDeleted ActivityLogsEvent = "location_deleted"
	ActivityLogsEventLocationChanged ActivityLogsEvent = "location_changed"

	ActivityLogsEventUserCreateRequested ActivityLogsEvent = "user_create_requested"
	ActivityLogsEventUserCreateFailed    ActivityLogsEvent = "user_create_failed"
	ActivityLogsEventUserCreated         ActivityLogsEvent = "user_created"

	ActivityLogsEventUserDeleteRequested ActivityLogsEvent = "user_delete_requested"
	ActivityLogsEventUserDeleteFailed    ActivityLogsEvent = "user_delete_failed"
	ActivityLogsEventUserDeleted         ActivityLogsEvent = "user_deleted"
)

// ActivityLogsService provides access to Activity Logs.
type ActivityLogsService service

// ActivityLogs represents an Activity Logs.
type ActivityLogs struct {
	ID        int               `json:"id"`
	Event     ActivityLogsEvent `json:"event"`
	UserEmail string            `json:"user_email"`
	CreatedAt string            `json:"created_at"`
	Data      interface{}       `json:"data"`
}

// ActivityLogsResponse represents an Activity Logs response.
type ActivityLogsResponse struct {
	paginatedResponse

	Data []ActivityLogs `json:"data"`
}

// FilterActivityLogs represent available filters for fetching list of Activity Logs.
type FilterActivityLogs struct {
	filter
}

// ByUserID filter Activity Logs by specified User ID.
func (f *FilterActivityLogs) ByUserID(id int) *FilterActivityLogs {
	f.addInt("filter[user_id]", id)
	return f
}

// ByEvent filter Activity Logs by event type.
func (f *FilterActivityLogs) ByEvent(event ActivityLogsEvent) *FilterActivityLogs {
	f.add("filter[event]", string(event))
	return f
}

// List return list of activity logs, filter can be nil.
func (s *ActivityLogsService) List(ctx context.Context, filter *FilterActivityLogs) (ActivityLogsResponse, error) {
	resp := ActivityLogsResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "activity_logs", &resp, withFilter(filter.data))
}
