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
  name = "Example short text"
  trigger_entity_type = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this investigation rule...
- `trigger_entity_type` (String) Which kind of new signal this rule decides about: Incident or Alert...

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource...
- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this investigation rule...
- `is_enabled` (Bool) Whether this rule is enabled...
- `monitors` (Set) Match only incidents/alerts from these monitors. Leave empty to match any monitor...
- `incident_severities` (Set) Match only incidents with these severities (incident rules only). Leave empty to match any severity...
- `alert_severities` (Set) Match only alerts with these severities (alert rules only). Leave empty to match any severity...
- `labels` (Set) Match only incidents/alerts that carry at least one of these labels. Leave empty to match any label...
- `monitor_labels` (Set) Match only when the incident/alert's monitor carries at least one of these labels — the natural way to scope rules to environments (e.g. staging vs production). Leave empty to match any monitor label...
- `title_pattern` (String) Case-insensitive regex matched against the entity's title. Leave empty to match any title...
- `description_pattern` (String) Case-insensitive regex matched against the entity's description. Leave empty to match any description...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_investigation_rule.example <id>
```
