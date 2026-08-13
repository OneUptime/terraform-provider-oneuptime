---
page_title: "oneuptime_status_page_scim Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage SCIM auto-provisioning for your status page
---

# oneuptime_status_page_scim (Resource)

Manage SCIM auto-provisioning for your status page

## Example Usage

```terraform
resource "oneuptime_status_page_scim" "example" {
  status_page_id = "123e4567-e89b-12d3-a456-426614174000"
  name = "Example short text"
  bearer_token = "This is an example of longer text content that might be stored in this field."
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `status_page_id` (String) A unique identifier for an object, represented as a UUID..
- `name` (String) Any friendly name for this SCIM configuration..
- `bearer_token` (String) Bearer token for SCIM authentication. Keep this secure...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description to help you remember..
- `auto_provision_users` (Bool) Automatically create status page users when they are added via SCIM..
- `auto_deprovision_users` (Bool) Automatically remove status page users when they are removed via SCIM..
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
terraform import oneuptime_status_page_scim.example <id>
```
