---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_domain Resource - oneuptime"
subcategory: ""
description: |-
  Domain resource
---

# oneuptime_domain (Resource)

Domain resource

## Example Usage

```terraform
resource "oneuptime_domain" "example" {
  domain = {
    id = "123e4567-e89b-12d3-a456-426614174000"
  }
}
```

## Schema

- `id` (String) Unique identifier for the resource. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Optional.
- `domain` (Map) Domain object. Required.
- `domain_verification_text` (String) Domain Verification Text. Optional.
- `is_verified` (Bool) Verified. Optional.
- `created_at` (Map) A date time object.. Computed.
- `updated_at` (Map) A date time object.. Computed.
- `deleted_at` (Map) A date time object.. Computed.
- `version` (Number) Version. Computed.
- `slug` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Domain], Update: [No access - you don't have permission for this operation]. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_domain.example <id>
```
