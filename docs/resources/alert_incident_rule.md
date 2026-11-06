---
page_title: "oneuptime_alert_incident_rule Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Declare an incident automatically from alerts that match, right away or once an alert has stayed unresolved for a while. The first rule that matches, from the top of the list down, decides.
---

# oneuptime_alert_incident_rule (Resource)

Declare an incident automatically from alerts that match, right away or once an alert has stayed unresolved for a while. The first rule that matches, from the top of the list down, decides.

## Example Usage

```terraform
resource "oneuptime_alert_incident_rule" "example" {
  name        = "Example alert incident rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this alert incident rule.

### Optional

- `alert_description_pattern` (String) Regular expression pattern to match alert descriptions. Leave empty to match any description.
- `alert_labels` (Set of String) Only declare incidents from alerts that have at least one of these labels. Leave empty to match alerts regardless of alert labels. IDs of `oneuptime_label` resources.
- `alert_severities` (Set of String) Only declare incidents from alerts with these severities. Leave empty to match alerts of any severity. IDs of `oneuptime_alert_severity` resources.
- `alert_title_pattern` (String) Regular expression pattern to match alert titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `declare_after_in_minutes` (Number) How long a matching alert must stay unresolved before the incident is declared, in whole minutes from when the alert was created. 0 declares it as soon as the alert is created. An alert resolved before the wait is up never becomes an incident. Defaults to `0`.
- `description` (String) Description of this alert incident rule.
- `incident_severity_id` (String) ID of the severity of the incidents this rule declares. Empty gives each incident the severity that corresponds to its alert's. The ID of a `oneuptime_incident_severity`.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `monitor_description_pattern` (String) Regular expression pattern to match monitor descriptions. Leave empty to match any monitor description.
- `monitor_labels` (Set of String) Only declare incidents from alerts raised by monitors that have at least one of these labels. Leave empty to match alerts regardless of monitor labels. IDs of `oneuptime_label` resources.
- `monitor_name_pattern` (String) Regular expression pattern to match monitor names. Leave empty to match any monitor name. Example: 'prod-.*' matches monitors starting with 'prod-'.
- `monitors` (Set of String) Only declare incidents from alerts raised by these monitors. Leave empty to match alerts from any monitor. IDs of `oneuptime_monitor` resources.
- `order` (Number) Where this rule sits in the list. Rules are checked from the top of the list down, lowest number first, and the first enabled one that matches an alert decides. A new rule is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.
- `show_incident_on_status_page` (Boolean) Whether the incidents this rule declares are shown on the status pages that list their monitor, and their subscribers notified, as for any incident declared with Show on Status Page on. Defaults to `true`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing alert incident rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_alert_incident_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_alert_incident_rule.example <id>
```
