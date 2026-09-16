---
page_title: "oneuptime_network_device_diagnostic Resource - oneuptime"
subcategory: "Other"
description: |-
  An on-demand ping or traceroute run against a Network Device from its probe. Transient: rows are deleted after two days.
---

# oneuptime_network_device_diagnostic (Resource)

An on-demand ping or traceroute run against a Network Device from its probe. Transient: rows are deleted after two days.

## Example Usage

```terraform
resource "oneuptime_network_device_diagnostic" "example" {
  network_device_id = "123e4567-e89b-12d3-a456-426614174000"
  diagnostic_type = "Example short text"
}
```

## Schema

### Required

- `network_device_id` (String) A unique identifier for an object, represented as a UUID..
- `diagnostic_type` (String) What to run against the device: "Ping" (ICMP echo: reachability, round-trip time, jitter, packet loss) or "Traceroute" (the hop-by-hop path from the probe to the device)...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `probe_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `hostname` (String) The hostname or IP address the probe reaches, copied from the device when the diagnostic is created. Managed by the server...
- `status` (String) Where this diagnostic is in its run: "Pending" (waiting for the probe), "In Progress" (claimed by the probe), "Completed" (a result is stored) or "Failed" (the probe could not run it; see Status Message). Managed by the server and the probe...
- `status_message` (String) Why a diagnostic Failed, e.g. the device has no usable hostname or the probe does not support this diagnostic. Managed by the probe...
- `ping_result` (String) For a Ping diagnostic: whether the device answered, the failure cause when it did not, and the packet statistics (min/avg/max round-trip time, jitter, packet loss). Managed by the probe...
- `trace_route_result` (String) For a Traceroute diagnostic: the DNS lookup and the hop-by-hop path from the probe to the device, in the same shape the Network monitor records. Managed by the probe...
- `started_at` (String) A date time object..
- `completed_at` (String) A date time object..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_network_device_diagnostic.example <id>
```
