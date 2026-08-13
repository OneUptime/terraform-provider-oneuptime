---
page_title: "oneuptime_metric_pipeline_rule Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Rules applied at metric ingest time to filter, drop, rename, enrich, redact, or sample metric data points.
---

# oneuptime_metric_pipeline_rule (Data Source)

Rules applied at metric ingest time to filter, drop, rename, enrich, redact, or sample metric data points. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_metric_pipeline_rule" "by_name" {
  name = "example-metric_pipeline_rule"
}

data "oneuptime_metric_pipeline_rule" "by_id" {
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
- `service_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description` (String) Description of what this rule does... Computed.
- `rule_type` (String) One of: Filter, Drop, RenameMetric, RenameAttribute, AddAttribute, RemoveAttribute, RedactAttribute, Sample... Computed.
- `filter_condition` (String) How to combine filters: 'All' requires every filter to match (AND), 'Any' requires at least one to match (OR)... Computed.
- `filters` (String) List of filters evaluated against each metric data point. An empty list matches every data point... Computed.
- `rename_from_key` (String) For RenameMetric: the existing metric name. For RenameAttribute: the existing attribute key... Computed.
- `rename_to_key` (String) For RenameMetric: the new metric name. For RenameAttribute: the new attribute key... Computed.
- `add_attribute_key` (String) For AddAttribute / RemoveAttribute / RedactAttribute: the attribute key to act on... Computed.
- `add_attribute_value` (String) For AddAttribute: the attribute value to set... Computed.
- `redact_replacement` (String) For RedactAttribute: the literal string to replace the value with. Defaults to [REDACTED]... Computed.
- `sample_percentage` (Number) For Sample: percentage of matched rows to keep (0-100). 100 keeps all... Computed.
- `is_enabled` (Bool) Whether this rule is active... Computed.
- `sort_order` (Number) Evaluation order within its scope (service-level or project-level)... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
