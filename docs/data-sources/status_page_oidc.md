---
page_title: "oneuptime_status_page_oidc Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage OpenID Connect (OIDC) authentication for your status page
---

# oneuptime_status_page_oidc (Data Source)

Manage OpenID Connect (OIDC) authentication for your status page

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one status page oidc may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_status_page_oidc" "example" {
  name = "Example status page oidc"
}

# Or by id:
data "oneuptime_status_page_oidc" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `client_id` (String) OIDC client ID issued by the identity provider.
- `client_secret` (String) OIDC client secret issued by the identity provider. Stored encrypted at rest.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Permissions - Create: [Project Owner, Project Admin, Create Status Page OIDC], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Status Page OIDC], Update: [Project Owner, Project Admin, Edit Status Page OIDC]
- `discovery_url` (String) OIDC discovery URL (typically ends in /.well-known/openid-configuration). Used to discover authorization, token, JWKS and userinfo endpoints.
- `email_claim_name` (String) Claim name in the ID token (or userinfo response) that contains the user's email address.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Permissions - Create: [Project Owner, Project Admin, Create Status Page OIDC], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Status Page OIDC], Update: [Project Owner, Project Admin, Edit Status Page OIDC]
- `is_tested` (Boolean) Permissions - Create: [Project Owner, Project Admin, Create Status Page OIDC], Read: [Project Owner, Project Admin, Read Status Page OIDC], Update: [No access - you don't have permission for this operation]
- `issuer_url` (String) Expected OIDC issuer URL. Must match the 'iss' claim in the ID token returned by the identity provider.
- `name` (String) Any friendly name of this object.
- `name_claim_name` (String) Claim name in the ID token (or userinfo response) that contains the user's display name.
- `scopes` (String) Space-separated list of OIDC scopes to request. Must include 'openid'.
- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
