---
page_title: "oneuptime_alert_episode_privacy_rule Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Configure rules for automatically marking matching alert episodes as private
---

# oneuptime_alert_episode_privacy_rule (Resource)

Configure rules for automatically marking matching alert episodes as private

## Example Usage

```terraform
resource "oneuptime_alert_episode_privacy_rule" "example" {
  name        = "Example alert episode privacy rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this alert episode privacy rule.

### Optional

- `alert_severities` (Set of String) Only trigger for episodes with these severities. Leave empty to match episodes of any severity. IDs of `oneuptime_alert_severity` resources.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this alert episode privacy rule.
- `episode_description_pattern` (String) Regex (case-insensitive) matched against the episode description. Leave empty to match any description.
- `episode_labels` (Set of String) Only trigger for episodes that have at least one of these labels. Leave empty to match regardless of episode labels. IDs of `oneuptime_label` resources.
- `episode_title_pattern` (String) Regex (case-insensitive) matched against the episode title. Leave empty to match any title.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing alert episode privacy rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_alert_episode_privacy_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_alert_episode_privacy_rule.example <id>
```
