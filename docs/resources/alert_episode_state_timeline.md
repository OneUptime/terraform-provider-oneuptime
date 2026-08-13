---
page_title: "oneuptime_alert_episode_state_timeline Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Change state of the alert episodes (Created to Acknowledged for example)
---

# oneuptime_alert_episode_state_timeline (Resource)

Change state of the alert episodes (Created to Acknowledged for example)

## Example Usage

```terraform
resource "oneuptime_alert_episode_state_timeline" "example" {
  alert_episode_id = "123e4567-e89b-12d3-a456-426614174000"
  alert_state_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `alert_episode_id` (String) A unique identifier for an object, represented as a UUID..
- `alert_state_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `root_cause` (String) What is the root cause of this status change?..
- `ends_at` (String) A date time object..
- `starts_at` (String) A date time object..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `is_owner_notified` (Bool) Are owners notified of state change?..
- `state_change_log` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Alert Admin, Alert Member, Alert Viewer, Read Alert Episode State Timeline], Update: [No access - you don't have permission for this operation].

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_alert_episode_state_timeline.example <id>
```
