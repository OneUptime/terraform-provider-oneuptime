---
page_title: "oneuptime_network_flow Data Source - oneuptime"
subcategory: "Other"
description: |-
  API endpoints for Network Flow
---

# oneuptime_network_flow (Data Source)

API endpoints for Network Flow

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one network flow may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_network_flow" "example" {
  network_device_id = "example-network-device-id"
}

# Or by id:
data "oneuptime_network_flow" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `dst_ip` (String) Destination IP.
- `dst_port` (Number) Destination Port.
- `exporter_ip` (String) Exporter IP.
- `flow_end_at` (String) Flow End.
- `flow_start_at` (String) Flow Start.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `ingested_at` (String) Ingested At.
- `input_interface_index` (Number) Input Interface Index.
- `network_device_id` (String) Network Device ID.
- `octets` (String) Octets.
- `output_interface_index` (Number) Output Interface Index.
- `packets` (String) Packets.
- `protocol` (Number) Protocol.
- `src_ip` (String) Source IP.
- `src_port` (Number) Source Port.

### Read-Only

- `project_id` (String) Project ID.
