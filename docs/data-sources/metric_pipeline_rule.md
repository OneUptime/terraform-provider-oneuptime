---
page_title: "oneuptime_metric_pipeline_rule Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Rules applied at metric ingest time to filter, drop, rename, enrich, redact, or sample metric data points.
---

# oneuptime_metric_pipeline_rule (Data Source)

Rules applied at metric ingest time to filter, drop, rename, enrich, redact, or sample metric data points.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one metric pipeline rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_metric_pipeline_rule" "example" {
  name = "Example metric pipeline rule"
}

# Or by id:
data "oneuptime_metric_pipeline_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `add_attribute_key` (String) For AddAttribute / RemoveAttribute / RedactAttribute: the attribute key to act on.
- `add_attribute_value` (String) For AddAttribute: the attribute value to set.
- `created_by_user_id` (String) ID of the user who created this rule. The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of what this rule does.
- `filter_condition` (String) How to combine filters: 'All' requires every filter to match (AND), 'Any' requires at least one to match (OR).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule is active.
- `name` (String) Friendly name for this rule.
- `redact_replacement` (String) For RedactAttribute: the literal string to replace the value with. Defaults to [REDACTED].
- `rename_from_key` (String) For RenameMetric: the existing metric name. For RenameAttribute: the existing attribute key.
- `rename_to_key` (String) For RenameMetric: the new metric name. For RenameAttribute: the new attribute key.
- `rule_type` (String) One of: Filter, Drop, RenameMetric, RenameAttribute, AddAttribute, RemoveAttribute, RedactAttribute, Sample.
- `sample_percentage` (Number) For Sample: percentage of matched rows to keep (0-100). 100 keeps all.
- `service_id` (String) Optional service ID scoping this rule. Null means the rule is project-wide. The ID of a `oneuptime_service`.
- `sort_order` (Number) Where this rule is evaluated among the project's metric pipeline rules, lowest number first. A new rule is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `filters` (String) List of filters evaluated against each metric data point. An empty list matches every data point. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of the project this metric pipeline rule belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
