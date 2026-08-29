---
page_title: "oneuptime_scheduled_maintenance_event Data Source - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Manage scheduled maintenance event for your project
---

# oneuptime_scheduled_maintenance_event (Data Source)

Manage scheduled maintenance event for your project Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_scheduled_maintenance_event" "by_name" {
  name = "example-scheduled_maintenance_event"
}

data "oneuptime_scheduled_maintenance_event" "by_id" {
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
- `title` (String) Title of this scheduled event... Computed.
- `description` (String) Description of this scheduled event that will show up on Status Page. This is in markdown... Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `monitors` (Set) List of monitors attached to this event.. Computed.
- `hosts` (Set) List of hosts affected by this event... Computed.
- `kubernetes_clusters` (Set) List of Kubernetes clusters affected by this event... Computed.
- `docker_hosts` (Set) List of Docker hosts affected by this event... Computed.
- `podman_hosts` (Set) List of Podman hosts affected by this event... Computed.
- `proxmox_clusters` (Set) List of Proxmox clusters affected by this event... Computed.
- `iot_fleets` (Set) List of IoT fleets affected by this event... Computed.
- `network_sites` (Set) List of network sites affected by this event. Their descendants are covered too... Computed.
- `docker_swarm_clusters` (Set) List of Docker Swarm clusters affected by this event... Computed.
- `ceph_clusters` (Set) List of Ceph clusters affected by this event... Computed.
- `services` (Set) List of services affected by this event... Computed.
- `status_pages` (Set) List of status pages to show this event on.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `current_scheduled_maintenance_state_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `change_monitor_status_to_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `starts_at` (String) A date time object.. Computed.
- `ends_at` (String) A date time object.. Computed.
- `subscriber_notification_status_on_event_scheduled` (String) Status of notification sent to subscribers when event was scheduled.. Computed.
- `subscriber_notification_status_message` (String) Status message for subscriber notifications when event is scheduled - includes success messages, failure reasons, or skip reasons.. Computed.
- `should_status_page_subscribers_be_notified_on_event_created` (Bool) Should subscribers be notified about this event creation?.. Computed.
- `should_status_page_subscribers_be_notified_when_event_changed_to_ongoing` (Bool) Should subscribers be notified about this event event is changed to ongoing?.. Computed.
- `should_status_page_subscribers_be_notified_when_event_changed_to_ended` (Bool) Should subscribers be notified about this event event is changed to ended?.. Computed.
- `custom_fields` (String) Custom Fields on this resource... Computed.
- `is_owner_notified_of_resource_creation` (Bool) Are owners notified of when this resource is created?.. Computed.
- `send_subscriber_notifications_on_before_the_event` (String) Should subscribers be notified before the event?.. Computed.
- `next_subscriber_notification_before_the_event_at` (String) A date time object.. Computed.
- `scheduled_maintenance_number` (Number) Scheduled Maintenance Number.. Computed.
- `scheduled_maintenance_number_with_prefix` (String) Scheduled maintenance number with prefix (e.g., 'SM-42' or '#42').. Computed.
- `is_visible_on_status_page` (Bool) Should this incident be visible on the status page?.. Computed.
- `enable_reminders` (Bool) Should reminder notifications be sent to owners while this scheduled maintenance event is still not complete? Reminders are sent based on the reminder rules configured for this project... Computed.
- `next_reminder_notification_at` (String) A date time object.. Computed.
- `reminder_notification_sent_count` (Number) How many reminder notifications have been sent to owners of this scheduled maintenance event so far... Computed.
