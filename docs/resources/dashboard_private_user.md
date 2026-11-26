---
page_title: "oneuptime_dashboard_private_user Resource - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  People who can sign in to open a dashboard's public link when it asks its visitors to sign in.
---

# oneuptime_dashboard_private_user (Resource)

People who can sign in to open a dashboard's public link when it asks its visitors to sign in.

## Example Usage

```terraform
resource "oneuptime_dashboard_private_user" "example" {
  dashboard_id = oneuptime_dashboard.example.id
  email        = "user@example.com"
}
```

## Schema

### Required

- `dashboard_id` (String) ID of the dashboard this private user may sign in to. The ID of a `oneuptime_dashboard`.
- `email` (String) The address this private user signs in with.

### Optional

- `password` (String, Sensitive) The private user's password, stored as a salted hash and never shown.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `has_password` (Boolean) Whether this private user has set a password.
- `id` (String) Unique identifier for the resource.
- `is_sso_user` (Boolean) Whether the dashboard's single sign-on added this user.
- `last_active` (String) When this private user last signed in to the dashboard.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing dashboard private user by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_dashboard_private_user.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_dashboard_private_user.example <id>
```
