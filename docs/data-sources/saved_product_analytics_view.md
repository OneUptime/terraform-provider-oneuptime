---
page_title: "oneuptime_saved_product_analytics_view Data Source - oneuptime"
subcategory: "Other"
description: |-
  Funnels and retention tables saved from a RUM application's product analytics, by name, to open again or show on a dashboard.
---

# oneuptime_saved_product_analytics_view (Data Source)

Funnels and retention tables saved from a RUM application's product analytics, by name, to open again or show on a dashboard.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one saved product analytics view may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_saved_product_analytics_view" "example" {
  name = "Example saved product analytics view"
}

# Or by id:
data "oneuptime_saved_product_analytics_view" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the user who saved this view. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `name` (String) What the team calls this funnel or retention table.
- `rum_application_id` (String) ID of the RUM application whose product analytics this saved view reads. The ID of a `oneuptime_rum_application`.
- `view_type` (String) What this saved view is: a funnel (funnel) or a retention table (retention).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `definition` (String) What it reads. A funnel: its steps (each an event's kind and name), the time to finish it (windowMs) and its segment (filters). A retention table: its startEvent, returnEvent, period (day, week or month) and segment (filters). The time range is not part of it: a page keeps its own, and a dashboard widget reads the dashboard's. A JSON value: read it with `jsondecode()`.
- `project_id` (String) ID of the project this saved view belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
