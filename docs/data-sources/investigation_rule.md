---
page_title: "oneuptime_investigation_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Choose which new incidents or alerts OneUptime AI investigates on its own. With no rule, every one is investigated; with rules, only those that match one.
---

# oneuptime_investigation_rule (Data Source)

Choose which new incidents or alerts OneUptime AI investigates on its own. With no rule, every one is investigated; with rules, only those that match one. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_investigation_rule" "by_name" {
  name = "example-investigation_rule"
}

data "oneuptime_investigation_rule" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource... Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description` (String) Description of this investigation rule... Computed.
- `is_enabled` (Bool) Whether this rule is enabled... Computed.
- `trigger_entity_type` (String) Which kind of new signal this rule decides about: Incident or Alert... Computed.
- `monitors` (Set) Match only incidents/alerts from these monitors. Leave empty to match any monitor... Computed.
- `incident_severities` (Set) Match only incidents with these severities (incident rules only). Leave empty to match any severity... Computed.
- `alert_severities` (Set) Match only alerts with these severities (alert rules only). Leave empty to match any severity... Computed.
- `labels` (Set) Match only incidents/alerts that carry at least one of these labels. Leave empty to match any label... Computed.
- `monitor_labels` (Set) Match only when the incident/alert's monitor carries at least one of these labels — the natural way to scope rules to environments (e.g. staging vs production). Leave empty to match any monitor label... Computed.
- `title_pattern` (String) Case-insensitive regex matched against the entity's title. Leave empty to match any title... Computed.
- `description_pattern` (String) Case-insensitive regex matched against the entity's description. Leave empty to match any description... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
