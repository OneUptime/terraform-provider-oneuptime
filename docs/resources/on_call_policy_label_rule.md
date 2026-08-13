---
page_title: "oneuptime_on_call_policy_label_rule Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Configure rules for automatically attaching labels to on-call policies when matching policies are created
---

# oneuptime_on_call_policy_label_rule (Resource)

Configure rules for automatically attaching labels to on-call policies when matching policies are created

## Example Usage

```terraform
resource "oneuptime_on_call_policy_label_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this on-call policy label rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this on-call policy label rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `on_call_duty_policy_labels` (Set) Only trigger for on-call policies that already have at least one of these labels. Leave empty to match regardless of labels...
- `on_call_duty_policy_name_pattern` (String) Regex (case-insensitive) matched against the on-call policy name. Leave empty to match any name...
- `on_call_duty_policy_description_pattern` (String) Regex (case-insensitive) matched against the on-call policy description. Leave empty to match any description...
- `labels_to_add` (Set) Labels to attach to the on-call policy when this rule matches. Already-attached labels are not duplicated...
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
terraform import oneuptime_on_call_policy_label_rule.example <id>
```
