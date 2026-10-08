---
page_title: "oneuptime_team_compliance_setting Resource - oneuptime"
subcategory: "Teams & Access"
description: |-
  Compliance settings for your OneUptime team
---

# oneuptime_team_compliance_setting (Resource)

Compliance settings for your OneUptime team

## Example Usage

```terraform
resource "oneuptime_team_compliance_setting" "example" {
  team_id   = oneuptime_team.example.id
  rule_type = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `rule_type` (String) Type of compliance rule.
- `team_id` (String) ID of Team this compliance setting belongs to. The ID of a `oneuptime_team`.

### Optional

- `alert_severities` (Set of String) Alert and alert episode on-call rules only: the severities members must have a rule for. Leave empty to require every alert severity. IDs of `oneuptime_alert_severity` resources.
- `enabled` (Boolean) Whether this compliance rule is enabled. Defaults to `false`.
- `incident_severities` (Set of String) Incident and incident episode on-call rules only: the severities members must have a rule for. Leave empty to require every incident severity. IDs of `oneuptime_incident_severity` resources.
- `notification_channel` (String) Deprecated: use notificationChannels. The first of the rule's notification channels, or empty when it accepts any channel. Sending this field without notificationChannels sets the rule to that one channel (Call, SMS, Push, Email, WhatsApp, Telegram, Slack, MicrosoftTeams or Webhook).
- `notification_channels` (String) On-call rules only: the channels members must be notified on, as a list - each member needs a rule on every one of them (Call, SMS, Push, Email, WhatsApp, Telegram, Slack, MicrosoftTeams or Webhook). Leave empty to accept any channel. A JSON value: write it with `jsonencode()`.
- `options` (String) Additional options for this compliance rule. A JSON value: write it with `jsonencode()`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing team compliance setting by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_team_compliance_setting.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_team_compliance_setting.example <id>
```
