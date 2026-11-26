---
page_title: "oneuptime_status_page_scim_log Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Logs of all SCIM provisioning operations for status pages.
---

# oneuptime_status_page_scim_log (Data Source)

Logs of all SCIM provisioning operations for status pages.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one status page scim log may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_status_page_scim_log" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_status_page_scim_log" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `http_method` (String) HTTP method used (GET, POST, PUT, PATCH, DELETE).
- `http_status_code` (Number) Response HTTP status code.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `log_body` (String) Detailed JSON with request/response data.
- `operation_type` (String) Type of SCIM operation (e.g., CreateUser, UpdateUser, DeleteUser, ListUsers, GetUser, BulkOperation).
- `request_path` (String) The SCIM endpoint path.
- `status` (String) Status of the SCIM operation.
- `status_message` (String) Short error or status description.
- `status_page_id` (String) ID of the Status Page. The ID of a `oneuptime_status_page`.
- `status_page_scim_id` (String) ID of your Status Page SCIM configuration. The ID of a `oneuptime_status_page_scim`.

### Read-Only

- `affected_user_email` (String) Email of the user affected by this operation.
- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
