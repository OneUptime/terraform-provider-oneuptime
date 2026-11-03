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

- `service_level_objective_id` (String) ID of the service level objective the team owns. The ID of a `oneuptime_service_level_objective`.
- `team_id` (String) ID of the team that owns the service level objective. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing service level objective team owner by its service level objective and its `team_id`, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_service_level_objective_team_owner.example
  id = "<service_level_objective_id>/<team_id>"
}
```

or on the command line:

```shell
terraform import oneuptime_service_level_objective_team_owner.example <service_level_objective_id>/<team_id>
```

Its own ID works too.
