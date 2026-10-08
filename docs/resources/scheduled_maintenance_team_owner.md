---
page_title: "oneuptime_scheduled_maintenance_team_owner Resource - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Add teams as owners to your Scheduled Maintenance event.
---

# oneuptime_scheduled_maintenance_team_owner (Resource)

Add teams as owners to your Scheduled Maintenance event.

## Example Usage

```terraform
resource "oneuptime_scheduled_maintenance_team_owner" "example" {
  team_id                  = oneuptime_team.example.id
  scheduled_maintenance_id = oneuptime_scheduled_maintenance_event.example.id
}
```

## Schema

### Required

- `scheduled_maintenance_id` (String) ID of your OneUptime ScheduledMaintenance in which this object belongs. The ID of a `oneuptime_scheduled_maintenance_event`.
- `team_id` (String) ID of your OneUptime Team in which this object belongs. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing scheduled maintenance team owner by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_scheduled_maintenance_team_owner.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_scheduled_maintenance_team_owner.example <id>
```
