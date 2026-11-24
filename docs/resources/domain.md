---
page_title: "oneuptime_domain Resource - oneuptime"
subcategory: "Organization"
description: |-
  Manage Custom Domains for your project
---

# oneuptime_domain (Resource)

Manage Custom Domains for your project

## Example Usage

```terraform
resource "oneuptime_domain" "example" {
  domain = "example.com"
}
```

## Schema

### Required

- `domain` (String) Domain - acmeinc.com for example.

### Optional

- `is_verified` (Boolean) Is this domain verified? OneUptime keeps it up to date once the record exists: the value given here is the one the record starts with, and a later change of it is never sent. Defaults to `false`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `domain_verification_text` (String) Verification text that you need to add to your domains TXT record to veify the domain.
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing domain by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_domain.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_domain.example <id>
```
