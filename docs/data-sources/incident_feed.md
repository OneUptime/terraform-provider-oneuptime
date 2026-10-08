---
page_title: "oneuptime_incident_feed Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Log of the entire incident state change. This is a log of all the incident state changes, public notes, more etc.
---

# oneuptime_incident_feed (Data Source)

Log of the entire incident state change. This is a log of all the incident state changes, public notes, more etc.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident feed may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_feed" "example" {
  incident_id = oneuptime_incident.example.id
}

# Or by id:
data "oneuptime_incident_feed" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `ai_run_id` (String) AI investigation run that produced this feed item when the event is a root cause.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `feed_info_in_markdown` (String) Log of the entire incident state change in Markdown.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_feed_event_type` (String) Incident Feed Event.
- `incident_id` (String) Relation to Incident ID in which this resource belongs. The ID of a `oneuptime_incident`.
- `more_information_in_markdown` (String) More information in Markdown.
- `user_id` (String) User who this feed belongs to (if this feed belongs to a User). The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `display_color` (String) Display color for the incident log.
- `posted_at` (String) Date and time when the feed was posted.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
