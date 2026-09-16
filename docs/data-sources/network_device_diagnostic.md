---
page_title: "oneuptime_network_device_diagnostic Data Source - oneuptime"
subcategory: "Other"
description: |-
  An on-demand ping or traceroute run against a Network Device from its probe. Transient: rows are deleted after two days.
---

# oneuptime_network_device_diagnostic (Data Source)

An on-demand ping or traceroute run against a Network Device from its probe. Transient: rows are deleted after two days. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_network_device_diagnostic" "by_name" {
  name = "example-network_device_diagnostic"
}

data "oneuptime_network_device_diagnostic" "by_id" {
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
- `network_device_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `probe_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `diagnostic_type` (String) What to run against the device: "Ping" (ICMP echo: reachability, round-trip time, jitter, packet loss) or "Traceroute" (the hop-by-hop path from the probe to the device)... Computed.
- `hostname` (String) The hostname or IP address the probe reaches, copied from the device when the diagnostic is created. Managed by the server... Computed.
- `status` (String) Where this diagnostic is in its run: "Pending" (waiting for the probe), "In Progress" (claimed by the probe), "Completed" (a result is stored) or "Failed" (the probe could not run it; see Status Message). Managed by the server and the probe... Computed.
- `status_message` (String) Why a diagnostic Failed, e.g. the device has no usable hostname or the probe does not support this diagnostic. Managed by the probe... Computed.
- `ping_result` (String) For a Ping diagnostic: whether the device answered, the failure cause when it did not, and the packet statistics (min/avg/max round-trip time, jitter, packet loss). Managed by the probe... Computed.
- `trace_route_result` (String) For a Traceroute diagnostic: the DNS lookup and the hop-by-hop path from the probe to the device, in the same shape the Network monitor records. Managed by the probe... Computed.
- `started_at` (String) A date time object.. Computed.
- `completed_at` (String) A date time object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
