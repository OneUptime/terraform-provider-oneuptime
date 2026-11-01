---
page_title: "oneuptime_data_security_setting Data Source - oneuptime"
subcategory: "Other"
description: |-
  How Data Security runs in a project: what it scans, how much it reads, and what it does when it finds a Critical issue.
---

# oneuptime_data_security_setting (Data Source)

How Data Security runs in a project: what it scans, how much it reads, and what it does when it finds a Critical issue.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data security setting may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_security_setting" "example" {
  alert_severity_id = oneuptime_alert_severity.example.id
}

# Or by id:
data "oneuptime_data_security_setting" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_severity_id` (String) The severity of the alerts raised. The project's most severe when empty. The ID of a `oneuptime_alert_severity`.
- `content_scan_interval_in_hours` (Number) How often a store's content is sampled again.
- `created_by_user_id` (String) ID of the user who created this settings. The ID of a `oneuptime_user` (see the data source).
- `declare_incidents_for_critical_issues` (Boolean) Declare an incident, which runs your incident process, when a Critical issue opens.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_severity_id` (String) The severity of the incidents declared. The project's most severe when empty. The ID of a `oneuptime_incident_severity`.
- `is_configuration_scan_enabled` (Boolean) Read stores' configuration and access, and run the checks.
- `is_content_scan_enabled` (Boolean) Sample stores' content for sensitive data, where the content permission was granted.
- `is_telemetry_scan_enabled` (Boolean) Sample the project's logs, traces, exceptions, replays and AI conversations for sensitive data.
- `keep_masked_samples` (Boolean) Keep up to three masked samples with each finding. Off keeps counts only.
- `max_bytes_per_object` (Number) The most bytes read from one object.
- `max_bytes_per_scan` (Number) The most bytes one content scan reads in all.
- `max_objects_per_scan` (Number) The most objects one content scan of a bucket reads.
- `max_rows_per_table` (Number) The most rows read from one table.
- `max_tables_per_scan` (Number) The most tables one database scan reads.
- `notify_owners` (Boolean) Tell a store's owners when an issue opens on it.
- `open_cases_for_critical_issues` (Boolean) Open a security case, grouped by store, when a Critical issue opens.
- `raise_alerts_for_critical_issues` (Boolean) Raise an alert, which pages whoever is on call, when a Critical issue opens.
- `telemetry_sample_size` (Number) How many recent records of each service and signal a telemetry scan reads.

### Read-Only

- `allowed_values` (String) Values that are never findings, such as your own support address. A JSON value: write it with `jsonencode()`.
- `created_at` (String) Date and Time when the object was created.
- `disabled_data_type_ids` (String) Built-in data types the project's scans do not look for. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of the project this settings belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
