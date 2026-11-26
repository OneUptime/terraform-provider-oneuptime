---
page_title: "oneuptime_metric_pipeline_rule Resource - oneuptime"
subcategory: "Logs, Traces & Metrics"
description: |-
  Rules applied at metric ingest time to filter, drop, rename, enrich, redact, or sample metric data points.
---

# oneuptime_metric_pipeline_rule (Resource)

Rules applied at metric ingest time to filter, drop, rename, enrich, redact, or sample metric data points.

## Example Usage

```terraform
resource "oneuptime_metric_pipeline_rule" "example" {
  name        = "Example metric pipeline rule"
  rule_type   = "Filter"
  description = "Managed by Terraform"
}
```

<a id="order"></a>
## Order

`sort_order` is this metric pipeline rule's place among the metric pipeline rules of its project, lowest first. Setting it places this metric pipeline rule and moves no other, so a configuration that numbers every metric pipeline rule of a project is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.metric_pipeline_rules`:

```terraform
sort_order = index(local.metric_pipeline_rules, each.key) + 1
```

- Give each metric pipeline rule a number of its own: two metric pipeline rules at one number are shown the older first.
- A metric pipeline rule created without `sort_order` goes after the last one of its project.
- Deleting a metric pipeline rule leaves a gap in the numbers, which changes nothing: metric pipeline rules go by their numbers, not by counting.
- Moving a metric pipeline rule in the OneUptime Dashboard numbers the project's metric pipeline rules 1..n again; the next plan shows that as changes to `sort_order`.


## Schema

### Required

- `name` (String) Friendly name for this rule.
- `rule_type` (String) One of: Filter, Drop, RenameMetric, RenameAttribute, AddAttribute, RemoveAttribute, RedactAttribute, Sample. Allowed values: `Filter`, `Drop`, `RenameMetric`, `RenameAttribute`, `AddAttribute`, `RemoveAttribute`, `RedactAttribute`, `Sample`.

### Optional

- `add_attribute_key` (String) For AddAttribute / RemoveAttribute / RedactAttribute: the attribute key to act on.
- `add_attribute_value` (String) For AddAttribute: the attribute value to set.
- `description` (String) Description of what this rule does.
- `filter_condition` (String) How to combine filters: 'All' requires every filter to match (AND), 'Any' requires at least one to match (OR). Allowed values: `All`, `Any`. Defaults to `All`.
- `filters` (String) List of filters evaluated against each metric data point. An empty list matches every data point. A JSON value: write it with `jsonencode()`.
- `is_enabled` (Boolean) Whether this rule is active. Defaults to `true`.
- `redact_replacement` (String) For RedactAttribute: the literal string to replace the value with. Defaults to [REDACTED].
- `rename_from_key` (String) For RenameMetric: the existing metric name. For RenameAttribute: the existing attribute key.
- `rename_to_key` (String) For RenameMetric: the new metric name. For RenameAttribute: the new attribute key.
- `sample_percentage` (Number) For Sample: percentage of matched rows to keep (0-100). 100 keeps all. Defaults to `100`.
- `service_id` (String) Optional service ID scoping this rule. Null means the rule is project-wide. The ID of a `oneuptime_service`.
- `sort_order` (Number) Where this rule is evaluated among the project's metric pipeline rules, lowest number first. A new rule is added to the end of the list. Setting a number places it there and moves no other one: two at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the user who created this rule. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
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
