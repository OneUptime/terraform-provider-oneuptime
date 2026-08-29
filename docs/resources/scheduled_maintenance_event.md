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
  title = "Example short text"
  starts_at = "2030-01-01T00:00:00Z"
  ends_at = "2030-01-01T00:00:00Z"
  description = "# Heading

This is **markdown** content"
}
```

## Schema

### Required

- `title` (String) Title of this scheduled event...
- `starts_at` (String) A date time object..
- `ends_at` (String) A date time object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this scheduled event that will show up on Status Page. This is in markdown...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `monitors` (Set) List of monitors attached to this event..
- `hosts` (Set) List of hosts affected by this event...
- `kubernetes_clusters` (Set) List of Kubernetes clusters affected by this event...
- `docker_hosts` (Set) List of Docker hosts affected by this event...
- `podman_hosts` (Set) List of Podman hosts affected by this event...
- `proxmox_clusters` (Set) List of Proxmox clusters affected by this event...
- `iot_fleets` (Set) List of IoT fleets affected by this event...
- `network_sites` (Set) List of network sites affected by this event. Their descendants are covered too...
- `docker_swarm_clusters` (Set) List of Docker Swarm clusters affected by this event...
- `ceph_clusters` (Set) List of Ceph clusters affected by this event...
- `services` (Set) List of services affected by this event...
- `status_pages` (Set) List of status pages to show this event on..
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `current_scheduled_maintenance_state_id` (String) A unique identifier for an object, represented as a UUID..
- `change_monitor_status_to_id` (String) A unique identifier for an object, represented as a UUID..
- `should_status_page_subscribers_be_notified_on_event_created` (Bool) Should subscribers be notified about this event creation?..
- `should_status_page_subscribers_be_notified_when_event_changed_to_ongoing` (Bool) Should subscribers be notified about this event event is changed to ongoing?..
- `should_status_page_subscribers_be_notified_when_event_changed_to_ended` (Bool) Should subscribers be notified about this event event is changed to ended?..
- `custom_fields` (String) Custom Fields on this resource...
- `send_subscriber_notifications_on_before_the_event` (String) Should subscribers be notified before the event?..
- `next_subscriber_notification_before_the_event_at` (String) A date time object..
- `is_visible_on_status_page` (Bool) Should this incident be visible on the status page?..
- `enable_reminders` (Bool) Should reminder notifications be sent to owners while this scheduled maintenance event is still not complete? Reminders are sent based on the reminder rules configured for this project...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `subscriber_notification_status_on_event_scheduled` (String) Status of notification sent to subscribers when event was scheduled..
- `subscriber_notification_status_message` (String) Status message for subscriber notifications when event is scheduled - includes success messages, failure reasons, or skip reasons..
- `is_owner_notified_of_resource_creation` (Bool) Are owners notified of when this resource is created?..
- `scheduled_maintenance_number` (Number) Scheduled Maintenance Number..
- `scheduled_maintenance_number_with_prefix` (String) Scheduled maintenance number with prefix (e.g., 'SM-42' or '#42')..
- `next_reminder_notification_at` (String) A date time object..
- `reminder_notification_sent_count` (Number) How many reminder notifications have been sent to owners of this scheduled maintenance event so far...

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_scheduled_maintenance_event.example <id>
```
