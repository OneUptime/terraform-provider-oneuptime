---
page_title: "oneuptime_kubernetes_cost_allocation Data Source - oneuptime"
subcategory: "Other"
description: |-
  API endpoints for Kubernetes Cost Allocation
---

# oneuptime_kubernetes_cost_allocation (Data Source)

API endpoints for Kubernetes Cost Allocation

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one kubernetes cost allocation may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_kubernetes_cost_allocation" "example" {
  kubernetes_cluster_id = "example-kubernetes-cluster-id"
}

# Or by id:
data "oneuptime_kubernetes_cost_allocation" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `cluster_name` (String) Cluster Name.
- `container_name` (String) Container Name.
- `controller_kind` (String) Controller Kind.
- `controller_name` (String) Controller Name.
- `cpu_core_hours` (Number) CPU Core Hours.
- `cpu_core_limit_average` (Number) CPU Core Limit Average.
- `cpu_core_request_average` (Number) CPU Core Request Average.
- `cpu_core_usage_average` (Number) CPU Core Usage Average.
- `cpu_cost` (Number) CPU Cost.
- `cpu_efficiency` (Number) CPU Efficiency.
- `currency` (String) Currency.
- `external_cost` (Number) External Cost.
- `gpu_cost` (Number) GPU Cost.
- `gpu_hours` (Number) GPU Hours.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `k8s_cluster_entity_key` (String) Kubernetes Cluster Entity Key.
- `kubernetes_cluster_id` (String) Kubernetes Cluster ID.
- `labels` (String) Labels.
- `load_balancer_cost` (Number) Load Balancer Cost.
- `namespace` (String) Namespace.
- `network_cost` (Number) Network Cost.
- `node_name` (String) Node Name.
- `pod_name` (String) Pod Name.
- `provider_id` (String) Provider ID.
- `pv_byte_hours` (Number) PV Byte Hours.
- `pv_cost` (Number) PV Cost.
- `ram_byte_hours` (Number) RAM Byte Hours.
- `ram_bytes_limit_average` (Number) RAM Bytes Limit Average.
- `ram_bytes_request_average` (Number) RAM Bytes Request Average.
- `ram_bytes_usage_average` (Number) RAM Bytes Usage Average.
- `ram_bytes_usage_max` (Number) RAM Bytes Usage Max.
- `ram_cost` (Number) RAM Cost.
- `ram_efficiency` (Number) RAM Efficiency.
- `shared_cost` (Number) Shared Cost.
- `shipment_chunk` (Number) Shipment Chunk.
- `shipment_id` (String) Shipment ID.
- `total_cost` (Number) Total Cost.
- `total_efficiency` (Number) Total Efficiency.
- `window_end` (String) Window End.
- `window_start` (String) Window Start.

### Read-Only

- `label_keys` (Set of String) Label Keys.
- `project_id` (String) Project ID.
