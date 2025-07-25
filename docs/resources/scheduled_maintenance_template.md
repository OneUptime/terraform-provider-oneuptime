---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_scheduled_maintenance_template Resource - oneuptime"
subcategory: ""
description: |-
  Scheduled maintenance template resource
---

# oneuptime_scheduled_maintenance_template (Resource)

Scheduled maintenance template resource

## Example Usage

```terraform
resource "oneuptime_scheduled_maintenance_template" "example" {
  template_name = "example-template_name"
  template_description = "example-template_description"
  title = "example-title"
  description = "Example resource"
}
```

## Schema

- `id` (String) Unique identifier for the resource. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Optional.
- `template_name` (String) Name. Required.
- `template_description` (String) Template Description. Required.
- `title` (String) Title. Required.
- `description` (String) Description. Optional.
- `monitors` (List) Monitors. Optional.
- `status_pages` (List) Status Pages. Optional.
- `labels` (List) Labels. Optional.
- `change_monitor_status_to_id` (String) A unique identifier for an object, represented as a UUID.. Optional.
- `first_event_scheduled_at` (Map) A date time object.. Optional.
- `first_event_starts_at` (Map) A date time object.. Optional.
- `first_event_ends_at` (Map) A date time object.. Optional.
- `recurring_interval` (Map) Recurring Interval. Optional.
- `is_recurring_event` (Bool) Is Recurring Event. Optional.
- `schedule_next_event_at` (Map) A date time object.. Optional.
- `should_status_page_subscribers_be_notified_on_event_created` (Bool) Should subscribers be notified when event is created?. Optional.
- `should_status_page_subscribers_be_notified_when_event_changed_to_ongoing` (Bool) Should subscribers be notified when event is changed to ongoing?. Optional.
- `should_status_page_subscribers_be_notified_when_event_changed_to_ended` (Bool) Should subscribers be notified when event is changed to ended?. Optional.
- `custom_fields` (Map) Custom Fields. Optional.
- `send_subscriber_notifications_on_before_the_event` (Map) Subscriber notifications before the event. Optional.
- `created_at` (Map) A date time object.. Computed.
- `updated_at` (Map) A date time object.. Computed.
- `deleted_at` (Map) A date time object.. Computed.
- `version` (Number) Version. Computed.
- `slug` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance Template], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance Template], Update: [No access - you don't have permission for this operation]. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_scheduled_maintenance_template.example <id>
```
