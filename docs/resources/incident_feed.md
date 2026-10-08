---
page_title: "oneuptime_incident_feed Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Log of the entire incident state change. This is a log of all the incident state changes, public notes, more etc.
---

# oneuptime_incident_feed (Resource)

Log of the entire incident state change. This is a log of all the incident state changes, public notes, more etc.

## Example Usage

```terraform
resource "oneuptime_incident_feed" "example" {
  incident_id              = oneuptime_incident.example.id
  feed_info_in_markdown    = "# Heading\n\nThis is **markdown** content"
  incident_feed_event_type = "Example short text"
  display_color            = "#ff0000"
}
```

## Schema

### Required

- `display_color` (String) Display color for the incident log.
- `feed_info_in_markdown` (String) Log of the entire incident state change in Markdown.
- `incident_feed_event_type` (String) Incident Feed Event.
- `incident_id` (String) Relation to Incident ID in which this resource belongs. The ID of a `oneuptime_incident`.

### Optional

- `more_information_in_markdown` (String) More information in Markdown.
- `posted_at` (String) Date and time when the feed was posted.
- `user_id` (String) User who this feed belongs to (if this feed belongs to a User). The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `ai_run_id` (String) AI investigation run that produced this feed item when the event is a root cause.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incident feed by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incident_feed.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incident_feed.example <id>
```
