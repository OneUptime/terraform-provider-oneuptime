---
page_title: "oneuptime_status_page_private_user Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage private users on your status page
---

# oneuptime_status_page_private_user (Resource)

Manage private users on your status page

## Example Usage

```terraform
resource "oneuptime_status_page_private_user" "example" {
  status_page_id = oneuptime_status_page.example.id
}
```

## Schema

### Required

- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.

### Optional

- `email` (String) Email.
- `is_sso_user` (Boolean) Did this user sign up via SSO? Defaults to `false`.
- `password` (String, Sensitive) Password.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing status page private user by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_status_page_private_user.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_status_page_private_user.example <id>
```
