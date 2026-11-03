---
page_title: "oneuptime_data_store_user_owner Resource - oneuptime"
subcategory: "Data Security"
description: |-
  Add users as owners to your data stores.
---

# oneuptime_data_store_user_owner (Resource)

Add users as owners to your data stores.

## Example Usage

```terraform
resource "oneuptime_data_store_user_owner" "example" {
  data_store_id = oneuptime_data_store.example.id
  user_email    = "alice@example.com"
}
```

## Schema

### Required

- `data_store_id` (String) ID of the data store the user owns. The ID of a `oneuptime_data_store`.

### Optional

- `user_email` (String) The email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) ID of the user who owns the data store. The ID of a `oneuptime_user` (see the data source). Or name them by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing data store user owner by its data store and the person's email address (or user ID), with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_data_store_user_owner.example
  id = "<data_store_id>/alice@example.com"
}
```

or on the command line:

```shell
terraform import oneuptime_data_store_user_owner.example <data_store_id>/alice@example.com
```

Imported by an address, the person's `user_email` stays as you wrote it. Its own ID works too.
