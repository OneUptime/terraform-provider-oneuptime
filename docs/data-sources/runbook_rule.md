---
page_title: "oneuptime_runbook_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Auto-attach runbooks to incidents, alerts, or scheduled maintenance events when they are created.
---

# oneuptime_runbook_rule (Data Source)

Auto-attach runbooks to incidents, alerts, or scheduled maintenance events when they are created. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_runbook_rule" "by_name" {
  name = "example-runbook_rule"
}

data "oneuptime_runbook_rule" "by_id" {
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
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description` (String) Description of this runbook rule... Computed.
- `is_enabled` (Bool) Whether this rule is enabled... Computed.
- `trigger_entity_type` (String) Entity type that triggers this rule on creation: Incident, Alert, or ScheduledMaintenance... Computed.
- `title_pattern` (String) Case-insensitive regex matched against the entity's title. Leave empty to match any title... Computed.
- `description_pattern` (String) Case-insensitive regex matched against the entity's description. Leave empty to match any description... Computed.
- `runbooks` (Set) Runbooks to start when this rule matches. Each runbook produces its own execution... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
