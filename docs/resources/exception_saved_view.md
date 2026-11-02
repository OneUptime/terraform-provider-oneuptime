---
page_title: "oneuptime_exception_saved_view Resource - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Save and reuse views of the exceptions list, including the current search, filters, time range, status, sort and page size.
---

# oneuptime_exception_saved_view (Resource)

Save and reuse views of the exceptions list, including the current search, filters, time range, status, sort and page size.

## Example Usage

```terraform
resource "oneuptime_exception_saved_view" "example" {
  name = "Example exception saved view"
}
```

## Schema

### Required

- `name` (String) Friendly name for this saved exceptions view.

### Optional

- `is_default` (Boolean) Whether this saved exceptions view should be applied by default. Defaults to `false`.
- `query` (String) Serialized exceptions list view state (search, filters, time range, page size, status, lens, assignee filter and sort) for this saved view. A JSON value: write it with `jsonencode()`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this saved exceptions view. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this saved exceptions view belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing exception saved view by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_exception_saved_view.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_exception_saved_view.example <id>
```
