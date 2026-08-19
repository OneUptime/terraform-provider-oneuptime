---
page_title: "oneuptime_status_page_oidc Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage OpenID Connect (OIDC) authentication for your status page
---

# oneuptime_status_page_oidc (Data Source)

Manage OpenID Connect (OIDC) authentication for your status page Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_status_page_oidc" "by_name" {
  name = "example-status_page_oidc"
}

data "oneuptime_status_page_oidc" "by_id" {
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
- `description` (String) Status page oidc description. Computed.
- `discovery_url` (String) OIDC discovery URL (typically ends in /.well-known/openid-configuration). Used to discover authorization, token, JWKS and userinfo endpoints... Computed.
- `issuer_url` (String) Expected OIDC issuer URL. Must match the 'iss' claim in the ID token returned by the identity provider... Computed.
- `client_id` (String) OIDC client ID issued by the identity provider... Computed.
- `client_secret` (String) OIDC client secret issued by the identity provider. Stored encrypted at rest... Computed.
- `scopes` (String) Space-separated list of OIDC scopes to request. Must include 'openid'... Computed.
- `email_claim_name` (String) Claim name in the ID token (or userinfo response) that contains the user's email address... Computed.
- `name_claim_name` (String) Claim name in the ID token (or userinfo response) that contains the user's display name... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_enabled` (Bool) Status page oidc is_enabled. Computed.
- `is_tested` (Bool) Status page oidc is_tested. Computed.
