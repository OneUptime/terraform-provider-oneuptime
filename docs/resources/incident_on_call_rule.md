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
