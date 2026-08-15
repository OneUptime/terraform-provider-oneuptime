---
page_title: "oneuptime_network_topology_suppression Resource - oneuptime"
subcategory: "Other"
description: |-
  Nodes hidden from the Network Topology map for the whole project. Display only — the device and its monitoring are untouched.
---

# oneuptime_network_topology_suppression (Resource)

Nodes hidden from the Network Topology map for the whole project. Display only — the device and its monitoring are untouched.

## Example Usage

```terraform
resource "oneuptime_network_topology_suppression" "example" {
  node_key = "Example short text"
}
```

## Schema

### Required

- `node_key` (String) The topology node id to hide. A device id for a managed device, 'unmanaged:<name>' for a discovery-protocol peer, or 'endpoint:<id>' for a discovered endpoint. Free text rather than a foreign key because two of the three are synthesised by the topology builder and have no row of their own...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `node_name` (String) What the node was called when it was hidden, so the hidden list is readable without rebuilding the graph...
- `reason` (String) Why this node was hidden — the note the next person needs to decide whether it should stay hidden...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_network_topology_suppression.example <id>
```
