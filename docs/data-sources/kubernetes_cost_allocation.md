---
page_title: "oneuptime_kubernetes_cost_allocation Data Source - oneuptime"
subcategory: "Other"
description: |-
  API endpoints for Kubernetes Cost Allocation
---

# oneuptime_kubernetes_cost_allocation (Data Source)

API endpoints for Kubernetes Cost Allocation Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_kubernetes_cost_allocation" "by_name" {
  name = "example-kubernetes_cost_allocation"
}

data "oneuptime_kubernetes_cost_allocation" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `project_id` (String) Project ID. Computed.
- `kubernetes_cluster_id` (String) Kubernetes Cluster ID. Computed.
- `cluster_name` (String) Cluster Name. Computed.
- `k8s_cluster_entity_key` (String) Kubernetes Cluster Entity Key. Computed.
- `window_start` (String) Window Start. Computed.
- `window_end` (String) Window End. Computed.
- `namespace` (String) Namespace. Computed.
- `controller_kind` (String) Controller Kind. Computed.
- `controller_name` (String) Controller Name. Computed.
- `pod_name` (String) Pod Name. Computed.
- `container_name` (String) Container Name. Computed.
- `node_name` (String) Node Name. Computed.
- `provider_id` (String) Provider ID. Computed.
- `labels` (String) Labels. Computed.
- `label_keys` (Set) Label Keys. Computed.
- `cpu_core_hours` (Number) CPU Core Hours. Computed.
- `cpu_core_request_average` (Number) CPU Core Request Average. Computed.
- `cpu_core_usage_average` (Number) CPU Core Usage Average. Computed.
- `cpu_core_limit_average` (Number) CPU Core Limit Average. Computed.
- `cpu_cost` (Number) CPU Cost. Computed.
- `gpu_hours` (Number) GPU Hours. Computed.
- `gpu_cost` (Number) GPU Cost. Computed.
- `ram_byte_hours` (Number) RAM Byte Hours. Computed.
- `ram_bytes_request_average` (Number) RAM Bytes Request Average. Computed.
- `ram_bytes_usage_average` (Number) RAM Bytes Usage Average. Computed.
- `ram_bytes_limit_average` (Number) RAM Bytes Limit Average. Computed.
- `ram_bytes_usage_max` (Number) RAM Bytes Usage Max. Computed.
- `ram_cost` (Number) RAM Cost. Computed.
- `pv_byte_hours` (Number) PV Byte Hours. Computed.
- `pv_cost` (Number) PV Cost. Computed.
- `network_cost` (Number) Network Cost. Computed.
- `load_balancer_cost` (Number) Load Balancer Cost. Computed.
- `shared_cost` (Number) Shared Cost. Computed.
- `external_cost` (Number) External Cost. Computed.
- `total_cost` (Number) Total Cost. Computed.
- `cpu_efficiency` (Number) CPU Efficiency. Computed.
- `ram_efficiency` (Number) RAM Efficiency. Computed.
- `total_efficiency` (Number) Total Efficiency. Computed.
- `currency` (String) Currency. Computed.
- `shipment_id` (String) Shipment ID. Computed.
- `shipment_chunk` (Number) Shipment Chunk. Computed.
