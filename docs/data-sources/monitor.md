---
page_title: "oneuptime_monitor Data Source - oneuptime"
subcategory: "Monitors"
description: |-
  Monitor is anything that monitors your API, Websites, IP, Network or more. You can also create static monitor that does not monitor anything.
---

# oneuptime_monitor (Data Source)

Monitor is anything that monitors your API, Websites, IP, Network or more. You can also create static monitor that does not monitor anything. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_monitor" "by_name" {
  name = "example-monitor"
}

data "oneuptime_monitor" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description` (String) Friendly description that will help you remember.. Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `depends_on_monitors` (Set) Parent monitors this monitor depends on. When a parent is offline (or in one of the configured suppression statuses), alerts and incidents from this monitor are suppressed at creation time — the monitor keeps evaluating and its status timeline still updates... Computed.
- `suppress_alerts_when_parent_monitor_statuses` (Set) Parent monitor statuses that suppress this monitor's alerts and incidents. When empty, statuses flagged as offline suppress (the default). Only used when Depends On Monitors is set... Computed.
- `monitor_template_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `auto_provisioned_network_device_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `monitor_type` (String) What is the type of this monitor? Website? API? etc... Computed.
- `current_monitor_status_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `monitor_steps` (Monitor_steps) MonitorSteps object. Computed.
- `monitoring_interval` (String) How often would you like OneUptime to monitor this resource?.. Computed.
- `custom_fields` (String) Custom Fields on this resource... Computed.
- `is_owner_notified_of_resource_creation` (Bool) Are owners notified of when this resource is created?.. Computed.
- `disable_active_monitoring` (Bool) Disable active monitoring for this resource?.. Computed.
- `incoming_request_monitor_heartbeat_checked_at` (String) A date time object.. Computed.
- `telemetry_monitor_next_monitor_at` (String) A date time object.. Computed.
- `telemetry_monitor_last_monitor_at` (String) A date time object.. Computed.
- `disable_active_monitoring_because_of_scheduled_maintenance_event` (Bool) Disable Monitoring because of Ongoing Scheduled Maintenance Event.. Computed.
- `disable_active_monitoring_because_of_manual_incident` (Bool) Disable Monitoring because of Incident which is creeated manually by user... Computed.
- `server_monitor_request_received_at` (String) A date time object.. Computed.
- `server_monitor_secret_key` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `incoming_request_secret_key` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `incoming_monitor_request` (String) Incoming Monitor Request for Incoming Request Monitor.. Computed.
- `incoming_email_secret_key` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `incoming_email_monitor_last_email_received_at` (String) A date time object.. Computed.
- `incoming_email_monitor_request` (String) This field is for Incoming Email Monitor only. Last email data received... Computed.
- `incoming_email_monitor_heartbeat_checked_at` (String) A date time object.. Computed.
- `server_monitor_response` (String) Server Monitor Response for Server Monitor.. Computed.
- `is_all_probes_disconnected_from_this_monitor` (Bool) All Probes Disconnected From This Monitor. Is this monitor not being monitored?.. Computed.
- `is_no_probe_enabled_on_this_monitor` (Bool) No Probe Enabled On This Monitor. Is this monitor not being monitored?.. Computed.
- `minimum_probe_agreement` (Number) Minimum number of probes that must agree on a status before the monitor status changes. If null, all enabled and connected probes must agree... Computed.
