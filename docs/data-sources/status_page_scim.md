---
page_title: "oneuptime_status_page_scim Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage SCIM auto-provisioning for your status page
---

# oneuptime_status_page_scim (Data Source)

Manage SCIM auto-provisioning for your status page

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one status page scim may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_status_page_scim" "example" {
  name = "Example status page scim"
}

# Or by id:
data "oneuptime_status_page_scim" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `auto_deprovision_users` (Boolean) Automatically remove status page users when they are removed via SCIM.
- `auto_provision_users` (Boolean) Automatically create status page users when they are added via SCIM.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description to help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Any friendly name for this SCIM configuration.
- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `oauth_client_secret_expires_at` (String) When the connection's OAuth 2.0 client secret stops working. No access token is issued with an expired secret.
- `oauth_previous_client_secret_expires_at` (String) When the previous OAuth 2.0 client secret stops working, unless it is deleted before.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
