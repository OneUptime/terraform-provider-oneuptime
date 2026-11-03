---
page_title: "oneuptime_probe_team_owner Resource - oneuptime"
subcategory: "Probes"
description: |-
  Add teams as owners to your probes.
---

# oneuptime_probe_team_owner (Resource)

Add teams as owners to your probes.

~> **Renamed:** this resource was called `oneuptime_probe_owner_team` before. The old name still works, but is deprecated. To switch, rename the resource in your configuration and add a `moved` block, so Terraform keeps the existing probe team owner:

```terraform
moved {
  from = oneuptime_probe_owner_team.example
  to   = oneuptime_probe_team_owner.example
}
```

## Example Usage

```terraform
resource "oneuptime_probe_team_owner" "example" {
  team_id  = oneuptime_team.example.id
  probe_id = oneuptime_probe.example.id
}
```

## Schema

### Required

- `probe_id` (String) ID of the probe the team owns. The ID of a `oneuptime_probe`.
- `team_id` (String) ID of the team that owns the probe. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing probe team owner by its probe and its `team_id`, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_probe_team_owner.example
  id = "<probe_id>/<team_id>"
}
```

or on the command line:

```shell
terraform import oneuptime_probe_team_owner.example <probe_id>/<team_id>
```

Its own ID works too.
