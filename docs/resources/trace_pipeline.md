---
page_title: "oneuptime_trace_pipeline Resource - oneuptime"
subcategory: "Logs, Traces & Metrics"
description: |-
  Configure server-side trace processing pipelines that transform spans at ingest time.
---

# oneuptime_trace_pipeline (Resource)

Configure server-side trace processing pipelines that transform spans at ingest time.

## Example Usage

```terraform
resource "oneuptime_trace_pipeline" "example" {
  name        = "Example trace pipeline"
  description = "Managed by Terraform"
}
```

<a id="order"></a>
## Order

`sort_order` is this trace pipeline's place among the trace pipelines of its project, lowest first. Setting it places this trace pipeline and moves no other, so a configuration that numbers every trace pipeline of a project is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.trace_pipelines`:

```terraform
sort_order = index(local.trace_pipelines, each.key) + 1
```

- Give each trace pipeline a number of its own: two trace pipelines at one number are shown the older first.
- A trace pipeline created without `sort_order` goes after the last one of its project.
- Deleting a trace pipeline leaves a gap in the numbers, which changes nothing: trace pipelines go by their numbers, not by counting.
- Moving a trace pipeline in the OneUptime Dashboard numbers the project's trace pipelines 1..n again; the next plan shows that as changes to `sort_order`.


## Schema

### Required

- `name` (String) Friendly name for this trace pipeline.

### Optional

- `description` (String) Description of what this trace pipeline does.
- `filter_query` (String) Filter expression that determines which spans this pipeline applies to.
- `is_enabled` (Boolean) Whether this trace pipeline is active. Defaults to `true`.
- `sort_order` (Number) Where this pipeline runs among the project's trace pipelines, lowest number first. A new pipeline is added to the end of the list. Setting a number places it there and moves no other one: two at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this trace pipeline. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this trace pipeline belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing trace pipeline by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_trace_pipeline.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_trace_pipeline.example <id>
```
