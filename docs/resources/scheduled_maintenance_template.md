---
page_title: "oneuptime_scheduled_maintenance_template Resource - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Manage scheduled maintenance templates for your project
---

# oneuptime_scheduled_maintenance_template (Resource)

Manage scheduled maintenance templates for your project

## Example Usage

```terraform
resource "oneuptime_scheduled_maintenance_template" "example" {
  template_name = "Example short text"
  template_description = "This is an example of longer text content that might be stored in this field."
  title = "Example short text"
  description = "# Heading

This is **markdown** content"
}
```

## Schema

### Required

- `template_name` (String) Name of the Scheduled Maintenance Template..
- `template_description` (String) Description of the Scheduled Maintenance Template..
- `title` (String) Title of this scheduled event...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this scheduled event that will show up on Status Page. This is a markdown field...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `monitors` (Set) List of monitors attached to this event..
- `hosts` (Set) List of hosts to pre-populate on scheduled maintenance events created from this template...
- `kubernetes_clusters` (Set) List of Kubernetes clusters to pre-populate on scheduled maintenance events created from this template...
- `docker_hosts` (Set) List of Docker hosts to pre-populate on scheduled maintenance events created from this template...
- `podman_hosts` (Set) List of Podman hosts to pre-populate on scheduled maintenance events created from this template...
- `services` (Set) List of services to pre-populate on scheduled maintenance events created from this template...
- `status_pages` (Set) List of status pages to show this event on..
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `change_monitor_status_to_id` (String) A unique identifier for an object, represented as a UUID..
- `first_event_scheduled_at` (String) A date time object..
- `first_event_starts_at` (String) A date time object..
- `first_event_ends_at` (String) A date time object..
- `recurring_interval` (String) How often should this event recur?..
- `is_recurring_event` (Bool) Is this a recurring event?..
- `should_status_page_subscribers_be_notified_on_event_created` (Bool) Should subscribers be notified about this event creation?..
- `should_status_page_subscribers_be_notified_when_event_changed_to_ongoing` (Bool) Should subscribers be notified about this event event is changed to ongoing?..
- `should_status_page_subscribers_be_notified_when_event_changed_to_ended` (Bool) Should subscribers be notified about this event event is changed to ended?..
- `custom_fields` (String) Custom Fields on this resource...
- `send_subscriber_notifications_on_before_the_event` (String) Should subscribers be notified before the event?..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `schedule_next_event_at` (String) A date time object..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_scheduled_maintenance_template.example <id>
```
