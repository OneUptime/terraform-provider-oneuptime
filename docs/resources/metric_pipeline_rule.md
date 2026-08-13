---
page_title: "oneuptime_metric_pipeline_rule Resource - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Rules applied at metric ingest time to filter, drop, rename, enrich, redact, or sample metric data points.
---

# oneuptime_metric_pipeline_rule (Resource)

Rules applied at metric ingest time to filter, drop, rename, enrich, redact, or sample metric data points.

## Example Usage

```terraform
resource "oneuptime_metric_pipeline_rule" "example" {
  name = jsonencode({
    "_type": "Name",
    "value": "John Doe"
  })
  rule_type = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name object.
- `rule_type` (String) One of: Filter, Drop, RenameMetric, RenameAttribute, AddAttribute, RemoveAttribute, RedactAttribute, Sample...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `service_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of what this rule does...
- `filter_condition` (String) How to combine filters: 'All' requires every filter to match (AND), 'Any' requires at least one to match (OR)...
- `filters` (String) List of filters evaluated against each metric data point. An empty list matches every data point...
- `rename_from_key` (String) For RenameMetric: the existing metric name. For RenameAttribute: the existing attribute key...
- `rename_to_key` (String) For RenameMetric: the new metric name. For RenameAttribute: the new attribute key...
- `add_attribute_key` (String) For AddAttribute / RemoveAttribute / RedactAttribute: the attribute key to act on...
- `add_attribute_value` (String) For AddAttribute: the attribute value to set...
- `redact_replacement` (String) For RedactAttribute: the literal string to replace the value with. Defaults to [REDACTED]...
- `sample_percentage` (Number) For Sample: percentage of matched rows to keep (0-100). 100 keeps all...
- `is_enabled` (Bool) Whether this rule is active...
- `sort_order` (Number) Evaluation order within its scope (service-level or project-level)...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_metric_pipeline_rule.example <id>
```
