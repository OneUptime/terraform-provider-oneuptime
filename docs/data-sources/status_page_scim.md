---
page_title: "oneuptime_status_page_scim Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage SCIM auto-provisioning for your status page
---

# oneuptime_status_page_scim (Data Source)

Manage SCIM auto-provisioning for your status page Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_status_page_scim" "by_name" {
  name = "example-status_page_scim"
}

data "oneuptime_status_page_scim" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `status_page_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description` (String) Friendly description to help you remember.. Computed.
- `bearer_token` (String) Bearer token for SCIM authentication. Keep this secure... Computed.
- `auto_provision_users` (Bool) Automatically create status page users when they are added via SCIM.. Computed.
- `auto_deprovision_users` (Bool) Automatically remove status page users when they are removed via SCIM.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
