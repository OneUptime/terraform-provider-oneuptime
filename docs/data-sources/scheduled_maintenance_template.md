---
page_title: "oneuptime_scheduled_maintenance_template Data Source - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Manage scheduled maintenance templates for your project
---

# oneuptime_scheduled_maintenance_template (Data Source)

Manage scheduled maintenance templates for your project Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_scheduled_maintenance_template" "by_name" {
  name = "example-scheduled_maintenance_template"
}

data "oneuptime_scheduled_maintenance_template" "by_id" {
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
- `template_name` (String) Name of the Scheduled Maintenance Template.. Computed.
- `template_description` (String) Description of the Scheduled Maintenance Template.. Computed.
- `title` (String) Title of this scheduled event... Computed.
- `description` (String) Description of this scheduled event that will show up on Status Page. This is a markdown field... Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `monitors` (Set) List of monitors attached to this event.. Computed.
- `hosts` (Set) List of hosts to pre-populate on scheduled maintenance events created from this template... Computed.
- `kubernetes_clusters` (Set) List of Kubernetes clusters to pre-populate on scheduled maintenance events created from this template... Computed.
- `docker_hosts` (Set) List of Docker hosts to pre-populate on scheduled maintenance events created from this template... Computed.
- `podman_hosts` (Set) List of Podman hosts to pre-populate on scheduled maintenance events created from this template... Computed.
- `services` (Set) List of services to pre-populate on scheduled maintenance events created from this template... Computed.
- `status_pages` (Set) List of status pages to show this event on.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `change_monitor_status_to_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `first_event_scheduled_at` (String) A date time object.. Computed.
- `first_event_starts_at` (String) A date time object.. Computed.
- `first_event_ends_at` (String) A date time object.. Computed.
- `recurring_interval` (String) How often should this event recur?.. Computed.
- `is_recurring_event` (Bool) Is this a recurring event?.. Computed.
- `schedule_next_event_at` (String) A date time object.. Computed.
- `should_status_page_subscribers_be_notified_on_event_created` (Bool) Should subscribers be notified about this event creation?.. Computed.
- `should_status_page_subscribers_be_notified_when_event_changed_to_ongoing` (Bool) Should subscribers be notified about this event event is changed to ongoing?.. Computed.
- `should_status_page_subscribers_be_notified_when_event_changed_to_ended` (Bool) Should subscribers be notified about this event event is changed to ended?.. Computed.
- `custom_fields` (String) Custom Fields on this resource... Computed.
- `send_subscriber_notifications_on_before_the_event` (String) Should subscribers be notified before the event?.. Computed.
