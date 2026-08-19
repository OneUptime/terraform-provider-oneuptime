---
page_title: "oneuptime_status_page_oidc Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage OpenID Connect (OIDC) authentication for your status page
---

# oneuptime_status_page_oidc (Resource)

Manage OpenID Connect (OIDC) authentication for your status page

## Example Usage

```terraform
resource "oneuptime_status_page_oidc" "example" {
  status_page_id = "123e4567-e89b-12d3-a456-426614174000"
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
  discovery_url = "https://www.example.com/path/to/resource?param=value"
  issuer_url = "This is an example of very long text content that might be stored in this field. It can contain a lot of information, such as detailed descriptions, comments, or any other lengthy text data that needs to be stored in the database."
  client_id = "Example short text"
  client_secret = "This is an example of longer text content that might be stored in this field."
  scopes = "Example short text"
  email_claim_name = "Example short text"
}
```

## Schema

### Required

- `status_page_id` (String) A unique identifier for an object, represented as a UUID..
- `name` (String) Any friendly name of this object..
- `description` (String) Status page oidc description.
- `discovery_url` (String) OIDC discovery URL (typically ends in /.well-known/openid-configuration). Used to discover authorization, token, JWKS and userinfo endpoints...
- `issuer_url` (String) Expected OIDC issuer URL. Must match the 'iss' claim in the ID token returned by the identity provider...
- `client_id` (String) OIDC client ID issued by the identity provider...
- `client_secret` (String) OIDC client secret issued by the identity provider. Stored encrypted at rest...
- `scopes` (String) Space-separated list of OIDC scopes to request. Must include 'openid'...
- `email_claim_name` (String) Claim name in the ID token (or userinfo response) that contains the user's email address...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `name_claim_name` (String) Claim name in the ID token (or userinfo response) that contains the user's display name...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_enabled` (Bool) Status page oidc is_enabled.
- `is_tested` (Bool) Status page oidc is_tested.

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
terraform import oneuptime_status_page_oidc.example <id>
```
