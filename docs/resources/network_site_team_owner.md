---
page_title: "oneuptime_network_site_team_owner Resource - oneuptime"
subcategory: "Network"
description: |-
  Add teams as owners to your network sites.
---

# oneuptime_network_site_team_owner (Resource)

Add teams as owners to your network sites.

## Example Usage

```terraform
resource "oneuptime_network_site_team_owner" "example" {
  team_id         = oneuptime_team.example.id
  network_site_id = oneuptime_network_site.example.id
}
```

## Schema

### Required

- `network_site_id` (String) ID of the network site the team owns. The ID of a `oneuptime_network_site`.
- `team_id` (String) ID of the team that owns the network site. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing network site team owner by its network site and its `team_id`, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_network_site_team_owner.example
  id = "<network_site_id>/<team_id>"
}
```

or on the command line:

```shell
terraform import oneuptime_network_site_team_owner.example <network_site_id>/<team_id>
```

Its own ID works too.
