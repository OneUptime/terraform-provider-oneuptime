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
  network_device_id = oneuptime_network_device.example.id
  diagnostic_type   = "Example short text"
}
```

## Schema

### Required

- `diagnostic_type` (String) What to run against the device: "Ping" (ICMP echo: reachability, round-trip time, jitter, packet loss) or "Traceroute" (the hop-by-hop path from the probe to the device).
- `network_device_id` (String) ID of the Network Device this diagnostic runs against. The ID of a `oneuptime_network_device`.

### Optional

- `probe_id` (String) ID of the Probe that runs this diagnostic. Defaults to the device's assigned probe. The ID of a `oneuptime_probe`.

### Read-Only

- `completed_at` (String) When the probe reported a result (or a failure) for this diagnostic. Managed by the server.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `hostname` (String) The hostname or IP address the probe reaches, copied from the device when the diagnostic is created. Managed by the server.
- `id` (String) Unique identifier for the resource.
- `ping_result` (String) For a Ping diagnostic: whether the device answered, the failure cause when it did not, and the packet statistics (min/avg/max round-trip time, jitter, packet loss). Managed by the probe. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `started_at` (String) When the probe claimed this diagnostic. Managed by the server.
- `status` (String) Where this diagnostic is in its run: "Pending" (waiting for the probe), "In Progress" (claimed by the probe), "Completed" (a result is stored) or "Failed" (the probe could not run it; see Status Message). Managed by the server and the probe.
- `status_message` (String) Why a diagnostic Failed, e.g. the device has no usable hostname or the probe does not support this diagnostic. Managed by the probe.
- `trace_route_result` (String) For a Traceroute diagnostic: the DNS lookup and the hop-by-hop path from the probe to the device, in the same shape the Network monitor records. Managed by the probe. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing network device diagnostic by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_network_device_diagnostic.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_network_device_diagnostic.example <id>
```
