---
page_title: "oneuptime_alert_episode_state_timeline Data Source - oneuptime"
subcategory: "Alerts"
description: |-
  Change state of the alert episodes (Created to Acknowledged for example)
---

# oneuptime_alert_episode_state_timeline (Data Source)

Change state of the alert episodes (Created to Acknowledged for example) Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_alert_episode_state_timeline" "by_name" {
  name = "example-alert_episode_state_timeline"
}

data "oneuptime_alert_episode_state_timeline" "by_id" {
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
- `alert_episode_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `alert_state_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_owner_notified` (Bool) Are owners notified of state change?.. Computed.
- `state_change_log` (String) Alert episode state timeline state_change_log. Computed.
- `root_cause` (String) What is the root cause of this status change?.. Computed.
- `ends_at` (String) A date time object.. Computed.
- `starts_at` (String) A date time object.. Computed.
