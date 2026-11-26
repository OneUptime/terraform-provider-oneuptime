---
page_title: "oneuptime_log_pipeline Resource - oneuptime"
subcategory: "Logs, Traces & Metrics"
description: |-
  Configure server-side log processing pipelines that transform logs at ingest time.
---

# oneuptime_log_pipeline (Resource)

Configure server-side log processing pipelines that transform logs at ingest time.

## Example Usage

```terraform
resource "oneuptime_log_pipeline" "example" {
  name        = "Example log pipeline"
  description = "Managed by Terraform"
}
```

<a id="order"></a>
## Order

`sort_order` is this log pipeline's place among the log pipelines of its project, lowest first. Setting it places this log pipeline and moves no other, so a configuration that numbers every log pipeline of a project is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.log_pipelines`:

```terraform
sort_order = index(local.log_pipelines, each.key) + 1
```

- Give each log pipeline a number of its own: two log pipelines at one number are shown the older first.
- A log pipeline created without `sort_order` goes after the last one of its project.
- Deleting a log pipeline leaves a gap in the numbers, which changes nothing: log pipelines go by their numbers, not by counting.
- Moving a log pipeline in the OneUptime Dashboard numbers the project's log pipelines 1..n again; the next plan shows that as changes to `sort_order`.


## Schema

### Required

- `name` (String) Friendly name for this log pipeline.

### Optional

- `description` (String) Description of what this log pipeline does.
- `filter_query` (String) Filter expression that determines which logs this pipeline applies to.
- `is_enabled` (Boolean) Whether this log pipeline is active. Defaults to `true`.
- `sort_order` (Number) Where this pipeline runs among the project's log pipelines, lowest number first. A new pipeline is added to the end of the list. Setting a number places it there and moves no other one: two at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the user who created this log pipeline. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of the project this log pipeline belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing log pipeline by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_log_pipeline.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_log_pipeline.example <id>
```
