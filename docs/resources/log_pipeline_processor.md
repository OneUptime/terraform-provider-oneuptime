---
page_title: "oneuptime_log_pipeline_processor Resource - oneuptime"
subcategory: "Logs, Traces & Metrics"
description: |-
  Individual processors within a log pipeline that transform log data during ingestion.
---

# oneuptime_log_pipeline_processor (Resource)

Individual processors within a log pipeline that transform log data during ingestion.

## Example Usage

```terraform
resource "oneuptime_log_pipeline_processor" "example" {
  log_pipeline_id = oneuptime_log_pipeline.example.id
  name            = "Example log pipeline processor"
  processor_type  = "GrokParser"
}
```

<a id="order"></a>
## Order

`sort_order` is this log pipeline processor's place among the log pipeline processors of its log pipeline, lowest first. Setting it places this log pipeline processor and moves no other, so a configuration that numbers every log pipeline processor of a log pipeline is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.log_pipeline_processors`:

```terraform
sort_order = index(local.log_pipeline_processors, each.key) + 1
```

- Give each log pipeline processor a number of its own: two log pipeline processors at one number are shown the older first.
- A log pipeline processor created without `sort_order` goes after the last one of its log pipeline.
- Deleting a log pipeline processor leaves a gap in the numbers, which changes nothing: log pipeline processors go by their numbers, not by counting.
- Moving a log pipeline processor in the OneUptime Dashboard numbers the log pipeline's log pipeline processors 1..n again; the next plan shows that as changes to `sort_order`.


## Schema

### Required

- `log_pipeline_id` (String) ID of the log pipeline this processor belongs to. The ID of a `oneuptime_log_pipeline`.
- `name` (String) Friendly name for this processor.
- `processor_type` (String) The type of processor: GrokParser, KeyValueParser, AttributeRemapper, SeverityRemapper, or CategoryProcessor. Allowed values: `GrokParser`, `KeyValueParser`, `AttributeRemapper`, `SeverityRemapper`, `CategoryProcessor`.

### Optional

- `configuration` (String) Processor-specific configuration as JSON (e.g., grok pattern, source/target fields, mapping rules). A JSON value: write it with `jsonencode()`.
- `is_enabled` (Boolean) Whether this processor is active. Defaults to `true`.
- `sort_order` (Number) Where this processor runs within its pipeline, lowest number first. A new processor is added to the end of the list. Setting a number places it there and moves no other one: two at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this log pipeline processor. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this log pipeline processor belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing log pipeline processor by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_log_pipeline_processor.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_log_pipeline_processor.example <id>
```
