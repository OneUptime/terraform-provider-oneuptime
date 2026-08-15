---
page_title: "oneuptime_network_device_link Resource - oneuptime"
subcategory: "Other"
description: |-
  Operator-declared links between two Network Devices, for cables LLDP and CDP cannot see: a device with discovery disabled, a device that does not speak either protocol, or one monitored by ping alone. Drawn on the topology map alongside discovered links, and merged with a discovered link between the same pair rather than duplicating it.
---

# oneuptime_network_device_link (Resource)

Operator-declared links between two Network Devices, for cables LLDP and CDP cannot see: a device with discovery disabled, a device that does not speak either protocol, or one monitored by ping alone. Drawn on the topology map alongside discovered links, and merged with a discovered link between the same pair rather than duplicating it.

## Example Usage

```terraform
resource "oneuptime_network_device_link" "example" {
  from_device_id = "123e4567-e89b-12d3-a456-426614174000"
  to_device_id = "123e4567-e89b-12d3-a456-426614174000"
  name = "Example short text"
}
```

## Schema

### Required

- `from_device_id` (String) A unique identifier for an object, represented as a UUID..
- `to_device_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `name` (String) Friendly name for this link..
- `from_port_name` (String) Port on the starting device, as free text. Nothing resolves it to an interface row — a hand-drawn link usually exists precisely because the port is not discoverable...
- `to_port_name` (String) Port on the ending device, as free text...
- `monitor_id` (String) A unique identifier for an object, represented as a UUID..
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
terraform import oneuptime_network_device_link.example <id>
```
