---
page_title: "oneuptime_incident_state_timeline Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Change state of the incidents (Created to Acknowledged for example)
---

# oneuptime_incident_state_timeline (Data Source)

Change state of the incidents (Created to Acknowledged for example)

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident state timeline may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_state_timeline" "example" {
  incident_id = oneuptime_incident.example.id
}

# Or by id:
data "oneuptime_incident_state_timeline" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) Relation to Incident ID in which this resource belongs. The ID of a `oneuptime_incident`.
- `incident_state_id` (String) Incident State ID Relation. Which incident state does this incident change to? The ID of a `oneuptime_incident_state`.
- `is_owner_notified` (Boolean) Are owners notified of state change?
- `root_cause` (String) What is the root cause of this status change?
- `should_status_page_subscribers_be_notified` (Boolean) Should subscribers be notified about this state change?
- `subscriber_notification_status` (String) Status of notification sent to subscribers about this incident state change.
- `subscriber_notification_status_message` (String) Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `ends_at` (String) When did this status change end?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `starts_at` (String) When did this status change? Correct this when the recorded time is wrong - every measurement derived from this timeline is recomputed from the corrected value.
- `state_change_log` (String) A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
