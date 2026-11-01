---
page_title: "oneuptime_data_security_setting Resource - oneuptime"
subcategory: "Other"
description: |-
  How Data Security runs in a project: what it scans, how much it reads, and what it does when it finds a Critical issue.
---

# oneuptime_data_security_setting (Resource)

How Data Security runs in a project: what it scans, how much it reads, and what it does when it finds a Critical issue.

## Example Usage

```terraform
resource "oneuptime_data_security_setting" "example" {

}
```

## Schema

### Optional

- `alert_severity_id` (String) The severity of the alerts raised. The project's most severe when empty. The ID of a `oneuptime_alert_severity`.
- `allowed_values` (String) Values that are never findings, such as your own support address. A JSON value: write it with `jsonencode()`.
- `content_scan_interval_in_hours` (Number) How often a store's content is sampled again. Defaults to `24`.
- `declare_incidents_for_critical_issues` (Boolean) Declare an incident, which runs your incident process, when a Critical issue opens. Defaults to `false`.
- `disabled_data_type_ids` (String) Built-in data types the project's scans do not look for. A JSON value: write it with `jsonencode()`.
- `incident_severity_id` (String) The severity of the incidents declared. The project's most severe when empty. The ID of a `oneuptime_incident_severity`.
- `is_configuration_scan_enabled` (Boolean) Read stores' configuration and access, and run the checks. Defaults to `true`.
- `is_content_scan_enabled` (Boolean) Sample stores' content for sensitive data, where the content permission was granted. Defaults to `true`.
- `is_telemetry_scan_enabled` (Boolean) Sample the project's logs, traces, exceptions, replays and AI conversations for sensitive data. Defaults to `true`.
- `keep_masked_samples` (Boolean) Keep up to three masked samples with each finding. Off keeps counts only. Defaults to `true`.
- `max_bytes_per_object` (Number) The most bytes read from one object. Defaults to `65536`.
- `max_bytes_per_scan` (Number) The most bytes one content scan reads in all. Defaults to `5.24288e+07`.
- `max_objects_per_scan` (Number) The most objects one content scan of a bucket reads. Defaults to `2000`.
- `max_rows_per_table` (Number) The most rows read from one table. Defaults to `1000`.
- `max_tables_per_scan` (Number) The most tables one database scan reads. Defaults to `200`.
- `notify_owners` (Boolean) Tell a store's owners when an issue opens on it. Defaults to `true`.
- `open_cases_for_critical_issues` (Boolean) Open a security case, grouped by store, when a Critical issue opens. Defaults to `true`.
- `raise_alerts_for_critical_issues` (Boolean) Raise an alert, which pages whoever is on call, when a Critical issue opens. Defaults to `false`.
- `telemetry_sample_size` (Number) How many recent records of each service and signal a telemetry scan reads. Defaults to `500`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this settings. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this settings belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing data security setting by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_data_security_setting.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_data_security_setting.example <id>
```
