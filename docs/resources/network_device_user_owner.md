---
page_title: "oneuptime_network_device_user_owner Resource - oneuptime"
subcategory: "Network"
description: |-
  Add users as owners to your network devices.
---

# oneuptime_network_device_user_owner (Resource)

Add users as owners to your network devices.

## Example Usage

```terraform
resource "oneuptime_network_device_user_owner" "example" {
  network_device_id = oneuptime_network_device.example.id
  user_email        = "alice@example.com"
}
```

## Schema

### Required

- `network_device_id` (String) ID of the network device the user owns. The ID of a `oneuptime_network_device`.

### Optional

- `user_email` (String) The email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) ID of the user who owns the network device. The ID of a `oneuptime_user` (see the data source). Or name them by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing network device user owner by its network device and the person's email address (or user ID), with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_network_device_user_owner.example
  id = "<network_device_id>/alice@example.com"
}
```

or on the command line:

```shell
terraform import oneuptime_network_device_user_owner.example <network_device_id>/alice@example.com
```

Imported by an address, the person's `user_email` stays as you wrote it. Its own ID works too.
