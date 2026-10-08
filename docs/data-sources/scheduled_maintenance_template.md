---
page_title: "oneuptime_scheduled_maintenance_template Data Source - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Manage scheduled maintenance templates for your project
---

# oneuptime_scheduled_maintenance_template (Data Source)

Manage scheduled maintenance templates for your project

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one scheduled maintenance template may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_scheduled_maintenance_template" "example" {
  template_name = "example-template-name"
}

# Or by id:
data "oneuptime_scheduled_maintenance_template" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `change_monitor_status_to_id` (String) Relation to Monitor Status Object ID. All monitors connected to this incident will be changed to this status when the event is ongoing. The ID of a `oneuptime_monitor_status`.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this scheduled event that will show up on Status Page. This is a markdown field.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_recurring_event` (Boolean) Is this a recurring event?
- `should_status_page_subscribers_be_notified_on_event_created` (Boolean) Should subscribers be notified about this event creation?
- `should_status_page_subscribers_be_notified_when_event_changed_to_ended` (Boolean) Should subscribers be notified about this event event is changed to ended?
- `should_status_page_subscribers_be_notified_when_event_changed_to_ongoing` (Boolean) Should subscribers be notified about this event event is changed to ongoing?
- `slug` (String) Friendly globally unique name for your object.
- `template_description` (String) Description of the Scheduled Maintenance Template.
- `template_name` (String) Name of the Scheduled Maintenance Template.
- `title` (String) Title of this scheduled event.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `custom_fields` (String) Custom Fields on this resource. A JSON value: write it with `jsonencode()`.
- `docker_hosts` (Set of String) List of Docker hosts to pre-populate on scheduled maintenance events created from this template. IDs of `oneuptime_docker_host` resources.
- `first_event_ends_at` (String) When does the first event end?
- `first_event_scheduled_at` (String) When would you like to schedule the first event?
- `first_event_starts_at` (String) When does the first event start?
- `hosts` (Set of String) List of hosts to pre-populate on scheduled maintenance events created from this template. IDs of `oneuptime_host` resources.
- `kubernetes_clusters` (Set of String) List of Kubernetes clusters to pre-populate on scheduled maintenance events created from this template. IDs of `oneuptime_kubernetes_cluster` resources.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `monitors` (Set of String) List of monitors attached to this event. IDs of `oneuptime_monitor` resources.
- `podman_hosts` (Set of String) List of Podman hosts to pre-populate on scheduled maintenance events created from this template. IDs of `oneuptime_podman_host` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `recurring_interval` (String) How often should this event recur? A JSON value: write it with `jsonencode()`.
- `schedule_next_event_at` (String) When is the next event scheduled?
- `send_subscriber_notifications_on_before_the_event` (String) Should subscribers be notified before the event? A JSON value: write it with `jsonencode()`.
- `services` (Set of String) List of services to pre-populate on scheduled maintenance events created from this template. IDs of `oneuptime_service` resources.
- `status_pages` (Set of String) List of status pages to show this event on. IDs of `oneuptime_status_page` resources.
- `updated_at` (String) Date and Time when the object was updated.
