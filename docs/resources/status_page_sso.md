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
  status_page_id = "123e4567-e89b-12d3-a456-426614174000"
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
  signature_method = "Example short text"
  digest_method = "Example short text"
  sign_on_url = "https://www.example.com/path/to/resource?param=value"
  issuer_url = "This is an example of very long text content that might be stored in this field. It can contain a lot of information, such as detailed descriptions, comments, or any other lengthy text data that needs to be stored in the database."
  public_certificate = "This is an example of very long text content that might be stored in this field. It can contain a lot of information, such as detailed descriptions, comments, or any other lengthy text data that needs to be stored in the database."
}
```

## Schema

### Required

- `status_page_id` (String) A unique identifier for an object, represented as a UUID..
- `name` (String) Any friendly name of this object..
- `description` (String) Status page sso description.
- `signature_method` (String) Status page sso signature_method.
- `digest_method` (String) Status page sso digest_method.
- `sign_on_url` (String) Status page sso sign_on_url.
- `issuer_url` (String) Status page sso issuer_url.
- `public_certificate` (String) Status page sso public_certificate.

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_enabled` (Bool) Status page sso is_enabled.
- `is_tested` (Bool) Status page sso is_tested.

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
terraform import oneuptime_status_page_sso.example <id>
```
