---
page_title: "oneuptime_service_user_owner Resource - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Add users as owners to your Service.
---

# oneuptime_service_user_owner (Resource)

Add users as owners to your Service.

## Example Usage

```terraform
resource "oneuptime_service_user_owner" "example" {
  user_id    = data.oneuptime_user.example.id
  service_id = oneuptime_service.example.id
}
```

## Schema

### Required

- `service_id` (String) ID of your OneUptime Service in which this object belongs. The ID of a `oneuptime_service`.
- `user_id` (String) ID of your OneUptime User in which this object belongs. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing service user owner by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_service_user_owner.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_service_user_owner.example <id>
```
