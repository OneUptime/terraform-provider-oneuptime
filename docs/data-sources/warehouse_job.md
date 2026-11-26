---
page_title: "oneuptime_warehouse_job Data Source - oneuptime"
subcategory: "Other"
description: |-
  A query a warehouse runs again and again, with what it spent over the last day, its failures and the watched tables it writes.
---

# oneuptime_warehouse_job (Data Source)

A query a warehouse runs again and again, with what it spent over the last day, its failures and the watched tables it writes.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one warehouse job may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_warehouse_job" "example" {
  name = "Example warehouse job"
}

# Or by id:
data "oneuptime_warehouse_job" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `bytes_billed` (Number) BigQuery bytes billed for it in the last day.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `credits` (Number) Snowflake compute credits it spent in the last day.
- `data_store_id` (String) ID of the warehouse the job runs in. The ID of a `oneuptime_data_store`.
- `failed_run_count` (Number) How many of its runs in the last day failed.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `job_key` (String) The warehouse's own key for the query, its values left out.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) The job's tag, or what it does and the table it writes: never its query.
- `run_count` (Number) How many times it ran in the last day.
- `team` (String) The team its query tag or label names.
- `user_name` (String) The login or service account that runs it.
- `warehouse_name` (String) The warehouse or reservation it runs on.
- `writes_only_unused_tables` (Boolean) Whether every table it writes has gone unread for 90 days: spend nobody uses.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_run_at` (String) When it last ran.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
- `written_tables` (String) The watched tables it writes, each as { dataAssetId, table }. A JSON value: read it with `jsondecode()`.
