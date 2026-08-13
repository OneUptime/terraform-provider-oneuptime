---
page_title: "oneuptime_runbook_owner_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically assigning owner users and teams when matching runbooks are created
---

# oneuptime_runbook_owner_rule (Resource)

Configure rules for automatically assigning owner users and teams when matching runbooks are created

## Example Usage

```terraform
resource "oneuptime_runbook_owner_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this runbook owner rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this runbook owner rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `notify_owners` (Bool) Send notifications to owner users and teams when they are added by this rule..
- `runbook_labels` (Set) Only trigger for runbooks that have at least one of these labels. Leave empty to match regardless of labels...
- `runbook_name_pattern` (String) Regex (case-insensitive) matched against the runbook name. Leave empty to match any name...
- `runbook_description_pattern` (String) Regex (case-insensitive) matched against the runbook description. Leave empty to match any description...
- `owner_users` (Set) Users to add as owners on the runbook when this rule matches...
- `owner_teams` (Set) Teams to add as owners on the runbook when this rule matches...
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
terraform import oneuptime_runbook_owner_rule.example <id>
```
