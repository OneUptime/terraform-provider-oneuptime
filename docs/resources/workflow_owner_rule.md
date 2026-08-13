---
page_title: "oneuptime_workflow_owner_rule Resource - oneuptime"
subcategory: "Workflows"
description: |-
  Configure rules for automatically assigning owner users and teams when matching workflows are created
---

# oneuptime_workflow_owner_rule (Resource)

Configure rules for automatically assigning owner users and teams when matching workflows are created

## Example Usage

```terraform
resource "oneuptime_workflow_owner_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this workflow owner rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this workflow owner rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `notify_owners` (Bool) Send notifications to owner users and teams when they are added by this rule..
- `workflow_labels` (Set) Only trigger for workflows that have at least one of these labels. Leave empty to match regardless of labels...
- `workflow_name_pattern` (String) Regex (case-insensitive) matched against the workflow name. Leave empty to match any name...
- `workflow_description_pattern` (String) Regex (case-insensitive) matched against the workflow description. Leave empty to match any description...
- `owner_users` (Set) Users to add as owners on the workflow when this rule matches...
- `owner_teams` (Set) Teams to add as owners on the workflow when this rule matches...
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
terraform import oneuptime_workflow_owner_rule.example <id>
```
