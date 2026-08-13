---
page_title: "oneuptime_runbook_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Auto-attach runbooks to incidents, alerts, or scheduled maintenance events when they are created.
---

# oneuptime_runbook_rule (Resource)

Auto-attach runbooks to incidents, alerts, or scheduled maintenance events when they are created.

## Example Usage

```terraform
resource "oneuptime_runbook_rule" "example" {
  name = "Example short text"
  trigger_entity_type = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this runbook rule...
- `trigger_entity_type` (String) Entity type that triggers this rule on creation: Incident, Alert, or ScheduledMaintenance...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this runbook rule...
- `is_enabled` (Bool) Whether this rule is enabled...
- `title_pattern` (String) Case-insensitive regex matched against the entity's title. Leave empty to match any title...
- `description_pattern` (String) Case-insensitive regex matched against the entity's description. Leave empty to match any description...
- `runbooks` (Set) Runbooks to start when this rule matches. Each runbook produces its own execution...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_runbook_rule.example <id>
```
