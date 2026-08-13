---
page_title: "oneuptime_incident_state_timeline Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Change state of the incidents (Created to Acknowledged for example)
---

# oneuptime_incident_state_timeline (Data Source)

Change state of the incidents (Created to Acknowledged for example) Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incident_state_timeline" "by_name" {
  name = "example-incident_state_timeline"
}

data "oneuptime_incident_state_timeline" "by_id" {
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
- `incident_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `incident_state_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `subscriber_notification_status` (String) Status of notification sent to subscribers about this incident state change.. Computed.
- `subscriber_notification_status_message` (String) Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons.. Computed.
- `should_status_page_subscribers_be_notified` (Bool) Should subscribers be notified about this state change?.. Computed.
- `is_owner_notified` (Bool) Are owners notified of state change?.. Computed.
- `state_change_log` (String) Incident state timeline state_change_log. Computed.
- `root_cause` (String) What is the root cause of this status change?.. Computed.
- `ends_at` (String) A date time object.. Computed.
- `starts_at` (String) A date time object.. Computed.
