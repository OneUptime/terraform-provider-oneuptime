---
page_title: "oneuptime_vcenter_team_owner Resource - oneuptime"
subcategory: "Other"
description: |-
  Add teams as owners to your vCenters.
---

# oneuptime_vcenter_team_owner (Resource)

Add teams as owners to your vCenters.

~> **Renamed:** this resource was called `oneuptime_v_center_team_owner` before. The old name still works, but is deprecated. To switch, rename the resource in your configuration and add a `moved` block, so Terraform keeps the existing vcenter team owner:

```terraform
moved {
  from = oneuptime_v_center_team_owner.example
  to   = oneuptime_vcenter_team_owner.example
}
```

## Example Usage

```terraform
resource "oneuptime_vcenter_team_owner" "example" {
  team_id            = oneuptime_team.example.id
  vmware_v_center_id = oneuptime_vcenter.example.id
}
```

## Schema

### Required

- `team_id` (String) ID of your OneUptime Team in which this object belongs. The ID of a `oneuptime_team`.
- `vmware_v_center_id` (String) ID of your OneUptime vCenter in which this object belongs. The ID of a `oneuptime_vcenter`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing vcenter team owner by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_vcenter_team_owner.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_vcenter_team_owner.example <id>
```
