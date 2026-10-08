---
page_title: "oneuptime_table_view Resource - oneuptime"
subcategory: "Other"
description: |-
  Table View is view settings for a table in a project. It contains columns, filters, and other settings.
---

# oneuptime_table_view (Resource)

Table View is view settings for a table in a project. It contains columns, filters, and other settings.

## Example Usage

```terraform
resource "oneuptime_table_view" "example" {
  name        = "Example table view"
  table_id    = "Example short text"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object.
- `table_id` (String) ID of the table this view is for.

### Optional

- `columns` (String) Which columns are shown, and in what order, for this table view. A JSON value: write it with `jsonencode()`.
- `description` (String) Friendly description that will help you remember.
- `facets` (String) Facet selections (owner, labels, status, etc.) for this table view. A JSON value: write it with `jsonencode()`.
- `items_on_page` (Number) Items on page. Defaults to `10`.
- `query` (String) Filters for this table view. A JSON value: write it with `jsonencode()`.
- `sort` (String) Sort for this table view. A JSON value: write it with `jsonencode()`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing table view by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_table_view.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_table_view.example <id>
```
