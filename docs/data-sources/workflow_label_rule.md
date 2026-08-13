---
page_title: "oneuptime_workflow_label_rule Data Source - oneuptime"
subcategory: "Workflows"
description: |-
  Configure rules for automatically attaching labels to workflows when matching workflows are created
---

# oneuptime_workflow_label_rule (Data Source)

Configure rules for automatically attaching labels to workflows when matching workflows are created Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_workflow_label_rule" "by_name" {
  name = "example-workflow_label_rule"
}

data "oneuptime_workflow_label_rule" "by_id" {
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
- `description` (String) Description of this workflow label rule.. Computed.
- `is_enabled` (Bool) Whether this rule is enabled.. Computed.
- `workflow_labels` (Set) Only trigger for workflows that already have at least one of these labels. Leave empty to match regardless of labels... Computed.
- `workflow_name_pattern` (String) Regex (case-insensitive) matched against the workflow name. Leave empty to match any name... Computed.
- `workflow_description_pattern` (String) Regex (case-insensitive) matched against the workflow description. Leave empty to match any description... Computed.
- `labels_to_add` (Set) Labels to attach to the workflow when this rule matches. Already-attached labels are not duplicated... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
