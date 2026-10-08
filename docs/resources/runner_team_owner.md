---
page_title: "oneuptime_runner_team_owner Resource - oneuptime"
subcategory: "Other"
description: |-
  Add teams as owners to your runbook agents.
---

# oneuptime_runner_team_owner (Resource)

Add teams as owners to your runbook agents.

## Example Usage

```terraform
resource "oneuptime_runner_team_owner" "example" {
  runner_id = oneuptime_runner.example.id
  team_id   = oneuptime_team.example.id
}
```

## Schema

### Required

- `runner_id` (String) ID of your OneUptime Runner in which this object belongs. The ID of a `oneuptime_runner`.
- `team_id` (String) ID of your OneUptime Team in which this object belongs. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing runner team owner by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_runner_team_owner.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_runner_team_owner.example <id>
```
