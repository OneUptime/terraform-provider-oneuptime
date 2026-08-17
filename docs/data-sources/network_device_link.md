---
page_title: "oneuptime_network_device_link Data Source - oneuptime"
subcategory: "Other"
description: |-
  Operator-declared links between two Network Devices, for cables LLDP and CDP cannot see: a device with discovery disabled, a device that does not speak either protocol, or one monitored by ping alone. Drawn on the topology map alongside discovered links, and merged with a discovered link between the same pair rather than duplicating it.
---

# oneuptime_network_device_link (Data Source)

Operator-declared links between two Network Devices, for cables LLDP and CDP cannot see: a device with discovery disabled, a device that does not speak either protocol, or one monitored by ping alone. Drawn on the topology map alongside discovered links, and merged with a discovered link between the same pair rather than duplicating it. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_network_device_link" "by_name" {
  name = "example-network_device_link"
}

data "oneuptime_network_device_link" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `from_device_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `to_device_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `parent_device_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `from_port_name` (String) Port on the starting device, as free text. Nothing resolves it to an interface row — a hand-drawn link usually exists precisely because the port is not discoverable... Computed.
- `to_port_name` (String) Port on the ending device, as free text... Computed.
- `monitor_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
