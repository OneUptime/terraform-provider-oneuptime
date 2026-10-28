---
page_title: "oneuptime_alert_on_call_rule Data Source - oneuptime"
subcategory: "Alerts"
description: |-
  Configure rules for automatically executing on-call duty policies when matching alerts are created
---

# oneuptime_alert_on_call_rule (Data Source)

Configure rules for automatically executing on-call duty policies when matching alerts are created

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one alert on call rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_alert_on_call_rule" "example" {
  name = "Example alert on call rule"
}

# Or by id:
data "oneuptime_alert_on_call_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `ai_hold_time_in_minutes` (Number) When When to Page is After OneUptime AI: how many minutes OneUptime AI may hold the page while it investigates and fixes (1 to 30). Empty: 10. If OneUptime AI has not resolved the record by then, the policies page as usual.
- `alert_description_pattern` (String) Regular expression pattern to match alert descriptions. Leave empty to match any description.
- `alert_title_pattern` (String) Regular expression pattern to match alert titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this alert on-call rule.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule is enabled.
- `monitor_description_pattern` (String) Regular expression pattern to match monitor descriptions. Leave empty to match any monitor description.
- `monitor_name_pattern` (String) Regular expression pattern to match monitor names. Leave empty to match any monitor name. Example: 'prod-.*' matches monitors starting with 'prod-'.
- `name` (String) Name of this alert on-call rule.
- `when_to_page` (String) When the on-call policies of this rule page: Right Away, as soon as a record matches it; Working Hours, only during the working hours of this rule - a record that matches it outside them is paged about when they begin; or After OneUptime AI, for alert and incident rules - when OneUptime AI may fix what the record is about, it responds first, and the policies page only if it has not resolved the record within the hold time.
- `working_hours_timezone` (String) The IANA time zone the working hours of this rule are in: Europe/London.

### Read-Only

- `alert_labels` (Set of String) Only trigger for alerts that have at least one of these labels. Leave empty to match alerts regardless of alert labels. IDs of `oneuptime_label` resources.
- `alert_severities` (Set of String) Only trigger for alerts with these severities. Leave empty to match alerts of any severity. IDs of `oneuptime_alert_severity` resources.
- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `monitor_labels` (Set of String) Only trigger for alerts from monitors that have at least one of these labels. Leave empty to match alerts regardless of monitor labels. IDs of `oneuptime_label` resources.
- `monitors` (Set of String) Only trigger for alerts from these monitors. Leave empty to match alerts from any monitor. IDs of `oneuptime_monitor` resources.
- `on_call_duty_policies` (Set of String) On-call duty policies to execute when an alert matches this rule. IDs of `oneuptime_on_call_policy` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
- `working_hours` (String) When this rule may page, when When to Page is Working Hours: the same daily or weekly windows as an on-call schedule layer's restriction times, read in the time zone of this rule. A JSON value: write it with `jsonencode()`.
