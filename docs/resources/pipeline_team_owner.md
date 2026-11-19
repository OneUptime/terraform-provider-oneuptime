---
page_title: "oneuptime_pipeline_team_owner Resource - oneuptime"
subcategory: "Other"
description: |-
  Add teams as owners to your data pipelines.
---

# oneuptime_pipeline_team_owner (Resource)

Add teams as owners to your data pipelines.

## Example Usage

```terraform
resource "oneuptime_pipeline_team_owner" "example" {
  team_id          = oneuptime_team.example.id
  data_pipeline_id = data.oneuptime_data_pipeline.example.id
}
```

## Schema

### Required

- `data_pipeline_id` (String) ID of the data pipeline the team owns. The ID of a `oneuptime_data_pipeline` (see the data source).
- `team_id` (String) ID of the team that owns the data pipeline. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing pipeline team owner by its data pipeline and its `team_id`, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_pipeline_team_owner.example
  id = "<data_pipeline_id>/<team_id>"
}
```

or on the command line:

```shell
terraform import oneuptime_pipeline_team_owner.example <data_pipeline_id>/<team_id>
```

Its own ID works too.
