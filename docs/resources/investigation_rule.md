---
page_title: "oneuptime_investigation_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Choose which new incidents or alerts OneUptime AI investigates on its own. With no rule, every one is investigated; with rules, only those that match one.
---

# oneuptime_investigation_rule (Resource)

Choose which new incidents or alerts OneUptime AI investigates on its own. With no rule, every one is investigated; with rules, only those that match one.

## Example Usage

```terraform
resource "oneuptime_investigation_rule" "example" {
  name                = "Example investigation rule"
  trigger_entity_type = "Example short text"
  description         = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this investigation rule.
- `trigger_entity_type` (String) Which kind of new signal this rule decides about: Incident or Alert.

### Optional

- `alert_severities` (Set of String) Match only alerts with these severities (alert rules only). Leave empty to match any severity. IDs of `oneuptime_alert_severity` resources.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this investigation rule.
- `description_pattern` (String) Case-insensitive regex matched against the entity's description. Leave empty to match any description.
- `incident_severities` (Set of String) Match only incidents with these severities (incident rules only). Leave empty to match any severity. IDs of `oneuptime_incident_severity` resources.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `labels` (Set of String) Match only incidents/alerts that carry at least one of these labels. Leave empty to match any label. IDs of `oneuptime_label` resources.
- `monitor_labels` (Set of String) Match only when the incident/alert's monitor carries at least one of these labels — the natural way to scope rules to environments (e.g. staging vs production). Leave empty to match any monitor label. IDs of `oneuptime_label` resources.
- `monitors` (Set of String) Match only incidents/alerts from these monitors. Leave empty to match any monitor. IDs of `oneuptime_monitor` resources.
- `title_pattern` (String) Case-insensitive regex matched against the entity's title. Leave empty to match any title.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing investigation rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_investigation_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_investigation_rule.example <id>
```
