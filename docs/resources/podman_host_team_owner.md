---
page_title: "oneuptime_podman_host_team_owner Resource - oneuptime"
subcategory: "Infrastructure"
description: |-
  Add teams as owners to your Podman hosts.
---

# oneuptime_podman_host_team_owner (Resource)

Add teams as owners to your Podman hosts.

## Example Usage

```terraform
resource "oneuptime_podman_host_team_owner" "example" {
  team_id        = oneuptime_team.example.id
  podman_host_id = oneuptime_podman_host.example.id
}
```

## Schema

### Required

- `podman_host_id` (String) ID of the Podman host the team owns. The ID of a `oneuptime_podman_host`.
- `team_id` (String) ID of the team that owns the Podman host. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing podman host team owner by its podman host and its `team_id`, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_podman_host_team_owner.example
  id = "<podman_host_id>/<team_id>"
}
```

or on the command line:

```shell
terraform import oneuptime_podman_host_team_owner.example <podman_host_id>/<team_id>
```

Its own ID works too.
