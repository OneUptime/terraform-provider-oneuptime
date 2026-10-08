---
page_title: "oneuptime_scheduled_maintenance_event Resource - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Manage scheduled maintenance event for your project
---

# oneuptime_scheduled_maintenance_event (Resource)

Manage scheduled maintenance event for your project

## Example Usage

```terraform
resource "oneuptime_scheduled_maintenance_event" "example" {
  title       = "Example short text"
  starts_at   = "2030-01-01T00:00:00Z"
  ends_at     = "2030-01-01T00:00:00Z"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `ends_at` (String) When does this event end?
- `starts_at` (String) When does this event start?
- `title` (String) Title of this scheduled event.

### Optional

- `ceph_clusters` (Set of String) List of Ceph clusters affected by this event. IDs of `oneuptime_ceph_cluster` resources.
- `change_monitor_status_to_id` (String) Relation to Monitor Status Object ID. The monitors attached to this event change to this status when the event starts, and back to operational when it ends. It can be changed until the event starts. The ID of a `oneuptime_monitor_status`.
- `current_scheduled_maintenance_state_id` (String) Scheduled Maintenance State ID. The state the event currently is in. The ID of a `oneuptime_scheduled_maintenance_state`.
- `custom_fields` (String) Custom Fields on this resource. A JSON value: write it with `jsonencode()`.
- `database_servers` (Set of String) List of databases affected by this event. IDs of `oneuptime_database` resources.
- `description` (String) Description of this scheduled event that will show up on Status Page. This is in markdown.
- `docker_hosts` (Set of String) List of Docker hosts affected by this event. IDs of `oneuptime_docker_host` resources.
- `docker_swarm_clusters` (Set of String) List of Docker Swarm clusters affected by this event. IDs of `oneuptime_docker_swarm_cluster` resources.
- `enable_reminders` (Boolean) Should reminder notifications be sent to owners while this scheduled maintenance event is still not complete? Reminders are sent based on the reminder rules configured for this project. Defaults to `true`.
- `hosts` (Set of String) List of hosts affected by this event. IDs of `oneuptime_host` resources.
- `iot_fleets` (Set of String) List of IoT fleets affected by this event. IDs of `oneuptime_iot_fleet` resources.
- `is_visible_on_status_page` (Boolean) Should this incident be visible on the status page? Defaults to `true`.
- `kubernetes_clusters` (Set of String) List of Kubernetes clusters affected by this event. IDs of `oneuptime_kubernetes_cluster` resources.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `monitors` (Set of String) List of monitors attached to this event. IDs of `oneuptime_monitor` resources.
- `network_sites` (Set of String) List of network sites affected by this event. Their descendants are covered too. IDs of `oneuptime_network_site` resources.
- `next_subscriber_notification_before_the_event_at` (String) When will the next notification to subscribers be sent out?
- `podman_hosts` (Set of String) List of Podman hosts affected by this event. IDs of `oneuptime_podman_host` resources.
- `proxmox_clusters` (Set of String) List of Proxmox clusters affected by this event. IDs of `oneuptime_proxmox_cluster` resources.
- `send_subscriber_notifications_on_before_the_event` (String) Should subscribers be notified before the event? A JSON value: write it with `jsonencode()`.
- `services` (Set of String) List of services affected by this event. IDs of `oneuptime_service` resources.
- `should_status_page_subscribers_be_notified_on_event_created` (Boolean) Should subscribers be notified about this event creation? Defaults to `true`.
- `should_status_page_subscribers_be_notified_when_event_changed_to_ended` (Boolean) Should subscribers be notified about this event event is changed to ended? Defaults to `true`.
- `should_status_page_subscribers_be_notified_when_event_changed_to_ongoing` (Boolean) Should subscribers be notified about this event event is changed to ongoing? Defaults to `true`.
- `status_pages` (Set of String) List of status pages to show this event on. IDs of `oneuptime_status_page` resources.
- `storage_arrays` (Set of String) List of storage arrays affected by this event. IDs of `oneuptime_storage_array` resources.
- `vmware_v_centers` (Set of String) List of vCenters affected by this event. IDs of `oneuptime_vcenter` resources.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified_of_resource_creation` (Boolean) Are owners notified of when this resource is created?
- `next_reminder_notification_at` (String) When will the next reminder notification be sent to owners of this scheduled maintenance event? This is set automatically based on the reminder rules configured for this project.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `reminder_notification_sent_count` (Number) How many reminder notifications have been sent to owners of this scheduled maintenance event so far.
- `scheduled_maintenance_number` (Number) Scheduled Maintenance Number.
- `scheduled_maintenance_number_with_prefix` (String) Scheduled maintenance number with prefix (e.g., 'SM-42' or '#42').
- `slug` (String) Friendly globally unique name for your object.
- `subscriber_notification_status_message` (String) Status message for subscriber notifications when event is scheduled - includes success messages, failure reasons, or skip reasons.
- `subscriber_notification_status_on_event_scheduled` (String) Status of notification sent to subscribers when event was scheduled.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing scheduled maintenance event by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_scheduled_maintenance_event.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_scheduled_maintenance_event.example <id>
```
