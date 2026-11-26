---
page_title: "oneuptime_dashboard_oidc Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Manage OpenID Connect (OIDC) sign-in for a dashboard that asks its visitors to sign in
---

# oneuptime_dashboard_oidc (Data Source)

Manage OpenID Connect (OIDC) sign-in for a dashboard that asks its visitors to sign in

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one dashboard oidc may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_dashboard_oidc" "example" {
  name = "Example dashboard oidc"
}

# Or by id:
data "oneuptime_dashboard_oidc" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `client_id` (String) OIDC client ID issued by the identity provider.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `dashboard_id` (String) ID of the dashboard this provider signs people in to. The ID of a `oneuptime_dashboard`.
- `description` (String) What this sign-in is for, such as which identity provider it uses.
- `discovery_url` (String) OIDC discovery URL (typically ends in /.well-known/openid-configuration). Used to discover authorization, token, JWKS and userinfo endpoints.
- `email_claim_name` (String) Claim name in the ID token (or userinfo response) that contains the user's email address.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) When on, visitors of a dashboard that asks them to sign in can sign in with this provider.
- `is_tested` (Boolean) Whether a sign-in with this provider has been tested. It does not turn the provider on.
- `issuer_url` (String) Expected OIDC issuer URL. Must match the 'iss' claim in the ID token returned by the identity provider.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Any friendly name of this object.
- `name_claim_name` (String) Claim name in the ID token (or userinfo response) that contains the user's display name.
- `scopes` (String) Space-separated list of OIDC scopes to request. Must include 'openid'.

### Read-Only

- `client_secret` (String, Sensitive) OIDC client secret issued by the identity provider. Stored encrypted at rest.
- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
