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
  team_id = "123e4567-e89b-12d3-a456-426614174000"
  rule_type = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `team_id` (String) A unique identifier for an object, represented as a UUID..
- `rule_type` (String) Type of compliance rule...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `enabled` (Bool) Whether this compliance rule is enabled...
- `options` (String) Additional options for this compliance rule...
- `notification_channel` (String) On-call rules only: the channel the member's rule must notify them on (Call, SMS, Push, Email, WhatsApp, Telegram, Slack, MicrosoftTeams or Webhook). Leave empty to accept any channel...
- `incident_severities` (Set) Incident and incident episode on-call rules only: the severities members must have a rule for. Leave empty to require every incident severity...
- `alert_severities` (Set) Alert and alert episode on-call rules only: the severities members must have a rule for. Leave empty to require every alert severity...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_team_compliance_setting.example <id>
```
