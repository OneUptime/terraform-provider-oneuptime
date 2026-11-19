---
page_title: "oneuptime_pipeline_user_owner Resource - oneuptime"
subcategory: "Other"
description: |-
  Add users as owners to your data pipelines.
---

# oneuptime_pipeline_user_owner (Resource)

Add users as owners to your data pipelines.

## Example Usage

```terraform
resource "oneuptime_pipeline_user_owner" "example" {
  data_pipeline_id = data.oneuptime_data_pipeline.example.id
  user_email       = "alice@example.com"
}
```

## Schema

### Required

- `data_pipeline_id` (String) ID of the data pipeline the user owns. The ID of a `oneuptime_data_pipeline` (see the data source).

### Optional

- `user_email` (String) The email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) ID of the user who owns the data pipeline. The ID of a `oneuptime_user` (see the data source). Or name them by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing pipeline user owner by its data pipeline and the person's email address (or user ID), with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_pipeline_user_owner.example
  id = "<data_pipeline_id>/alice@example.com"
}
```

or on the command line:

```shell
terraform import oneuptime_pipeline_user_owner.example <data_pipeline_id>/alice@example.com
```

Imported by an address, the person's `user_email` stays as you wrote it. Its own ID works too.
