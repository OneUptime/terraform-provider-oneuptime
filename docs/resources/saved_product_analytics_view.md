---
page_title: "oneuptime_saved_product_analytics_view Resource - oneuptime"
subcategory: "Other"
description: |-
  Funnels and retention tables saved from a RUM application's product analytics, by name, to open again or show on a dashboard.
---

# oneuptime_saved_product_analytics_view (Resource)

Funnels and retention tables saved from a RUM application's product analytics, by name, to open again or show on a dashboard.

## Example Usage

```terraform
resource "oneuptime_saved_product_analytics_view" "example" {
  rum_application_id = oneuptime_rum_application.example.id
  name               = "Example saved product analytics view"
  view_type          = "funnel"
}
```

## Schema

### Required

- `name` (String) What the team calls this funnel or retention table.
- `rum_application_id` (String) ID of the RUM application whose product analytics this saved view reads. The ID of a `oneuptime_rum_application`.
- `view_type` (String) What this saved view is: a funnel (funnel) or a retention table (retention). Allowed values: `funnel`, `retention`.

### Optional

- `definition` (String) What it reads. A funnel: its steps (each an event's kind and name), the time to finish it (windowMs) and its segment (filters). A retention table: its startEvent, returnEvent, period (day, week or month) and segment (filters). The time range is not part of it: a page keeps its own, and a dashboard widget reads the dashboard's. A JSON value: write it with `jsonencode()`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who saved this view. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this saved view belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing saved product analytics view by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_saved_product_analytics_view.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_saved_product_analytics_view.example <id>
```
