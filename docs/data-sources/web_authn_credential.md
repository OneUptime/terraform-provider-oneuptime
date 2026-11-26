---
page_title: "oneuptime_web_authn_credential Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  WebAuthn credentials for users (security keys)
---

# oneuptime_web_authn_credential (Data Source)

WebAuthn credentials for users (security keys)

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one web authn credential may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_web_authn_credential" "example" {
  name = "Example web authn credential"
}

# Or by id:
data "oneuptime_web_authn_credential" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_passkey` (Boolean) Whether this credential was registered for passwordless sign-in. Older credentials have no recorded purpose.
- `is_verified` (Boolean) Is this WebAuthn credential verified and validated.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Name of the WebAuthn credential.
- `user_id` (String) User ID who owns this WebAuthn credential. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `updated_at` (String) Date and Time when the object was updated.
