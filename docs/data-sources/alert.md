---
page_title: "oneuptime_alert Data Source - oneuptime"
subcategory: "Alerts"
description: |-
  Manage alerts for your project
---

# oneuptime_alert (Data Source)

Manage alerts for your project

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one alert may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_alert" "example" {
  title = "example-title"
}

# Or by id:
data "oneuptime_alert" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_episode_id` (String) The ID of the latest episode this alert is a member of, if any. Read-only: set by OneUptime when the alert is added to or removed from an episode's members (Alert Episode Member). The ID of a `oneuptime_alert_episode`.
- `alert_number` (Number) Alert Number.
- `alert_number_with_prefix` (String) Alert number with prefix (e.g., 'ALT-42' or '#42').
- `alert_severity_id` (String) Alert Severity ID. The ID of a `oneuptime_alert_severity`.
- `created_by_probe_id` (String) If this alert was created by a Probe, this is the ID of the probe that created it. The ID of a `oneuptime_probe`.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `created_criteria_id` (String) If this alert was created by a Probe, this is the ID of the criteria that created it.
- `current_alert_state_id` (String) Current Alert State ID. The ID of a `oneuptime_alert_state`.
- `description` (String) Short description of this alert. This will be visible on the status page. This is in markdown.
- `enable_reminders` (Boolean) Should reminder notifications be sent to owners while this alert is still open? Reminders are sent based on the reminder rules configured for this project.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_created_automatically` (Boolean) Is this alert created by OneUptime Probe or Workers automatically (and not created manually by a user)?
- `is_owner_notified_of_alert_creation` (Boolean) Are owners notified of when this alert is created?
- `is_private` (Boolean) If true, this alert is only visible to its owners (users in 'owner users' and members of 'owner teams'), project admins, and project owners.
- `monitor_id` (String) ID of the monitor this alert belongs to. The ID of a `oneuptime_monitor`.
- `monitor_status_when_this_alert_was_created_id` (String) Monitor Status ID when this alert was created. The ID of a `oneuptime_monitor_status`.
- `remediation_notes` (String) Notes on how to remediate this alert. This is in markdown.
- `reminder_notification_sent_count` (Number) How many reminder notifications have been sent to owners of this alert so far.
- `root_cause` (String) What is the root cause of this alert?
- `series_fingerprint` (String) For metric monitors with per-series alerting (e.g. grouped by host.name), this is a stable hash of the series label values so one alert is created per affected series.
- `title` (String) Title of this alert.

### Read-Only

- `ceph_clusters` (Set of String) List of Ceph clusters affected by this alert. IDs of `oneuptime_ceph_cluster` resources.
- `created_at` (String) Date and Time when the object was created.
- `created_state_log` (String) A JSON value: write it with `jsonencode()`.
- `custom_fields` (String) Custom Fields on this resource. A JSON value: write it with `jsonencode()`.
- `database_servers` (Set of String) List of databases affected by this alert. IDs of `oneuptime_database` resources.
- `docker_hosts` (Set of String) List of Docker hosts affected by this alert. IDs of `oneuptime_docker_host` resources.
- `docker_resources` (Set of String) List of Docker resources (containers, images, networks, volumes) affected by this alert.
- `docker_swarm_clusters` (Set of String) List of Docker Swarm clusters affected by this alert. IDs of `oneuptime_docker_swarm_cluster` resources.
- `hosts` (Set of String) List of hosts affected by this alert. IDs of `oneuptime_host` resources.
- `impact_started_at` (String) When customer impact actually began. Left blank until someone records it - never inferred, because a guessed value is worse than no value.
- `iot_fleets` (Set of String) List of IoT fleets affected by this alert. IDs of `oneuptime_iot_fleet` resources.
- `kubernetes_clusters` (Set of String) List of Kubernetes clusters affected by this alert. IDs of `oneuptime_kubernetes_cluster` resources.
- `kubernetes_containers` (Set of String) List of Kubernetes containers affected by this alert.
- `kubernetes_resources` (Set of String) List of Kubernetes resources (pods, deployments, nodes, etc.) affected by this alert.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `monitor_summary` (String) The monitor summary captured at the moment this alert was created - the same card the monitor page shows, frozen so it survives the monitor log being aged out. A JSON value: write it with `jsonencode()`.
- `next_reminder_notification_at` (String) When will the next reminder notification be sent to owners of this alert? This is set automatically based on the reminder rules configured for this project.
- `on_call_duty_policies` (Set of String) List of on-call duty policies affected by this alert. IDs of `oneuptime_on_call_policy` resources.
- `podman_hosts` (Set of String) List of Podman hosts affected by this alert. IDs of `oneuptime_podman_host` resources.
- `podman_resources` (Set of String) List of Podman resources (containers, images, networks, volumes) affected by this alert.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `proxmox_clusters` (Set of String) List of Proxmox clusters affected by this alert. IDs of `oneuptime_proxmox_cluster` resources.
- `series_labels` (String) Attribute key/value pairs that identify the affected series (e.g. {host.name: prod-db-01}) when this alert was created from a per-series metric breach. A JSON value: write it with `jsonencode()`.
- `service_level_objectives` (Set of String) List of Service Level Objectives (SLOs) affected by this alert. IDs of `oneuptime_service_level_objective` resources.
- `services` (Set of String) List of services affected by this alert. IDs of `oneuptime_service` resources.
- `storage_arrays` (Set of String) List of storage arrays affected by this alert. IDs of `oneuptime_storage_array` resources.
- `telemetry_query` (String) Telemetry query for this alert. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
- `vmware_v_centers` (Set of String) List of vCenters affected by this alert. IDs of `oneuptime_vcenter` resources.
