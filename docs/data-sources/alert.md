---
page_title: "oneuptime_alert Data Source - oneuptime"
subcategory: "Alerts"
description: |-
  Manage alerts for your project
---

# oneuptime_alert (Data Source)

Manage alerts for your project Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_alert" "by_name" {
  name = "example-alert"
}

data "oneuptime_alert" "by_id" {
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
- `title` (String) Title of this alert.. Computed.
- `description` (String) Short description of this alert. This will be visible on the status page. This is in markdown... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `monitor_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policies` (Set) List of on-call duty policies affected by this alert... Computed.
- `hosts` (Set) List of hosts affected by this alert... Computed.
- `kubernetes_clusters` (Set) List of Kubernetes clusters affected by this alert... Computed.
- `kubernetes_resources` (Set) List of Kubernetes resources (pods, deployments, nodes, etc.) affected by this alert... Computed.
- `kubernetes_containers` (Set) List of Kubernetes containers affected by this alert... Computed.
- `docker_hosts` (Set) List of Docker hosts affected by this alert... Computed.
- `podman_hosts` (Set) List of Podman hosts affected by this alert... Computed.
- `proxmox_clusters` (Set) List of Proxmox clusters affected by this alert... Computed.
- `iot_fleets` (Set) List of IoT fleets affected by this alert... Computed.
- `docker_swarm_clusters` (Set) List of Docker Swarm clusters affected by this alert... Computed.
- `ceph_clusters` (Set) List of Ceph clusters affected by this alert... Computed.
- `docker_resources` (Set) List of Docker resources (containers, images, networks, volumes) affected by this alert... Computed.
- `podman_resources` (Set) List of Podman resources (containers, images, networks, volumes) affected by this alert... Computed.
- `services` (Set) List of services affected by this alert... Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `current_alert_state_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `alert_severity_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `monitor_status_when_this_alert_was_created_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `custom_fields` (String) Custom Fields on this resource... Computed.
- `is_owner_notified_of_alert_creation` (Bool) Are owners notified of when this alert is created?.. Computed.
- `root_cause` (String) What is the root cause of this alert?.. Computed.
- `created_state_log` (String) Alert created_state_log. Computed.
- `created_criteria_id` (String) If this alert was created by a Probe, this is the ID of the criteria that created it... Computed.
- `series_fingerprint` (String) For metric monitors with per-series alerting (e.g. grouped by host.name), this is a stable hash of the series label values so one alert is created per affected series... Computed.
- `series_labels` (String) Attribute key/value pairs that identify the affected series (e.g. {host.name: prod-db-01}) when this alert was created from a per-series metric breach... Computed.
- `monitor_summary` (String) The monitor summary captured at the moment this alert was created - the same card the monitor page shows, frozen so it survives the monitor log being aged out... Computed.
- `created_by_probe_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_created_automatically` (Bool) Is this alert created by OneUptime Probe or Workers automatically (and not created manually by a user)?.. Computed.
- `remediation_notes` (String) Notes on how to remediate this alert. This is in markdown... Computed.
- `telemetry_query` (String) Telemetry query for this alert.. Computed.
- `alert_number` (Number) Alert Number.. Computed.
- `alert_number_with_prefix` (String) Alert number with prefix (e.g., 'ALT-42' or '#42').. Computed.
- `alert_episode_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_private` (Bool) If true, this alert is only visible to its owners (users in 'owner users' and members of 'owner teams'), project admins, and project owners... Computed.
- `enable_reminders` (Bool) Should reminder notifications be sent to owners while this alert is still open? Reminders are sent based on the reminder rules configured for this project... Computed.
- `next_reminder_notification_at` (String) A date time object.. Computed.
- `reminder_notification_sent_count` (Number) How many reminder notifications have been sent to owners of this alert so far... Computed.
