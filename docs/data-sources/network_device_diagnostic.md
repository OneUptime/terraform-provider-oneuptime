---
page_title: "oneuptime_network_device_diagnostic Data Source - oneuptime"
subcategory: "Network"
description: |-
  An on-demand ping, traceroute or SNMP walk run against a Network Device from its probe. Transient: rows are deleted after two days.
---

# oneuptime_network_device_diagnostic (Data Source)

An on-demand ping, traceroute or SNMP walk run against a Network Device from its probe. Transient: rows are deleted after two days.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one network device diagnostic may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_network_device_diagnostic" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_network_device_diagnostic" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `diagnostic_type` (String) What to run against the device: "Ping" (ICMP echo: reachability, round-trip time, jitter, packet loss), "Traceroute" (the hop-by-hop path from the probe to the device) or "SNMP Walk" (every OID the device answers under a starting OID, read with its SNMP credentials).
- `hostname` (String) The hostname or IP address the probe reaches, copied from the device when the diagnostic is created. Managed by the server.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `network_device_id` (String) ID of the Network Device this diagnostic runs against. The ID of a `oneuptime_network_device`.
- `probe_id` (String) ID of the Probe that runs this diagnostic. Defaults to the device's assigned probe. The ID of a `oneuptime_probe`.
- `snmp_walk_after_oid` (String) For an SNMP Walk that continues an earlier one: the last OID the earlier walk read. The walk reads what follows it under the start OID.
- `snmp_walk_root_oid` (String) For an SNMP Walk: the OID whose subtree is walked, such as 1.3.6.1.2.1.2 for the interfaces. Defaults to 1.3.6.1, everything the device exposes.
- `status` (String) Where this diagnostic is in its run: "Pending" (waiting for the probe), "In Progress" (claimed by the probe), "Completed" (a result is stored) or "Failed" (the probe could not run it; see Status Message). Managed by the server and the probe.
- `status_message` (String) Why a diagnostic Failed, e.g. the device has no usable hostname or the probe does not support this diagnostic. Managed by the probe.

### Read-Only

- `completed_at` (String) When the probe reported a result (or a failure) for this diagnostic. Managed by the server.
- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `ping_result` (String) For a Ping diagnostic: whether the device answered, the failure cause when it did not, and the packet statistics (min/avg/max round-trip time, jitter, packet loss). Managed by the probe. A JSON value: read it with `jsondecode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `snmp_walk_result` (String) For an SNMP Walk: each OID the device answered with its type and value, why the walk stopped, where to continue from, and, when the device refused or did not answer, why. Managed by the probe. A JSON value: read it with `jsondecode()`.
- `started_at` (String) When the probe claimed this diagnostic. Managed by the server.
- `trace_route_result` (String) For a Traceroute diagnostic: the DNS lookup and the hop-by-hop path from the probe to the device, in the same shape the Network monitor records. Managed by the probe. A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.
