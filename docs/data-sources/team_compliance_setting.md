---
page_title: "oneuptime_team_compliance_setting Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Compliance settings for your OneUptime team
---

# oneuptime_team_compliance_setting (Data Source)

Compliance settings for your OneUptime team Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_team_compliance_setting" "by_name" {
  name = "example-team_compliance_setting"
}

data "oneuptime_team_compliance_setting" "by_id" {
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
- `team_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `rule_type` (String) Type of compliance rule... Computed.
- `enabled` (Bool) Whether this compliance rule is enabled... Computed.
- `options` (String) Additional options for this compliance rule... Computed.
- `notification_channel` (String) On-call rules only: the channel the member's rule must notify them on (Call, SMS, Push, Email, WhatsApp, Telegram, Slack, MicrosoftTeams or Webhook). Leave empty to accept any channel... Computed.
- `incident_severities` (Set) Incident and incident episode on-call rules only: the severities members must have a rule for. Leave empty to require every incident severity... Computed.
- `alert_severities` (Set) Alert and alert episode on-call rules only: the severities members must have a rule for. Leave empty to require every alert severity... Computed.
