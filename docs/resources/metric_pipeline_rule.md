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
  name        = "Example metric pipeline rule"
  rule_type   = "Example short text"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Friendly name for this rule.
- `rule_type` (String) One of: Filter, Drop, RenameMetric, RenameAttribute, AddAttribute, RemoveAttribute, RedactAttribute, Sample.

### Optional

- `add_attribute_key` (String) For AddAttribute / RemoveAttribute / RedactAttribute: the attribute key to act on.
- `add_attribute_value` (String) For AddAttribute: the attribute value to set.
- `description` (String) Description of what this rule does.
- `filter_condition` (String) How to combine filters: 'All' requires every filter to match (AND), 'Any' requires at least one to match (OR). Defaults to `All`.
- `filters` (String) List of filters evaluated against each metric data point. An empty list matches every data point. A JSON value: write it with `jsonencode()`.
- `is_enabled` (Boolean) Whether this rule is active. Defaults to `true`.
- `redact_replacement` (String) For RedactAttribute: the literal string to replace the value with. Defaults to [REDACTED].
- `rename_from_key` (String) For RenameMetric: the existing metric name. For RenameAttribute: the existing attribute key.
- `rename_to_key` (String) For RenameMetric: the new metric name. For RenameAttribute: the new attribute key.
- `sample_percentage` (Number) For Sample: percentage of matched rows to keep (0-100). 100 keeps all. Defaults to `100`.
- `service_id` (String) Optional service ID scoping this rule. Null means the rule is project-wide. The ID of a `oneuptime_service`.
- `sort_order` (Number) Where this rule is evaluated among the project's metric pipeline rules, lowest number first. A new rule is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this rule. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this metric pipeline rule belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing metric pipeline rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_metric_pipeline_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_metric_pipeline_rule.example <id>
```
