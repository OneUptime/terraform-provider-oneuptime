---
page_title: "oneuptime_network_flow Data Source - oneuptime"
subcategory: "Other"
description: |-
  API endpoints for Network Flow
---

# oneuptime_network_flow (Data Source)

API endpoints for Network Flow Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_network_flow" "by_name" {
  name = "example-network_flow"
}

data "oneuptime_network_flow" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `project_id` (String) Project ID. Computed.
- `network_device_id` (String) Network Device ID. Computed.
- `exporter_ip` (String) Exporter IP. Computed.
- `src_ip` (String) Source IP. Computed.
- `dst_ip` (String) Destination IP. Computed.
- `src_port` (Number) Source Port. Computed.
- `dst_port` (Number) Destination Port. Computed.
- `protocol` (Number) Protocol. Computed.
- `input_interface_index` (Number) Input Interface Index. Computed.
- `output_interface_index` (Number) Output Interface Index. Computed.
- `octets` (String) Octets. Computed.
- `packets` (String) Packets. Computed.
- `flow_start_at` (String) Flow Start. Computed.
- `flow_end_at` (String) Flow End. Computed.
- `ingested_at` (String) Ingested At. Computed.
