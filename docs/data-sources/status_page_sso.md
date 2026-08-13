---
page_title: "oneuptime_status_page_sso Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Configure Status Page SSO
---

# oneuptime_status_page_sso (Data Source)

Configure Status Page SSO Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_status_page_sso" "by_name" {
  name = "example-status_page_sso"
}

data "oneuptime_status_page_sso" "by_id" {
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
- `description` (String) Status page sso description. Computed.
- `signature_method` (String) Status page sso signature_method. Computed.
- `digest_method` (String) Status page sso digest_method. Computed.
- `sign_on_url` (String) Status page sso sign_on_url. Computed.
- `issuer_url` (String) Status page sso issuer_url. Computed.
- `public_certificate` (String) Status page sso public_certificate. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_enabled` (Bool) Status page sso is_enabled. Computed.
- `is_tested` (Bool) Status page sso is_tested. Computed.
