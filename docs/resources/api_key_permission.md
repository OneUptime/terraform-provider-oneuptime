---
page_title: "oneuptime_api_key_permission Resource - oneuptime"
subcategory: "Teams & Access"
description: |-
  Permissions for your API Keys
---

# oneuptime_api_key_permission (Resource)

Permissions for your API Keys

## Example Usage

```terraform
resource "oneuptime_api_key_permission" "example" {
  api_key_id = oneuptime_api_key.example.id
}
```

## Schema

### Required

- `api_key_id` (String) ID of API Key resource in which this object belongs. The ID of a `oneuptime_api_key`.

### Optional

- `is_block_permission` (Boolean) Defaults to `false`.
- `labels` (Set of String) Relation to Labels Array where this permission is scoped at. IDs of `oneuptime_label` resources.
- `permission` (String) Permission. You can find list of permissions on the Permissions page. A JSON value: write it with `jsonencode()`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing api key permission by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_api_key_permission.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_api_key_permission.example <id>
```
