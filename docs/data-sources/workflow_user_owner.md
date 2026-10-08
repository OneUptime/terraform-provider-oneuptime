---
page_title: "oneuptime_workflow_user_owner Data Source - oneuptime"
subcategory: "Workflows"
description: |-
  Add users as owners to your workflows.
---

# oneuptime_workflow_user_owner (Data Source)

Add users as owners to your workflows.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one workflow user owner may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_workflow_user_owner" "example" {
  user_id = data.oneuptime_user.example.id
}

# Or by id:
data "oneuptime_workflow_user_owner" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `user_id` (String) ID of your OneUptime User in which this object belongs. The ID of a `oneuptime_user` (see the data source).
- `workflow_id` (String) ID of your OneUptime Workflow in which this object belongs. The ID of a `oneuptime_workflow`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
