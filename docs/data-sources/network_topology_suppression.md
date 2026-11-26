---
page_title: "oneuptime_network_topology_suppression Data Source - oneuptime"
subcategory: "Network"
description: |-
  Nodes hidden from the Network Topology map for the whole project. Display only — the device and its monitoring are untouched.
---

# oneuptime_network_topology_suppression (Data Source)

Nodes hidden from the Network Topology map for the whole project. Display only — the device and its monitoring are untouched.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one network topology suppression may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_network_topology_suppression" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_network_topology_suppression" "by_id" {
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
- `node_key` (String) The topology node id to hide. A device id for a managed device, 'unmanaged:<name>' for a discovery-protocol peer, or 'endpoint:<id>' for a discovered endpoint. Free text rather than a foreign key because two of the three are synthesised by the topology builder and have no row of their own.
- `node_name` (String) What the node was called when it was hidden, so the hidden list is readable without rebuilding the graph.
- `reason` (String) Why this node was hidden — the note the next person needs to decide whether it should stay hidden.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
