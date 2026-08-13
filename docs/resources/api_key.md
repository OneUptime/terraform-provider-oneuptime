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
  name = "Example short text"
  expires_at = "2030-01-01T00:00:00Z"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..
- `expires_at` (String) A date time object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `api_key` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_api_key.example <id>
```
