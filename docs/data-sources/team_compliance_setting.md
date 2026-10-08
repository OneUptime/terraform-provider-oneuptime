---
page_title: "oneuptime_team_compliance_setting Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Compliance settings for your OneUptime team
---

# oneuptime_team_compliance_setting (Data Source)

Compliance settings for your OneUptime team

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one team compliance setting may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_team_compliance_setting" "example" {
  team_id = oneuptime_team.example.id
}

# Or by id:
data "oneuptime_team_compliance_setting" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `enabled` (Boolean) Whether this compliance rule is enabled.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `notification_channel` (String) Deprecated: use notificationChannels. The first of the rule's notification channels, or empty when it accepts any channel. Sending this field without notificationChannels sets the rule to that one channel (Call, SMS, Push, Email, WhatsApp, Telegram, Slack, MicrosoftTeams or Webhook).
- `rule_type` (String) Type of compliance rule.
- `team_id` (String) ID of Team this compliance setting belongs to. The ID of a `oneuptime_team`.

### Read-Only

- `alert_severities` (Set of String) Alert and alert episode on-call rules only: the severities members must have a rule for. Leave empty to require every alert severity. IDs of `oneuptime_alert_severity` resources.
- `created_at` (String) Date and Time when the object was created.
- `incident_severities` (Set of String) Incident and incident episode on-call rules only: the severities members must have a rule for. Leave empty to require every incident severity. IDs of `oneuptime_incident_severity` resources.
- `notification_channels` (String) On-call rules only: the channels members must be notified on, as a list - each member needs a rule on every one of them (Call, SMS, Push, Email, WhatsApp, Telegram, Slack, MicrosoftTeams or Webhook). Leave empty to accept any channel. A JSON value: write it with `jsonencode()`.
- `options` (String) Additional options for this compliance rule. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
