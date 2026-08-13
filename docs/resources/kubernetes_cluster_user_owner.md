---
page_title: "oneuptime_kubernetes_cluster_user_owner Resource - oneuptime"
subcategory: "Other"
description: |-
  Add users as owners to your Kubernetes clusters.
---

# oneuptime_kubernetes_cluster_user_owner (Resource)

Add users as owners to your Kubernetes clusters.

## Example Usage

```terraform
resource "oneuptime_kubernetes_cluster_user_owner" "example" {
  user_id = "123e4567-e89b-12d3-a456-426614174000"
  kubernetes_cluster_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `user_id` (String) A unique identifier for an object, represented as a UUID..
- `kubernetes_cluster_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_owner_notified` (Bool) Are owners notified of this resource ownership?..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_kubernetes_cluster_user_owner.example <id>
```
