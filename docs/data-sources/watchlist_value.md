---
page_title: "oneuptime_watchlist_value Data Source - oneuptime"
subcategory: "Security"
description: |-
  The values on a watchlist the team keeps, each with an optional note, matched by the detection rules that read the list.
---

# oneuptime_watchlist_value (Data Source)

The values on a watchlist the team keeps, each with an optional note, matched by the detection rules that read the list.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one watchlist value may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_watchlist_value" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_watchlist_value" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the person who added the value. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `note` (String) Why the value is on the list, for whoever reads it later.
- `security_watchlist_id` (String) ID of the watchlist this value is on. The ID of a `oneuptime_watchlist`.
- `value` (String) The value, at most 256 characters, matched whatever the case. * stands for any run of characters (svc-* is every service account), and an IP address list takes ranges such as 10.20.0.0/16.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `expires_at` (String) When the value leaves the list: from then on no rule, monitor, hunt or search reads it, and OneUptime removes it. Empty: it stays until someone removes it.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of the project this value belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
