---
page_title: "oneuptime_service_ai_note Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  What a service's team tells OneUptime AI about it. Every investigation of an incident or an alert on the service reads its notes.
---

# oneuptime_service_ai_note (Data Source)

What a service's team tells OneUptime AI about it. Every investigation of an incident or an alert on the service reads its notes.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one service ai note may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_service_ai_note" "example" {
  service_id = oneuptime_service.example.id
}

# Or by id:
data "oneuptime_service_ai_note" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the person who wrote the note. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `note` (String) What the team knows about the service that OneUptime AI should read when it investigates it, in plain text, at most 2000 characters.
- `service_id` (String) ID of the service the note is about. The ID of a `oneuptime_service`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this note belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
