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
  domain = jsonencode({
    "_type": "Domain",
    "value": "example.com"
  })
}
```

## Schema

### Required

- `domain` (String) Domain object.

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_verified` (Bool) Is this domain verified?..
- `domain_verification_text` (String) Verification text that you need to add to your domains TXT record to veify the domain...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_domain.example <id>
```
