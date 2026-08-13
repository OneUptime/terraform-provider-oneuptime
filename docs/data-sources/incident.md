---
page_title: "oneuptime_incident Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Manage incidents for your project
---

# oneuptime_incident (Data Source)

Manage incidents for your project Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incident" "by_name" {
  name = "example-incident"
}

data "oneuptime_incident" "by_id" {
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
- `title` (String) Title of this incident.. Computed.
- `description` (String) Short description of this incident. This is in markdown and will be visible on the status page... Computed.
- `declared_at` (String) A date time object.. Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `monitors` (Set) List of monitors affected by this incident.. Computed.
- `hosts` (Set) List of hosts affected by this incident... Computed.
- `kubernetes_clusters` (Set) List of Kubernetes clusters affected by this incident... Computed.
- `kubernetes_resources` (Set) List of Kubernetes resources (pods, deployments, nodes, etc.) affected by this incident... Computed.
- `kubernetes_containers` (Set) List of Kubernetes containers affected by this incident... Computed.
- `docker_hosts` (Set) List of Docker hosts affected by this incident... Computed.
- `podman_hosts` (Set) List of Podman hosts affected by this incident... Computed.
- `proxmox_clusters` (Set) List of Proxmox clusters affected by this incident... Computed.
- `iot_fleets` (Set) List of IoT fleets affected by this incident... Computed.
- `docker_swarm_clusters` (Set) List of Docker Swarm clusters affected by this incident... Computed.
- `ceph_clusters` (Set) List of Ceph clusters affected by this incident... Computed.
- `docker_resources` (Set) List of Docker resources (containers, images, networks, volumes) affected by this incident... Computed.
- `podman_resources` (Set) List of Podman resources (containers, images, networks, volumes) affected by this incident... Computed.
- `services` (Set) List of services affected by this incident... Computed.
- `on_call_duty_policies` (Set) List of on-call duty policies affected by this incident... Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `current_incident_state_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `incident_severity_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `change_monitor_status_to_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `subscriber_notification_status_on_incident_created` (String) Status of notification sent to subscribers about this incident.. Computed.
- `subscriber_notification_status_message` (String) Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons.. Computed.
- `subscriber_notification_status_on_postmortem_published` (String) Status of notification sent to subscribers about this incident postmortem.. Computed.
- `subscriber_notification_status_message_on_postmortem_published` (String) Status message for subscriber notifications on postmortem published - includes success messages, failure reasons, or skip reasons.. Computed.
- `should_status_page_subscribers_be_notified_on_incident_created` (Bool) Should subscribers be notified about this incident?.. Computed.
- `custom_fields` (String) Custom Fields on this resource... Computed.
- `is_owner_notified_of_resource_creation` (Bool) Are owners notified of when this resource is created?.. Computed.
- `root_cause` (String) What is the root cause of this incident?.. Computed.
- `postmortem_note` (String) Document the postmortem summary for this incident... Computed.
- `show_postmortem_on_status_page` (Bool) Should the postmortem note and attachments be visible on the status page once published?.. Computed.
- `notify_subscribers_on_postmortem_published` (Bool) Should subscribers be notified when the postmortem is published?.. Computed.
- `postmortem_posted_at` (String) A date time object.. Computed.
- `postmortem_attachments` (Set) Files that accompany the postmortem note and can be shared publicly when enabled... Computed.
- `created_state_log` (String) Incident created_state_log. Computed.
- `created_criteria_id` (String) If this incident was created by a Probe, this is the ID of the criteria that created it... Computed.
- `created_incident_template_id` (String) If this incident was created by a Probe, this is the ID of the incident template that was used for creation... Computed.
- `series_fingerprint` (String) For metric monitors with per-series alerting (e.g. grouped by host.name), this is a stable hash of the series label values so one incident is created per affected series... Computed.
- `series_labels` (String) Attribute key/value pairs that identify the affected series (e.g. {host.name: prod-db-01}) when this incident was created from a per-series metric breach... Computed.
- `monitor_summary` (String) The monitor summary captured at the moment this incident was created - the same card the monitor page shows, frozen so it survives the monitor log being aged out... Computed.
- `created_by_probe_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_created_automatically` (Bool) Is this incident created by OneUptime Probe or Workers automatically (and not created manually by a user)?.. Computed.
- `remediation_notes` (String) Notes on how to remediate this incident. This is in markdown... Computed.
- `telemetry_query` (String) Telemetry query for this incident.. Computed.
- `incident_number` (Number) Incident Number.. Computed.
- `incident_number_with_prefix` (String) Incident number with prefix (e.g., 'INC-42' or '#42').. Computed.
- `is_visible_on_status_page` (Bool) Should this incident be visible on the status page?.. Computed.
- `is_private` (Bool) If true, this incident is only visible to its owners (users in 'owner users' and members of 'owner teams'), project admins, and project owners. Private incidents are hidden from status pages... Computed.
- `enable_reminders` (Bool) Should reminder notifications be sent to owners while this incident is still open? Reminders are sent based on the reminder rules configured for this project... Computed.
- `next_reminder_notification_at` (String) A date time object.. Computed.
- `reminder_notification_sent_count` (Number) How many reminder notifications have been sent to owners of this incident so far... Computed.
- `incident_episode_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
