---
page_title: "oneuptime_openstack_cloud_team_owner Resource - oneuptime"
subcategory: "Infrastructure"
description: |-
  Add teams as owners to your OpenStack clouds.
---

# oneuptime_openstack_cloud_team_owner (Resource)

Add teams as owners to your OpenStack clouds.

~> **Renamed:** this resource was called `oneuptime_open_stack_cloud_team_owner` before. The old name still works, but is deprecated. To switch, rename the resource in your configuration and add a `moved` block, so Terraform keeps the existing openstack cloud team owner:

```terraform
moved {
  from = oneuptime_open_stack_cloud_team_owner.example
  to   = oneuptime_openstack_cloud_team_owner.example
}
```

## Example Usage

```terraform
resource "oneuptime_openstack_cloud_team_owner" "example" {
  team_id             = oneuptime_team.example.id
  open_stack_cloud_id = oneuptime_openstack_cloud.example.id
}
```

## Schema

### Required

- `open_stack_cloud_id` (String) ID of the OpenStack cloud the team owns. The ID of a `oneuptime_openstack_cloud`.
- `team_id` (String) ID of the team that owns the OpenStack cloud. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing openstack cloud team owner by its openstack cloud and its `team_id`, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_openstack_cloud_team_owner.example
  id = "<open_stack_cloud_id>/<team_id>"
}
```

or on the command line:

```shell
terraform import oneuptime_openstack_cloud_team_owner.example <open_stack_cloud_id>/<team_id>
```

Its own ID works too.
