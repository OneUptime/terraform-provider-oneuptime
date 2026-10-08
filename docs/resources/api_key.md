---
page_title: "oneuptime_api_key Resource - oneuptime"
subcategory: "Teams & Access"
description: |-
  Manage API Keys for your project
---

# oneuptime_api_key (Resource)

Manage API Keys for your project

## Example Usage

```terraform
resource "oneuptime_api_key" "example" {
  name        = "Example api key"
  expires_at  = "2030-01-01T00:00:00Z"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `expires_at` (String) Date and Time when this API Key expires.
- `name` (String) Any friendly name of this object.

### Optional

- `description` (String) Friendly description that will help you remember.

### Read-Only

- `api_key` (String) Secret API Key.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing api key by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_api_key.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_api_key.example <id>
```
