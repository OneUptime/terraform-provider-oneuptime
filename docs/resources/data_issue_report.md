---
page_title: "oneuptime_data_issue_report Resource - oneuptime"
subcategory: "Other"
description: |-
  Someone saying a table's data looks wrong, and whether OneUptime had already raised it.
---

# oneuptime_data_issue_report (Resource)

Someone saying a table's data looks wrong, and whether OneUptime had already raised it.

## Example Usage

```terraform
resource "oneuptime_data_issue_report" "example" {
  data_asset_id = oneuptime_data_asset.example.id
  note          = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `data_asset_id` (String) ID of the table whose data looks wrong. The ID of a `oneuptime_data_asset`.
- `note` (String) What looks wrong, in the words of whoever noticed.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of who reported it. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this report belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
- `was_caught_first` (Boolean) Whether an alert on this table, or a table upstream of it, was already open when this was reported.

## Import

Import an existing data issue report by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_data_issue_report.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_data_issue_report.example <id>
```
