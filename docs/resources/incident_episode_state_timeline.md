---
page_title: "oneuptime_incident_episode_state_timeline Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Change state of the incident episodes (Created to Acknowledged for example)
---

# oneuptime_incident_episode_state_timeline (Resource)

Change state of the incident episodes (Created to Acknowledged for example)

## Example Usage

```terraform
resource "oneuptime_incident_episode_state_timeline" "example" {
  incident_episode_id = oneuptime_incident_episode.example.id
  incident_state_id   = oneuptime_incident_state.example.id
}
```

## Schema

### Required

- `incident_episode_id` (String) Relation to Incident Episode ID in which this resource belongs. The ID of a `oneuptime_incident_episode`.
- `incident_state_id` (String) Incident State ID Relation. Which incident state does this episode change to? The ID of a `oneuptime_incident_state`.

### Optional

- `ends_at` (String) When did this status change end?
- `root_cause` (String) What is the root cause of this status change?
- `should_status_page_subscribers_be_notified` (Boolean) Should status page subscribers be notified about this state change? Defaults to `true`.
- `starts_at` (String) When did this status change?

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of state change?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `state_change_log` (String) A JSON value: write it with `jsonencode()`.
- `subscriber_notification_status` (String) Status of notification sent to subscribers about this state change.
- `subscriber_notification_status_message` (String) Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incident episode state timeline by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incident_episode_state_timeline.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incident_episode_state_timeline.example <id>
```
