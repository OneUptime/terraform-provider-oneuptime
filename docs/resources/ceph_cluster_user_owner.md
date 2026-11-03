---
page_title: "oneuptime_ceph_cluster_user_owner Resource - oneuptime"
subcategory: "Infrastructure"
description: |-
  Add users as owners to your Ceph clusters.
---

# oneuptime_ceph_cluster_user_owner (Resource)

Add users as owners to your Ceph clusters.

## Example Usage

```terraform
resource "oneuptime_ceph_cluster_user_owner" "example" {
  ceph_cluster_id = oneuptime_ceph_cluster.example.id
  user_email      = "alice@example.com"
}
```

## Schema

### Required

- `ceph_cluster_id` (String) ID of the Ceph cluster the user owns. The ID of a `oneuptime_ceph_cluster`.

### Optional

- `user_email` (String) The email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) ID of the user who owns the Ceph cluster. The ID of a `oneuptime_user` (see the data source). Or name them by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing ceph cluster user owner by its ceph cluster and the person's email address (or user ID), with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_ceph_cluster_user_owner.example
  id = "<ceph_cluster_id>/alice@example.com"
}
```

or on the command line:

```shell
terraform import oneuptime_ceph_cluster_user_owner.example <ceph_cluster_id>/alice@example.com
```

Imported by an address, the person's `user_email` stays as you wrote it. Its own ID works too.
