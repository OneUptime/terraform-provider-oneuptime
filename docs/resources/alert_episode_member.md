---
page_title: "oneuptime_alert_episode_member Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Link between alerts and episodes
---

# oneuptime_alert_episode_member (Resource)

Link between alerts and episodes

## Example Usage

```terraform
resource "oneuptime_alert_episode_member" "example" {
  alert_episode_id = "123e4567-e89b-12d3-a456-426614174000"
  alert_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `alert_episode_id` (String) A unique identifier for an object, represented as a UUID..
- `alert_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `added_at` (String) A date time object..
- `added_by` (String) How this alert was added to the episode (rule, manual, or api)..
- `added_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `matched_rule_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `is_owner_notified_of_alert_added` (Bool) Has the owner been notified that this alert was added to the episode?..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_alert_episode_member.example <id>
```
