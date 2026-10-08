---
page_title: "oneuptime_incident Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Manage incidents for your project
---

# oneuptime_incident (Resource)

Manage incidents for your project

## Example Usage

```terraform
resource "oneuptime_incident" "example" {
  title                = "This is an example of longer text content that might be stored in this field."
  incident_severity_id = oneuptime_incident_severity.example.id
  description          = "Managed by Terraform"
}
```

## Schema

### Required

- `incident_severity_id` (String) Incident Severity ID. The ID of a `oneuptime_incident_severity`.
- `title` (String) Title of this incident.

### Optional

- `ceph_clusters` (Set of String) List of Ceph clusters affected by this incident. IDs of `oneuptime_ceph_cluster` resources.
- `change_monitor_status_to_id` (String) Relation to Monitor Status Object ID. All monitors connected to this incident will be changed to this status when the incident is created. The ID of a `oneuptime_monitor_status`.
- `current_incident_state_id` (String) Current Incident State ID. The ID of a `oneuptime_incident_state`.
- `custom_fields` (String) The incident's custom field values, keyed by each incident custom field's name. When a user or an API key creates or updates an incident, each value it sets or changes must fit its field - a number for a Number field, true or false for a Boolean, one of the options for a Dropdown, and so on - or the request is refused. Values left as they were, keys that are not the name of a field and empty values are not checked. Required on Create is not enforced here: it applies to the dashboard's Declare Incident form only. A JSON value: write it with `jsonencode()`.
- `database_servers` (Set of String) List of databases affected by this incident. IDs of `oneuptime_database` resources.
- `declared_at` (String) Date and time when this incident was declared.
- `description` (String) Short description of this incident. This is in markdown and will be visible on the status page.
- `docker_hosts` (Set of String) List of Docker hosts affected by this incident. IDs of `oneuptime_docker_host` resources.
- `docker_resources` (Set of String) List of Docker resources (containers, images, networks, volumes) affected by this incident.
- `docker_swarm_clusters` (Set of String) List of Docker Swarm clusters affected by this incident. IDs of `oneuptime_docker_swarm_cluster` resources.
- `enable_reminders` (Boolean) Should reminder notifications be sent to owners while this incident is still open? Reminders are sent based on the reminder rules configured for this project. Defaults to `true`.
- `hosts` (Set of String) List of hosts affected by this incident. IDs of `oneuptime_host` resources.
- `impact_started_at` (String) When customer impact actually began. Left blank until someone records it - never inferred, because a guessed value is worse than no value.
- `iot_fleets` (Set of String) List of IoT fleets affected by this incident. IDs of `oneuptime_iot_fleet` resources.
- `is_private` (Boolean) If true, this incident is only visible to its owners (users in 'owner users' and members of 'owner teams'), project admins, and project owners. Private incidents are hidden from status pages. Defaults to `false`.
- `is_visible_on_status_page` (Boolean) Should this incident be visible on the status page? Defaults to `true`.
- `kubernetes_clusters` (Set of String) List of Kubernetes clusters affected by this incident. IDs of `oneuptime_kubernetes_cluster` resources.
- `kubernetes_containers` (Set of String) List of Kubernetes containers affected by this incident.
- `kubernetes_resources` (Set of String) List of Kubernetes resources (pods, deployments, nodes, etc.) affected by this incident.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `monitors` (Set of String) List of monitors affected by this incident. IDs of `oneuptime_monitor` resources.
- `notify_subscribers_on_postmortem_published` (Boolean) Should subscribers be notified when the postmortem is published? Defaults to `true`.
- `on_call_duty_policies` (Set of String) List of on-call duty policies affected by this incident. IDs of `oneuptime_on_call_policy` resources.
- `podman_hosts` (Set of String) List of Podman hosts affected by this incident. IDs of `oneuptime_podman_host` resources.
- `podman_resources` (Set of String) List of Podman resources (containers, images, networks, volumes) affected by this incident.
- `postmortem_attachments` (Set of String) Files that accompany the postmortem note and can be shared publicly when enabled. IDs of `oneuptime_file` resources.
- `postmortem_note` (String) Document the postmortem summary for this incident.
- `postmortem_posted_at` (String) Timestamp that will be shown alongside the published postmortem on the status page.
- `proxmox_clusters` (Set of String) List of Proxmox clusters affected by this incident. IDs of `oneuptime_proxmox_cluster` resources.
- `remediation_notes` (String) Notes on how to remediate this incident. This is in markdown.
- `root_cause` (String) What is the root cause of this incident?
- `service_level_objectives` (Set of String) List of Service Level Objectives (SLOs) affected by this incident. IDs of `oneuptime_service_level_objective` resources.
- `services` (Set of String) List of services affected by this incident. IDs of `oneuptime_service` resources.
- `should_status_page_subscribers_be_notified_on_incident_created` (Boolean) Should subscribers be notified about this incident? Defaults to `true`.
- `show_postmortem_on_status_page` (Boolean) Should the postmortem note and attachments be visible on the status page once published? Defaults to `false`.
- `status_pages` (Set of String) Limit this incident to these status pages. When set, the incident is shown on, and notifies the subscribers of, only these pages among the status pages that list its monitors. Leave empty to reach every status page that lists its monitors. IDs of `oneuptime_status_page` resources.
- `storage_arrays` (Set of String) List of storage arrays affected by this incident. IDs of `oneuptime_storage_array` resources.
- `telemetry_query` (String) Telemetry query for this incident. A JSON value: write it with `jsonencode()`.
- `vmware_v_centers` (Set of String) List of vCenters affected by this incident. IDs of `oneuptime_vcenter` resources.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_probe_id` (String) If this incident was created by a Probe, this is the ID of the probe that created it. The ID of a `oneuptime_probe`.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `created_criteria_id` (String) If this incident was created by a Probe, this is the ID of the criteria that created it.
- `created_incident_template_id` (String) If this incident was created by a Probe, this is the ID of the incident template that was used for creation.
- `created_state_log` (String) A JSON value: write it with `jsonencode()`.
- `holds_monitors` (Boolean) Whether this incident is holding its monitors - keeping them in its monitor status, or their monitoring paused - so that resolving it gives them back: their monitoring resumes and their status returns to operational. True from when the incident is declared open, or from when an edit while it is open puts its monitors in its monitor status. False for an incident declared already resolved, which never held them, and once a resolve has given them back. Empty for incidents from before it was recorded, which give their monitors back when they are resolved. Set by OneUptime; it cannot be written.
- `id` (String) Unique identifier for the resource.
- `incident_episode_id` (String) ID of the latest Incident Episode this incident is a member of. Read-only: set by OneUptime when the incident is added to or removed from an episode's members (Incident Episode Member). The ID of a `oneuptime_incident_episode`.
- `incident_number` (Number) Incident Number.
- `incident_number_with_prefix` (String) Incident number with prefix (e.g., 'INC-42' or '#42').
- `is_created_automatically` (Boolean) Is this incident created by OneUptime Probe or Workers automatically (and not created manually by a user)?
- `is_owner_notified_of_resource_creation` (Boolean) Are owners notified of when this resource is created?
- `is_scoped_to_status_pages` (Boolean) Whether this incident is limited to the status pages in Status Pages. Derived from Status Pages; any value sent for it is ignored.
- `monitor_summary` (String) The monitor summary captured at the moment this incident was created - the same card the monitor page shows, frozen so it survives the monitor log being aged out. A JSON value: write it with `jsonencode()`.
- `next_reminder_notification_at` (String) When will the next reminder notification be sent to owners of this incident? This is set automatically based on the reminder rules configured for this project.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `reminder_notification_sent_count` (Number) How many reminder notifications have been sent to owners of this incident so far.
- `series_fingerprint` (String) For metric monitors with per-series alerting (e.g. grouped by host.name), this is a stable hash of the series label values so one incident is created per affected series.
- `series_labels` (String) Attribute key/value pairs that identify the affected series (e.g. {host.name: prod-db-01}) when this incident was created from a per-series metric breach. A JSON value: write it with `jsonencode()`.
- `slug` (String) Friendly globally unique name for your object.
- `status_pages_notified_on_creation` (String) IDs of the status pages whose subscribers were sent the notification that this incident was created. A JSON value: write it with `jsonencode()`.
- `subscriber_notification_status_message` (String) Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons.
- `subscriber_notification_status_message_on_postmortem_published` (String) Status message for subscriber notifications on postmortem published - includes success messages, failure reasons, or skip reasons.
- `subscriber_notification_status_on_incident_created` (String) Status of notification sent to subscribers about this incident.
- `subscriber_notification_status_on_postmortem_published` (String) Status of notification sent to subscribers about this incident postmortem.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incident by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incident.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incident.example <id>
```
