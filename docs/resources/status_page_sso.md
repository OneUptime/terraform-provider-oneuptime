---
page_title: "oneuptime_status_page_sso Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Configure Status Page SSO
---

# oneuptime_status_page_sso (Resource)

Configure Status Page SSO

## Example Usage

```terraform
resource "oneuptime_status_page_sso" "example" {
  status_page_id     = oneuptime_status_page.example.id
  name               = "Example status page sso"
  description        = "Managed by Terraform"
  signature_method   = "Example short text"
  digest_method      = "Example short text"
  sign_on_url        = "https://www.example.com/path/to/resource?param=value"
  issuer_url         = "This is an example of very long text content that might be stored in this field. It can contain a lot of information, such as detailed descriptions, comments, or any other lengthy text data that needs to be stored in the database."
  public_certificate = "This is an example of very long text content that might be stored in this field. It can contain a lot of information, such as detailed descriptions, comments, or any other lengthy text data that needs to be stored in the database."
}
```

## Schema

### Required

- `description` (String)
- `digest_method` (String)
- `issuer_url` (String)
- `name` (String) Any friendly name of this object.
- `public_certificate` (String)
- `sign_on_url` (String)
- `signature_method` (String)
- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.

### Optional

- `is_enabled` (Boolean) Defaults to `false`.
- `is_tested` (Boolean) Defaults to `false`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing status page sso by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_status_page_sso.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_status_page_sso.example <id>
```
