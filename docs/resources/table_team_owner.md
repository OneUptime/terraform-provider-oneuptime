---
page_title: "oneuptime_table_team_owner Resource - oneuptime"
subcategory: "Other"
description: |-
  Add teams as owners to your tables.
---

# oneuptime_table_team_owner (Resource)

Add teams as owners to your tables.

## Example Usage

```terraform
resource "oneuptime_table_team_owner" "example" {
  team_id       = oneuptime_team.example.id
  data_asset_id = oneuptime_data_asset.example.id
}
```

## Schema

### Required

- `data_asset_id` (String) ID of the table the team owns. The ID of a `oneuptime_data_asset`.
- `team_id` (String) ID of the team that owns the table. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing table team owner by its data asset and its `team_id`, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_table_team_owner.example
  id = "<data_asset_id>/<team_id>"
}
```

or on the command line:

```shell
terraform import oneuptime_table_team_owner.example <data_asset_id>/<team_id>
```

Its own ID works too.
