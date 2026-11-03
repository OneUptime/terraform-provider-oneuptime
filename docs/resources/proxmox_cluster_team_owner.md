---
page_title: "oneuptime_proxmox_cluster_team_owner Resource - oneuptime"
subcategory: "Infrastructure"
description: |-
  Add teams as owners to your Proxmox clusters.
---

# oneuptime_proxmox_cluster_team_owner (Resource)

Add teams as owners to your Proxmox clusters.

## Example Usage

```terraform
resource "oneuptime_proxmox_cluster_team_owner" "example" {
  team_id            = oneuptime_team.example.id
  proxmox_cluster_id = oneuptime_proxmox_cluster.example.id
}
```

## Schema

### Required

- `proxmox_cluster_id` (String) ID of the Proxmox cluster the team owns. The ID of a `oneuptime_proxmox_cluster`.
- `team_id` (String) ID of the team that owns the Proxmox cluster. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing proxmox cluster team owner by its proxmox cluster and its `team_id`, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_proxmox_cluster_team_owner.example
  id = "<proxmox_cluster_id>/<team_id>"
}
```

or on the command line:

```shell
terraform import oneuptime_proxmox_cluster_team_owner.example <proxmox_cluster_id>/<team_id>
```

Its own ID works too.
