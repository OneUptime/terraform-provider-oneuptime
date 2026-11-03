---
page_title: "oneuptime_kubernetes_cluster_team_owner Resource - oneuptime"
subcategory: "Infrastructure"
description: |-
  Add teams as owners to your Kubernetes clusters.
---

# oneuptime_kubernetes_cluster_team_owner (Resource)

Add teams as owners to your Kubernetes clusters.

## Example Usage

```terraform
resource "oneuptime_kubernetes_cluster_team_owner" "example" {
  team_id               = oneuptime_team.example.id
  kubernetes_cluster_id = oneuptime_kubernetes_cluster.example.id
}
```

## Schema

### Required

- `kubernetes_cluster_id` (String) ID of the Kubernetes cluster the team owns. The ID of a `oneuptime_kubernetes_cluster`.
- `team_id` (String) ID of the team that owns the Kubernetes cluster. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing kubernetes cluster team owner by its kubernetes cluster and its `team_id`, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_kubernetes_cluster_team_owner.example
  id = "<kubernetes_cluster_id>/<team_id>"
}
```

or on the command line:

```shell
terraform import oneuptime_kubernetes_cluster_team_owner.example <kubernetes_cluster_id>/<team_id>
```

Its own ID works too.
