---
page_title: "oneuptime_alert_episode_label_rule Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Configure rules for automatically attaching labels to alert episodes when matching episodes are created
---

# oneuptime_alert_episode_label_rule (Resource)

Configure rules for automatically attaching labels to alert episodes when matching episodes are created

## Example Usage

```terraform
resource "oneuptime_alert_episode_label_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this alert episode label rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this alert episode label rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `alert_severities` (Set) Only trigger for episodes with these severities. Leave empty to match any severity...
- `episode_labels` (Set) Only trigger for episodes that already have at least one of these labels...
- `episode_title_pattern` (String) Regex (case-insensitive) matched against the episode title...
- `episode_description_pattern` (String) Regex (case-insensitive) matched against the episode description...
- `labels_to_add` (Set) Labels to attach to the episode when this rule matches. Already-attached labels are not duplicated...
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
terraform import oneuptime_alert_episode_label_rule.example <id>
```
