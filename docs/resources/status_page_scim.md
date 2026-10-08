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
  status_page_id = oneuptime_status_page.example.id
  name           = "Example status page scim"
  bearer_token   = "This is an example of longer text content that might be stored in this field."
  description    = "Managed by Terraform"
}
```

## Schema

### Required

- `bearer_token` (String) Bearer token for SCIM authentication. Keep this secure.
- `name` (String) Any friendly name for this SCIM configuration.
- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.

### Optional

- `auto_deprovision_users` (Boolean) Automatically remove status page users when they are removed via SCIM. Defaults to `true`.
- `auto_provision_users` (Boolean) Automatically create status page users when they are added via SCIM. Defaults to `true`.
- `description` (String) Friendly description to help you remember.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing status page scim by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_status_page_scim.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_status_page_scim.example <id>
```
