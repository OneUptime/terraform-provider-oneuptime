---
page_title: "oneuptime_dashboard_private_user Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  People who can sign in to open a dashboard's public link when it asks its visitors to sign in.
---

# oneuptime_dashboard_private_user (Data Source)

People who can sign in to open a dashboard's public link when it asks its visitors to sign in.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one dashboard private user may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_dashboard_private_user" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_dashboard_private_user" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `dashboard_id` (String) ID of the dashboard this private user may sign in to. The ID of a `oneuptime_dashboard`.
- `has_password` (Boolean) Whether this private user has set a password.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_sso_user` (Boolean) Whether the dashboard's single sign-on added this user.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `email` (String) The address this private user signs in with.
- `last_active` (String) When this private user last signed in to the dashboard.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
