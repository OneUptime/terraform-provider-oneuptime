---
page_title: "oneuptime_hunt_run Data Source - oneuptime"
subcategory: "Security"
description: |-
  Each scheduled run of a hunt: the events it searched, how many it found, or why it could not run. Kept for 90 days.
---

# oneuptime_hunt_run (Data Source)

Each scheduled run of a hunt: the events it searched, how many it found, or why it could not run. Kept for 90 days.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one hunt run may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_hunt_run" "example" {
  security_hunt_id = oneuptime_hunt.example.id
}

# Or by id:
data "oneuptime_hunt_run" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `error` (String) Why the run could not be made. Empty when it ran; its window is searched again by the next run.
- `event_count` (Number) How many events the run's search found in its window. Empty when the run could not be made.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `security_hunt_id` (String) ID of the hunt that ran. The ID of a `oneuptime_hunt`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this run's hunt belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
- `window_end_at` (String) The end of the events the run searched, two minutes before it ran.
- `window_start_at` (String) The start of the events the run searched: the moment after the run before it stopped.
