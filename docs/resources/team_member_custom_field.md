---
page_title: "oneuptime_team_member_custom_field Resource - oneuptime"
subcategory: "Teams & Access"
description: |-
  Manage custom fields for your team members.
---

# oneuptime_team_member_custom_field (Resource)

Manage custom fields for your team members.

## Example Usage

```terraform
resource "oneuptime_team_member_custom_field" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description of this custom field that will help you remember..
- `custom_field_type` (String) Is this field Text, Number or Boolean?..
- `dropdown_options` (String) Options and optional colors for dropdown fields. Plain one-per-line values remain supported...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_team_member_custom_field.example <id>
```
