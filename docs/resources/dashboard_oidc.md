---
page_title: "oneuptime_dashboard_oidc Resource - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Manage OpenID Connect (OIDC) sign-in for a dashboard that asks its visitors to sign in
---

# oneuptime_dashboard_oidc (Resource)

Manage OpenID Connect (OIDC) sign-in for a dashboard that asks its visitors to sign in

## Example Usage

```terraform
resource "oneuptime_dashboard_oidc" "example" {
  dashboard_id     = oneuptime_dashboard.example.id
  name             = "Example dashboard oidc"
  description      = "Managed by Terraform"
  discovery_url    = "https://www.example.com/path/to/resource?param=value"
  issuer_url       = "https://accounts.google.com"
  client_id        = "1234567890-abcdefgh.apps.googleusercontent.com"
  client_secret    = "This is an example of longer text content that might be stored in this field."
  scopes           = "openid email profile"
  email_claim_name = "email"
}
```

## Schema

### Required

- `client_id` (String) OIDC client ID issued by the identity provider.
- `client_secret` (String, Sensitive) OIDC client secret issued by the identity provider. Stored encrypted at rest.
- `dashboard_id` (String) ID of the dashboard this provider signs people in to. The ID of a `oneuptime_dashboard`.
- `description` (String) What this sign-in is for, such as which identity provider it uses.
- `discovery_url` (String) OIDC discovery URL (typically ends in /.well-known/openid-configuration). Used to discover authorization, token, JWKS and userinfo endpoints.
- `email_claim_name` (String) Claim name in the ID token (or userinfo response) that contains the user's email address.
- `issuer_url` (String) Expected OIDC issuer URL. Must match the 'iss' claim in the ID token returned by the identity provider.
- `name` (String) Any friendly name of this object.
- `scopes` (String) Space-separated list of OIDC scopes to request. Must include 'openid'.

### Optional

- `is_enabled` (Boolean) When on, visitors of a dashboard that asks them to sign in can sign in with this provider. Defaults to `false`.
- `is_tested` (Boolean) Whether a sign-in with this provider has been tested. It does not turn the provider on. Defaults to `false`.
- `name_claim_name` (String) Claim name in the ID token (or userinfo response) that contains the user's display name.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing dashboard oidc by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_dashboard_oidc.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_dashboard_oidc.example <id>
```
