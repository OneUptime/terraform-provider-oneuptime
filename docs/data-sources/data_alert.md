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
  alert_id = oneuptime_alert.example.id
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
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `data_asset_id` (String) ID of the table the alert is about, when it is about one. The ID of a `oneuptime_data_asset`.
- `data_pipeline_id` (String) ID of the pipeline the alert is about, when it is about one. The ID of a `oneuptime_data_pipeline` (see the data source).
- `data_store_id` (String) ID of the database the alert's table lives in, or of the database itself when the whole of it stopped reporting. The ID of a `oneuptime_data_store`.
- `detected_after_in_seconds` (Number) How long after a table's update was due the alert opened; empty for an alert about anything else.
- `feedback` (String) What whoever handled the alert said of it: useful, expected, or not an issue.
- `feedback_by_user_id` (String) Who gave the feedback.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_cause_confirmed` (Boolean) Whether the likely cause OneUptime named was the real one, once someone says.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `feedback_at` (String) When the feedback was given.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
