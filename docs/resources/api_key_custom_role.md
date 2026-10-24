---
page_title: "oneuptime_api_key_custom_role Resource - oneuptime"
subcategory: "Teams & Access"
description: |-
  The custom roles an API key is given. The key holds every permission of the role except those with the Owned scope: a key owns nothing.
---

# oneuptime_api_key_custom_role (Resource)

The custom roles an API key is given. The key holds every permission of the role except those with the Owned scope: a key owns nothing.

## Example Usage

```terraform
resource "oneuptime_api_key_custom_role" "example" {
  api_key_id     = oneuptime_api_key.example.id
  custom_role_id = oneuptime_custom_role.example.id
}
```

## Schema

### Required

- `api_key_id` (String) ID of the API key given the role. The ID of a `oneuptime_api_key`.
- `custom_role_id` (String) ID of the custom role given to the API key. The ID of a `oneuptime_custom_role`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing api key custom role by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_api_key_custom_role.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_api_key_custom_role.example <id>
```
