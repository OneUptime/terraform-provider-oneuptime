---
page_title: "oneuptime_network_site_status_event Data Source - oneuptime"
subcategory: "Network"
description: |-
  History of the rolled-up health status of a Network Site (Operational to Offline for example)
---

# oneuptime_network_site_status_event (Data Source)

History of the rolled-up health status of a Network Site (Operational to Offline for example)

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one network site status event may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_network_site_status_event" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_network_site_status_event" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `monitor_status_id` (String) Relation to Monitor Status ID Resource in which this object belongs. The ID of a `oneuptime_monitor_status`.
- `site_id` (String) Relation to Network Site ID Resource in which this object belongs. The ID of a `oneuptime_network_site`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `ends_at` (String) When did this status change end?
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `starts_at` (String) When did this status change?
- `updated_at` (String) Date and Time when the object was updated.
