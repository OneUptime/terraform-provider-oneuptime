---
page_title: "oneuptime_metric_saved_view Resource - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Save and reuse metrics explorer views, including the current search, filters, time range, and page size.
---

# oneuptime_metric_saved_view (Resource)

Save and reuse metrics explorer views, including the current search, filters, time range, and page size.

## Example Usage

```terraform
resource "oneuptime_metric_saved_view" "example" {
  name = "Example metric saved view"
}
```

## Schema

### Required

- `name` (String) Friendly name for this saved metric view.

### Optional

- `is_default` (Boolean) Whether this saved metric view should be applied by default. Defaults to `false`.
- `query` (String) Serialized metrics explorer view state (search, filters, time range, page size) for this saved view. A JSON value: write it with `jsonencode()`.
- `view_type` (String) Which surface this saved view belongs to ('list' or 'explorer'). Null means 'list' — rows created before this column existed all came from the metric list page.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this saved metric view. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this saved metric view belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing metric saved view by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_metric_saved_view.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_metric_saved_view.example <id>
```
