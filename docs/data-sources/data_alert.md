---
page_title: "oneuptime_data_alert Data Source - oneuptime"
subcategory: "Other"
description: |-
  What a data monitor's alert is about - a table, a pipeline or a store - and the alert it is a consequence of, found by lineage.
---

# oneuptime_data_alert (Data Source)

What a data monitor's alert is about - a table, a pipeline or a store - and the alert it is a consequence of, found by lineage.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data alert may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_alert" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_data_alert" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_id` (String) ID of the alert this is about. The ID of a `oneuptime_alert`.
- `cause_alert_id` (String) ID of the alert this one is a consequence of, found by lineage. The ID of a `oneuptime_alert`.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `data_asset_id` (String) ID of the table the alert is about, when it is about one. The ID of a `oneuptime_data_asset`.
- `data_pipeline_id` (String) ID of the pipeline the alert is about, when it is about one. The ID of a `oneuptime_data_pipeline` (see the data source).
- `data_store_id` (String) ID of the database the alert's table lives in, or of the database itself when the whole of it stopped reporting. The ID of a `oneuptime_data_store`.
- `detected_after_in_seconds` (Number) How long after a table's update was due the alert opened; empty for an alert about anything else.
- `feedback` (String) What whoever handled the alert said of it: useful, expected, or not an issue.
- `feedback_by_user_id` (String) Who gave the feedback.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_cause_confirmed` (Boolean) Whether the likely cause OneUptime named was the real one, once someone says.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `feedback_at` (String) When the feedback was given.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
