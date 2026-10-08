---
page_title: "oneuptime_status_page_scim Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage SCIM auto-provisioning for your status page
---

# oneuptime_status_page_scim (Data Source)

Manage SCIM auto-provisioning for your status page

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one status page scim may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_status_page_scim" "example" {
  name = "Example status page scim"
}

# Or by id:
data "oneuptime_status_page_scim" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `auto_deprovision_users` (Boolean) Automatically remove status page users when they are removed via SCIM.
- `auto_provision_users` (Boolean) Automatically create status page users when they are added via SCIM.
- `bearer_token` (String) Bearer token for SCIM authentication. Keep this secure.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description to help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `name` (String) Any friendly name for this SCIM configuration.
- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
