---
page_title: "oneuptime_audit_log Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  API endpoints for Audit Log
---

# oneuptime_audit_log (Data Source)

API endpoints for Audit Log

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one audit log may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_audit_log" "example" {
  resource_type = "example-resource-type"
}

# Or by id:
data "oneuptime_audit_log" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `action` (String) Action.
- `api_key_id` (String) API Key ID.
- `api_key_name` (String) API Key Name.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `mcp_client_name` (String) MCP Client Name.
- `mcp_o_auth_grant_id` (String) MCP Client Authorization ID.
- `resource_id` (String) Resource ID.
- `resource_name` (String) Resource Name.
- `resource_type` (String) Resource Type.
- `root_resource_id` (String) Root Resource ID.
- `root_resource_type` (String) Root Resource Type.
- `user_email` (String) User Email.
- `user_id` (String) User ID.
- `user_name` (String) User Name.
- `user_type` (String) User Type.
- `workflow_id` (String) Workflow ID.
- `workflow_name` (String) Workflow Name.

### Read-Only

- `changes` (Set of String) Changes.
- `project_id` (String) Project ID.
