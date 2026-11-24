---
page_title: "oneuptime_form_private_user Resource - oneuptime"
subcategory: "Forms"
description: |-
  People who can sign in to open a form's public link when it asks its visitors to sign in.
---

# oneuptime_form_private_user (Resource)

People who can sign in to open a form's public link when it asks its visitors to sign in.

## Example Usage

```terraform
resource "oneuptime_form_private_user" "example" {
  form_id = oneuptime_form.example.id
  email   = "user@example.com"
}
```

## Schema

### Required

- `email` (String) The address this private user signs in with.
- `form_id` (String) ID of the form this private user may sign in to. The ID of a `oneuptime_form`.

### Optional

- `password` (String, Sensitive) The private user's password, stored as a salted hash and never shown.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `has_password` (Boolean) Whether this private user has set a password.
- `id` (String) Unique identifier for the resource.
- `is_sso_user` (Boolean) Whether the form's single sign-on added this user.
- `last_active` (String) When this private user last signed in to the form.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing form private user by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_form_private_user.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_form_private_user.example <id>
```
