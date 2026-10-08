---
page_title: "oneuptime_domain Data Source - oneuptime"
subcategory: "Organization"
description: |-
  Manage Custom Domains for your project
---

# oneuptime_domain (Data Source)

Manage Custom Domains for your project

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one domain may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_domain" "example" {
  slug = "example-slug"
}

# Or by id:
data "oneuptime_domain" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `domain_verification_text` (String) Verification text that you need to add to your domains TXT record to veify the domain.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_verified` (Boolean) Is this domain verified?
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `domain` (String) Domain - acmeinc.com for example.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
