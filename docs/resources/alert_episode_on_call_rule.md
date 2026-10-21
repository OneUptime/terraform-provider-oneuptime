---
page_title: "oneuptime_alert_episode_on_call_rule Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Configure rules for automatically executing on-call duty policies when matching alert episodes are created
---

# oneuptime_alert_episode_on_call_rule (Resource)

Configure rules for automatically executing on-call duty policies when matching alert episodes are created

## Example Usage

```terraform
resource "oneuptime_alert_episode_on_call_rule" "example" {
  name        = "Example alert episode on call rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this alert episode on-call rule.

### Optional

- `alert_severities` (Set of String) Only trigger for episodes with these severities. Leave empty to match any severity. IDs of `oneuptime_alert_severity` resources.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this alert episode on-call rule.
- `episode_description_pattern` (String) Regex (case-insensitive) matched against the episode description.
- `episode_labels` (Set of String) Only trigger for episodes that have at least one of these labels. IDs of `oneuptime_label` resources.
- `episode_title_pattern` (String) Regex (case-insensitive) matched against the episode title.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `on_call_duty_policies` (Set of String) On-call duty policies to execute when an alert episode matches this rule. IDs of `oneuptime_on_call_policy` resources.
- `when_to_page` (String) When the on-call policies of this rule page: Right Away, as soon as a record matches it, or Working Hours, only during the working hours of this rule - a record that matches it outside them is paged about when they begin. Defaults to `Right Away`.
- `working_hours` (String) When this rule may page, when When to Page is Working Hours: the same daily or weekly windows as an on-call schedule layer's restriction times, read in the time zone of this rule. A JSON value: write it with `jsonencode()`.
- `working_hours_timezone` (String) The IANA time zone the working hours of this rule are in: Europe/London.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing alert episode on call rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_alert_episode_on_call_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_alert_episode_on_call_rule.example <id>
```
