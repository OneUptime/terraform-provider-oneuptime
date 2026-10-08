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
  status_page_id   = oneuptime_status_page.example.id
  name             = "Example status page oidc"
  description      = "Managed by Terraform"
  discovery_url    = "https://www.example.com/path/to/resource?param=value"
  issuer_url       = "This is an example of very long text content that might be stored in this field. It can contain a lot of information, such as detailed descriptions, comments, or any other lengthy text data that needs to be stored in the database."
  client_id        = "Example short text"
  client_secret    = "This is an example of longer text content that might be stored in this field."
  scopes           = "Example short text"
  email_claim_name = "Example short text"
}
```

## Schema

### Required

- `client_id` (String) OIDC client ID issued by the identity provider.
- `client_secret` (String) OIDC client secret issued by the identity provider. Stored encrypted at rest.
- `description` (String)
- `discovery_url` (String) OIDC discovery URL (typically ends in /.well-known/openid-configuration). Used to discover authorization, token, JWKS and userinfo endpoints.
- `email_claim_name` (String) Claim name in the ID token (or userinfo response) that contains the user's email address.
- `issuer_url` (String) Expected OIDC issuer URL. Must match the 'iss' claim in the ID token returned by the identity provider.
- `name` (String) Any friendly name of this object.
- `scopes` (String) Space-separated list of OIDC scopes to request. Must include 'openid'.
- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.

### Optional

- `is_enabled` (Boolean) Defaults to `false`.
- `is_tested` (Boolean) Defaults to `false`.
- `name_claim_name` (String) Claim name in the ID token (or userinfo response) that contains the user's display name.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing status page oidc by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_status_page_oidc.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_status_page_oidc.example <id>
```
