---
page_title: "oneuptime_alert Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Manage alerts for your project
---

# oneuptime_alert (Resource)

Manage alerts for your project

## Example Usage

```terraform
resource "oneuptime_alert" "example" {
  title = "This is an example of longer text content that might be stored in this field."
  alert_severity_id = "123e4567-e89b-12d3-a456-426614174000"
  description = "# Heading

This is **markdown** content"
}
```

## Schema

### Required

- `title` (String) Title of this alert..
- `alert_severity_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Short description of this alert. This will be visible on the status page. This is in markdown...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `monitor_id` (String) A unique identifier for an object, represented as a UUID..
- `on_call_duty_policies` (Set) List of on-call duty policies affected by this alert...
- `hosts` (Set) List of hosts affected by this alert...
- `kubernetes_clusters` (Set) List of Kubernetes clusters affected by this alert...
- `kubernetes_resources` (Set) List of Kubernetes resources (pods, deployments, nodes, etc.) affected by this alert...
- `kubernetes_containers` (Set) List of Kubernetes containers affected by this alert...
- `docker_hosts` (Set) List of Docker hosts affected by this alert...
- `podman_hosts` (Set) List of Podman hosts affected by this alert...
- `proxmox_clusters` (Set) List of Proxmox clusters affected by this alert...
- `iot_fleets` (Set) List of IoT fleets affected by this alert...
- `docker_swarm_clusters` (Set) List of Docker Swarm clusters affected by this alert...
- `ceph_clusters` (Set) List of Ceph clusters affected by this alert...
- `docker_resources` (Set) List of Docker resources (containers, images, networks, volumes) affected by this alert...
- `podman_resources` (Set) List of Podman resources (containers, images, networks, volumes) affected by this alert...
- `services` (Set) List of services affected by this alert...
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `current_alert_state_id` (String) A unique identifier for an object, represented as a UUID..
- `monitor_status_when_this_alert_was_created_id` (String) A unique identifier for an object, represented as a UUID..
- `custom_fields` (String) Custom Fields on this resource...
- `root_cause` (String) What is the root cause of this alert?..
- `remediation_notes` (String) Notes on how to remediate this alert. This is in markdown...
- `telemetry_query` (String) Telemetry query for this alert..
- `alert_episode_id` (String) A unique identifier for an object, represented as a UUID..
- `is_private` (Bool) If true, this alert is only visible to its owners (users in 'owner users' and members of 'owner teams'), project admins, and project owners...
- `enable_reminders` (Bool) Should reminder notifications be sent to owners while this alert is still open? Reminders are sent based on the reminder rules configured for this project...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `is_owner_notified_of_alert_creation` (Bool) Are owners notified of when this alert is created?..
- `created_state_log` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert], Update: [No access - you don't have permission for this operation].
- `created_criteria_id` (String) If this alert was created by a Probe, this is the ID of the criteria that created it...
- `series_fingerprint` (String) For metric monitors with per-series alerting (e.g. grouped by host.name), this is a stable hash of the series label values so one alert is created per affected series...
- `series_labels` (String) Attribute key/value pairs that identify the affected series (e.g. {host.name: prod-db-01}) when this alert was created from a per-series metric breach...
- `monitor_summary` (String) The monitor summary captured at the moment this alert was created - the same card the monitor page shows, frozen so it survives the monitor log being aged out...
- `created_by_probe_id` (String) A unique identifier for an object, represented as a UUID..
- `is_created_automatically` (Bool) Is this alert created by OneUptime Probe or Workers automatically (and not created manually by a user)?..
- `alert_number` (Number) Alert Number..
- `alert_number_with_prefix` (String) Alert number with prefix (e.g., 'ALT-42' or '#42')..
- `next_reminder_notification_at` (String) A date time object..
- `reminder_notification_sent_count` (Number) How many reminder notifications have been sent to owners of this alert so far...

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_alert.example <id>
```
