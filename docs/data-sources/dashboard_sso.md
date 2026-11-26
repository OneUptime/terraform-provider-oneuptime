---
page_title: "oneuptime_dashboard_sso Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Configure SAML single sign-on for a dashboard that asks its visitors to sign in
---

# oneuptime_dashboard_sso (Data Source)

Configure SAML single sign-on for a dashboard that asks its visitors to sign in

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one dashboard sso may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_dashboard_sso" "example" {
  name = "Example dashboard sso"
}

# Or by id:
data "oneuptime_dashboard_sso" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `dashboard_id` (String) ID of the dashboard this provider signs people in to. The ID of a `oneuptime_dashboard`.
- `description` (String) What this sign-in is for, such as which identity provider it uses.
- `digest_method` (String) The digest algorithm your identity provider's SAML signatures use.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) When on, visitors of a dashboard that asks them to sign in can sign in with this provider.
- `is_tested` (Boolean) Whether a sign-in with this provider has been tested. It does not turn the provider on.
- `issuer_url` (String) Your identity provider's issuer (entity ID), as its SAML responses name it.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Any friendly name of this object.
- `public_certificate` (String) Your identity provider's signing certificate, in PEM, to check its SAML responses with.
- `sign_on_url` (String) Your identity provider's single sign-on URL, where visitors are sent to sign in.
- `signature_method` (String) The algorithm your identity provider signs SAML responses with.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
