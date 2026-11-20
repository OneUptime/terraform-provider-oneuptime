---
page_title: "oneuptime_trace_pipeline_processor Resource - oneuptime"
subcategory: "Logs, Traces & Metrics"
description: |-
  Individual processors within a trace pipeline that transform span data during ingestion.
---

# oneuptime_trace_pipeline_processor (Resource)

Individual processors within a trace pipeline that transform span data during ingestion.

## Example Usage

```terraform
resource "oneuptime_trace_pipeline_processor" "example" {
  trace_pipeline_id = oneuptime_trace_pipeline.example.id
  name              = "Example trace pipeline processor"
  processor_type    = "AttributeRemapper"
}
```

<a id="order"></a>
## Order

`sort_order` is this trace pipeline processor's place among the trace pipeline processors of its trace pipeline, lowest first. Setting it places this trace pipeline processor and moves no other, so a configuration that numbers every trace pipeline processor of a trace pipeline is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.trace_pipeline_processors`:

```terraform
sort_order = index(local.trace_pipeline_processors, each.key) + 1
```

- Give each trace pipeline processor a number of its own: two trace pipeline processors at one number are shown the older first.
- A trace pipeline processor created without `sort_order` goes after the last one of its trace pipeline.
- Deleting a trace pipeline processor leaves a gap in the numbers, which changes nothing: trace pipeline processors go by their numbers, not by counting.
- Moving a trace pipeline processor in the OneUptime Dashboard numbers the trace pipeline's trace pipeline processors 1..n again; the next plan shows that as changes to `sort_order`.


## Schema

### Required

- `name` (String) Friendly name for this processor.
- `processor_type` (String) The type of processor: AttributeRemapper, SpanNameRemapper, StatusRemapper, SpanKindRemapper, or CategoryProcessor. Allowed values: `AttributeRemapper`, `SpanNameRemapper`, `StatusRemapper`, `SpanKindRemapper`, `CategoryProcessor`.
- `trace_pipeline_id` (String) ID of the trace pipeline this processor belongs to. The ID of a `oneuptime_trace_pipeline`.

### Optional

- `configuration` (String) Processor-specific configuration as JSON (e.g., source/target fields, mapping rules). A JSON value: write it with `jsonencode()`.
- `is_enabled` (Boolean) Whether this processor is active. Defaults to `true`.
- `sort_order` (Number) Where this processor runs within its pipeline, lowest number first. A new processor is added to the end of the list. Setting a number places it there and moves no other one: two at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this trace pipeline processor. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this trace pipeline processor belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing trace pipeline processor by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_trace_pipeline_processor.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_trace_pipeline_processor.example <id>
```
