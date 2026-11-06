---
page_title: "oneuptime_data_pipeline Data Source - oneuptime"
subcategory: "Other"
description: |-
  A job that moves data - an Airflow task, a Spark job, a dbt model - as its OpenLineage events describe it, with how its last run went.
---

# oneuptime_data_pipeline (Data Source)

A job that moves data - an Airflow task, a Spark job, a dbt model - as its OpenLineage events describe it, with how its last run went.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data pipeline may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_pipeline" "example" {
  job_namespace = "example-job-namespace"
}

# Or by id:
data "oneuptime_data_pipeline" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) What the pipeline does and who relies on it, in your own words.
- `expected_run_interval_in_seconds` (Number) The longest the pipeline goes between runs on its schedule, in seconds: the longest gap between the starts of its last 50 runs, once there are five.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `job_name` (String) The job's OpenLineage name, such as orders_daily.load_orders.
- `job_namespace` (String) The job's OpenLineage namespace: where it runs, such as airflow://prod.
- `last_error_message` (String) The error the latest failed run reported.
- `last_run_duration_in_seconds` (Number) How long the latest finished run took, in seconds.
- `last_run_id` (String) The OpenLineage run id of the pipeline's latest run.
- `last_run_state` (String) How the latest run stands: Running, Complete, Fail or Abort.
- `pipeline_key` (String) The job's namespace and name together: what its events are matched by.
- `run_count` (Number) How many runs OneUptime has seen.
- `typical_run_duration_in_seconds` (Number) How long a run usually takes, in seconds: the median of the last 20 completed runs, once there are five.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `first_seen_at` (String) When the pipeline's first event arrived.
- `inputs` (String) The datasets the pipeline reads, as namespace and name. A JSON value: read it with `jsondecode()`.
- `last_failure_at` (String) When a run last failed or was aborted.
- `last_run_ended_at` (String) When the latest run ended.
- `last_run_started_at` (String) When the latest run started.
- `last_success_at` (String) When a run last completed.
- `outputs` (String) The datasets the pipeline writes, as namespace and name. A JSON value: read it with `jsondecode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
