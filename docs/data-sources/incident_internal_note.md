---
page_title: "oneuptime_incident_internal_note Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Manage internal notes for your incident
---

# oneuptime_incident_internal_note (Data Source)

Manage internal notes for your incident

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident internal note may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_internal_note" "example" {
  incident_id = oneuptime_incident.example.id
}

# Or by id:
data "oneuptime_incident_internal_note" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) Relation to Incident ID in which this resource belongs. The ID of a `oneuptime_incident`.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `note` (String) Notes in markdown.
- `posted_from_slack_message_id` (String) Unique identifier for the Slack message this note was created from (channel_id:message_ts). Used to prevent duplicate notes when multiple users react to the same message.

### Read-Only

- `attachments` (Set of String) Files attached to this note. IDs of `oneuptime_file` resources.
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
