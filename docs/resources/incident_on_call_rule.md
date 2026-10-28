---
page_title: "oneuptime_incident_on_call_rule Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Configure rules for automatically executing on-call duty policies when matching incidents are created
---

# oneuptime_incident_on_call_rule (Resource)

Configure rules for automatically executing on-call duty policies when matching incidents are created

## Example Usage

```terraform
resource "oneuptime_incident_on_call_rule" "example" {
  name        = "Example incident on call rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this incident on-call rule.

### Optional

- `ai_hold_time_in_minutes` (Number) When When to Page is After OneUptime AI: how many minutes OneUptime AI may hold the page while it investigates and fixes (1 to 30). Empty: 10. If OneUptime AI has not resolved the record by then, the policies page as usual.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this incident on-call rule.
- `incident_description_pattern` (String) Regular expression pattern to match incident descriptions. Leave empty to match any description.
- `incident_labels` (Set of String) Only trigger for incidents that have at least one of these labels. Leave empty to match incidents regardless of incident labels. IDs of `oneuptime_label` resources.
- `incident_severities` (Set of String) Only trigger for incidents with these severities. Leave empty to match incidents of any severity. IDs of `oneuptime_incident_severity` resources.
- `incident_title_pattern` (String) Regular expression pattern to match incident titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `monitor_description_pattern` (String) Regular expression pattern to match monitor descriptions. Leave empty to match any monitor description.
- `monitor_labels` (Set of String) Only trigger for incidents from monitors that have at least one of these labels. Leave empty to match incidents regardless of monitor labels. IDs of `oneuptime_label` resources.
- `monitor_name_pattern` (String) Regular expression pattern to match monitor names. Leave empty to match any monitor name. Example: 'prod-.*' matches monitors starting with 'prod-'.
- `monitors` (Set of String) Only trigger for incidents from these monitors. Leave empty to match incidents from any monitor. IDs of `oneuptime_monitor` resources.
- `on_call_duty_policies` (Set of String) On-call duty policies to execute when an incident matches this rule. IDs of `oneuptime_on_call_policy` resources.
- `when_to_page` (String) When the on-call policies of this rule page: Right Away, as soon as a record matches it; Working Hours, only during the working hours of this rule - a record that matches it outside them is paged about when they begin; or After OneUptime AI, for alert and incident rules - when OneUptime AI may fix what the record is about, it responds first, and the policies page only if it has not resolved the record within the hold time. Defaults to `Right Away`.
- `working_hours` (String) When this rule may page, when When to Page is Working Hours: the same daily or weekly windows as an on-call schedule layer's restriction times, read in the time zone of this rule. A JSON value: write it with `jsonencode()`.
- `working_hours_timezone` (String) The IANA time zone the working hours of this rule are in: Europe/London.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incident on call rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incident_on_call_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incident_on_call_rule.example <id>
```
