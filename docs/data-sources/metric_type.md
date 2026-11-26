---
page_title: "oneuptime_metric_type Data Source - oneuptime"
subcategory: "Logs, Traces & Metrics"
description: |-
  List of all the metrics ingested with OpenTelemetry
---

# oneuptime_metric_type (Data Source)

List of all the metrics ingested with OpenTelemetry

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one metric type may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_metric_type" "example" {
  name = "Example metric type"
}

# Or by id:
data "oneuptime_metric_type" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `aggregation_temporality` (String) OpenTelemetry aggregation temporality of this metric (Delta or Cumulative), as reported at ingest. Null when unknown.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Metric description.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_monotonic` (Boolean) Whether this metric is a monotonic counter (only ever increases), as reported by OpenTelemetry at ingest. Null when the instrument type does not carry monotonicity (e.g. gauges).
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Any friendly name of this object.
- `unit` (String) Metric description.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `services` (Set of String) List of services this metric is related to. IDs of `oneuptime_service` resources.
- `updated_at` (String) Date and Time when the object was updated.
