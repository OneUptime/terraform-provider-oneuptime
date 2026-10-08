---
page_title: "oneuptime_monitor Data Source - oneuptime"
subcategory: "Monitors"
description: |-
  Monitor is anything that monitors your API, Websites, IP, Network or more. You can also create static monitor that does not monitor anything.
---

# oneuptime_monitor (Data Source)

Monitor is anything that monitors your API, Websites, IP, Network or more. You can also create static monitor that does not monitor anything.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one monitor may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_monitor" "example" {
  name = "Example monitor"
}

# Or by id:
data "oneuptime_monitor" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `auto_provisioned_network_device_id` (String) ID of the Network Device that caused this monitor to be provisioned automatically. The ID of a `oneuptime_network_device`.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `current_monitor_status_id` (String) Whats the current status ID of this monitor? The ID of a `oneuptime_monitor_status`.
- `description` (String) Friendly description that will help you remember.
- `disable_active_monitoring` (Boolean) Disable active monitoring for this resource?
- `disable_active_monitoring_because_of_manual_incident` (Boolean) Disable Monitoring because of Incident which is creeated manually by user.
- `disable_active_monitoring_because_of_scheduled_maintenance_event` (Boolean) Disable Monitoring because of Ongoing Scheduled Maintenance Event.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incoming_email_custom_local_part` (String) This field is for Incoming Email Monitor only. A custom name for this monitor's inbound email address: the part before the @, on the server's inbound email domain. While set, it replaces the generated monitor-{secret key} address. Must be unique across all monitors. Set to null to go back to the generated address.
- `incoming_email_secret_key` (String) This field is for Incoming Email Monitor only. Secret Key used to generate unique email address.
- `incoming_request_secret_key` (String) This field is for Incoming Request Monitor only. Secret Key to authenticate the request.
- `is_all_probes_disconnected_from_this_monitor` (Boolean) All Probes Disconnected From This Monitor. Is this monitor not being monitored?
- `is_archived` (Boolean) Archived monitors are hidden from monitor lists and status pages, are not checked, and open no incidents or alerts. Unarchiving resumes monitoring.
- `is_no_probe_enabled_on_this_monitor` (Boolean) No Probe Enabled On This Monitor. Is this monitor not being monitored?
- `is_owner_notified_of_resource_creation` (Boolean) Are owners notified of when this resource is created?
- `minimum_probe_agreement` (Number) Minimum number of probes that must agree on a status before the monitor status changes. If null, all enabled and connected probes must agree.
- `monitor_template_id` (String) ID of the Monitor Template this monitor was created from. Null for monitors not created from a template. The ID of a `oneuptime_monitor_template`.
- `monitor_type` (String) What is the type of this monitor? Website? API? etc.
- `monitoring_interval` (String) How often would you like OneUptime to monitor this resource? A 5-field cron expression, not a label: "*/5 * * * *" is every five minutes.
- `name` (String) Any friendly name for this monitor.
- `network_alert_policy_id` (String) ID of the Network Alert Policy that provisioned this monitor, when one did. The ID of a `oneuptime_network_alert_policy`.
- `server_monitor_secret_key` (String) This field is for Server Monitor only. Secret Key to authenticate the request.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `archived_at` (String) When this monitor was archived. Empty while it is not archived.
- `ceph_clusters` (Set of String) Ceph clusters this monitor watches. Incidents and alerts it creates are linked to them, so OneUptime AI can investigate and fix them there. IDs of `oneuptime_ceph_cluster` resources.
- `created_at` (String) Date and Time when the object was created.
- `custom_fields` (String) Custom Fields on this resource. A JSON value: write it with `jsonencode()`.
- `database_servers` (Set of String) Databases this monitor watches. Incidents and alerts it creates are linked to them, so OneUptime AI can investigate and fix them there. IDs of `oneuptime_database` resources.
- `depends_on_monitors` (Set of String) Parent monitors this monitor depends on. When a parent is offline (or in one of the configured suppression statuses), alerts and incidents from this monitor are suppressed at creation time — the monitor keeps evaluating and its status timeline still updates. IDs of `oneuptime_monitor` resources.
- `docker_hosts` (Set of String) Docker hosts this monitor watches. Incidents and alerts it creates are linked to them, so OneUptime AI can investigate and fix them there. IDs of `oneuptime_docker_host` resources.
- `docker_swarm_clusters` (Set of String) Docker Swarm clusters this monitor watches. Incidents and alerts it creates are linked to them, so OneUptime AI can investigate and fix them there. IDs of `oneuptime_docker_swarm_cluster` resources.
- `hosts` (Set of String) Hosts this monitor watches. Incidents and alerts it creates are linked to them, so OneUptime AI can investigate and fix them there. IDs of `oneuptime_host` resources.
- `incoming_email_monitor_heartbeat_checked_at` (String) This field is for Incoming Email monitor only. When was the last time we checked the heartbeat?
- `incoming_email_monitor_last_email_received_at` (String) This field is for Incoming Email Monitor only. When was the last email received?
- `incoming_email_monitor_request` (String) This field is for Incoming Email Monitor only. Last email data received. A JSON value: write it with `jsonencode()`.
- `incoming_monitor_request` (String) Incoming Monitor Request for Incoming Request Monitor. A JSON value: write it with `jsonencode()`.
- `incoming_request_monitor_heartbeat_checked_at` (String) This field is for Incoming Request monitor only. When was the last time we checked the heartbeat?
- `iot_fleets` (Set of String) IoT fleets this monitor watches. Incidents and alerts it creates are linked to them, so OneUptime AI can investigate and fix them there. IDs of `oneuptime_iot_fleet` resources.
- `kubernetes_clusters` (Set of String) Kubernetes clusters this monitor watches. Incidents and alerts it creates are linked to them, so OneUptime AI can investigate and fix them there. IDs of `oneuptime_kubernetes_cluster` resources.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `monitor_steps` (String) What would you like to monitor and what is the criteria?
- `podman_hosts` (Set of String) Podman hosts this monitor watches. Incidents and alerts it creates are linked to them, so OneUptime AI can investigate and fix them there. IDs of `oneuptime_podman_host` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `proxmox_clusters` (Set of String) Proxmox clusters this monitor watches. Incidents and alerts it creates are linked to them, so OneUptime AI can investigate and fix them there. IDs of `oneuptime_proxmox_cluster` resources.
- `server_monitor_request_received_at` (String) This field is for Server Monitor only. When was the last time we received a request?
- `server_monitor_response` (String) Server Monitor Response for Server Monitor. A JSON value: write it with `jsonencode()`.
- `services` (Set of String) Services this monitor watches. Incidents and alerts it creates are linked to them, so OneUptime AI can investigate and fix them there. IDs of `oneuptime_service` resources.
- `storage_arrays` (Set of String) Storage arrays this monitor watches. Incidents and alerts it creates are linked to them. IDs of `oneuptime_storage_array` resources.
- `suppress_alerts_when_parent_monitor_statuses` (Set of String) Parent monitor statuses that suppress this monitor's alerts and incidents. When empty, statuses flagged as offline suppress (the default). Only used when Depends On Monitors is set. IDs of `oneuptime_monitor_status` resources.
- `telemetry_monitor_last_monitor_at` (String) This field is for Telemetry Monitor only. When was the last time we monitored?
- `telemetry_monitor_next_monitor_at` (String) This field is for Telemetry Monitor only. When is the next time we should monitor?
- `updated_at` (String) Date and Time when the object was updated.
- `vmware_v_centers` (Set of String) VMware vCenters this monitor watches. Incidents and alerts it creates are linked to them, so OneUptime AI can investigate and fix them there. IDs of `oneuptime_vcenter` resources.
