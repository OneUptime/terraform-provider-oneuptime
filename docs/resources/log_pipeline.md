---
page_title: "oneuptime_log_pipeline Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure server-side log processing pipelines that transform logs at ingest time.
---

# oneuptime_log_pipeline (Resource)

Configure server-side log processing pipelines that transform logs at ingest time.

## Example Usage

```terraform
resource "oneuptime_log_pipeline" "example" {
  name        = "Example log pipeline"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Friendly name for this log pipeline.

### Optional

- `description` (String) Description of what this log pipeline does.
- `filter_query` (String) Filter expression that determines which logs this pipeline applies to.
- `is_enabled` (Boolean) Whether this log pipeline is active. Defaults to `true`.
- `sort_order` (Number) Where this pipeline runs among the project's log pipelines, lowest number first. A new pipeline is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this log pipeline. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this log pipeline belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing log pipeline by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_log_pipeline.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_log_pipeline.example <id>
```
