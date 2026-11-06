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
  data_pipeline_id = data.oneuptime_data_pipeline.example.id
}

# Or by id:
data "oneuptime_data_pipeline_run" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `attempt_number` (Number) Which try of its task this run is, when its scheduler says (Airflow's try number).
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `data_pipeline_id` (String) ID of the pipeline this is a run of. The ID of a `oneuptime_data_pipeline` (see the data source).
- `duration_in_seconds` (Number) How long the run took, in seconds.
- `error_message` (String) The error a failed run reported.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `run_id` (String) The run's OpenLineage id.
- `state` (String) How the run stands: Running, Complete, Fail or Abort.
- `trace_id` (String) The trace of a run found in traces: a CronJob or a scheduled function.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `ended_at` (String) When the run ended.
- `inputs` (String) The datasets the run read, as namespace and name. A JSON value: read it with `jsondecode()`.
- `output_data_asset_ids` (String) The watched tables among the datasets the run wrote. A JSON value: read it with `jsondecode()`.
- `outputs` (String) The datasets the run wrote, as namespace and name. A JSON value: read it with `jsondecode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `scheduled_at` (String) When the run was scheduled for, when its producer says (OpenLineage's nominalTime facet): a backfilled run's scheduled time, not when it ran.
- `started_at` (String) When the run started.
- `updated_at` (String) Date and Time when the object was updated.
