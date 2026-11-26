---
page_title: "oneuptime_data_pipeline_run Data Source - oneuptime"
subcategory: "Other"
description: |-
  One run of a data pipeline: when it started and ended, how it went, and the datasets it read and wrote.
---

# oneuptime_data_pipeline_run (Data Source)

One run of a data pipeline: when it started and ended, how it went, and the datasets it read and wrote.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data pipeline run may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_pipeline_run" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_data_pipeline_run" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `attempt_number` (Number) Which try of its task this run is, when its scheduler says (Airflow's try number).
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `data_pipeline_id` (String) ID of the pipeline this is a run of. The ID of a `oneuptime_data_pipeline` (see the data source).
- `duration_in_seconds` (Number) How long the run took, in seconds.
- `error_message` (String) The error a failed run reported.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `run_id` (String) The run's OpenLineage id.
- `state` (String) How the run stands: Running, Complete, Fail or Abort.
- `trace_id` (String) The trace of a run found in traces: a CronJob or a scheduled function.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `ended_at` (String) When the run ended.
- `inputs` (String) The datasets the run read, as namespace and name. A JSON value: read it with `jsondecode()`.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `output_data_asset_ids` (String) The watched tables among the datasets the run wrote. A JSON value: read it with `jsondecode()`.
- `outputs` (String) The datasets the run wrote, as namespace and name. A JSON value: read it with `jsondecode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `scheduled_at` (String) When the run was scheduled for, when its producer says (OpenLineage's nominalTime facet): a backfilled run's scheduled time, not when it ran.
- `started_at` (String) When the run started.
- `updated_at` (String) Date and Time when the object was updated.
