---
page_title: "oneuptime_service_level_objective_team_owner Resource - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Add teams as owners to your Service Level Objectives.
---

# oneuptime_service_level_objective_team_owner (Resource)

Add teams as owners to your Service Level Objectives.

## Example Usage

```terraform
resource "oneuptime_service_level_objective_team_owner" "example" {
  team_id                    = oneuptime_team.example.id
  service_level_objective_id = oneuptime_service_level_objective.example.id
}
```

## Schema

### Required

- `service_level_objective_id` (String) ID of your OneUptime Service Level Objective in which this object belongs. The ID of a `oneuptime_service_level_objective`.
- `team_id` (String) ID of your OneUptime Team in which this object belongs. The ID of a `oneuptime_team`.

### Optional

- `is_owner_notified` (Boolean) Are owners notified of this resource ownership? Defaults to `false`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing service level objective team owner by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_service_level_objective_team_owner.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_service_level_objective_team_owner.example <id>
```
