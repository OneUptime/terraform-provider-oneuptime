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
  rule_type = "HasNotificationEmailMethod"
}
```

## Schema

### Required

- `rule_type` (String) Type of compliance rule. Allowed values: `HasNotificationEmailMethod`, `HasNotificationSMSMethod`, `HasNotificationCallMethod`, `HasNotificationPushMethod`, `HasNotificationWhatsAppMethod`, `HasNotificationTelegramMethod`, `HasNotificationSlackMethod`, `HasNotificationMicrosoftTeamsMethod`, `HasNotificationMattermostMethod`, `HasNotificationWebhookMethod`, `HasIncidentOnCallRules`, `HasAlertOnCallRules`, `HasIncidentEpisodeOnCallRules`, `HasAlertEpisodeOnCallRules`.
- `team_id` (String) ID of Team this compliance setting belongs to. The ID of a `oneuptime_team`.

### Optional

- `alert_severities` (Set of String) Alert and alert episode on-call rules only: the severities members must have a rule for. Leave empty to require every alert severity. IDs of `oneuptime_alert_severity` resources.
- `enabled` (Boolean) Whether this compliance rule is enabled. Defaults to `false`.
- `incident_severities` (Set of String) Incident and incident episode on-call rules only: the severities members must have a rule for. Leave empty to require every incident severity. IDs of `oneuptime_incident_severity` resources.
- `notification_channel` (String) Deprecated: use notificationChannels. The first of the rule's notification channels, or empty when it accepts any channel. Sending this field without notificationChannels sets the rule to that one channel (Call, SMS, Push, Email, WhatsApp, Telegram, Slack, MicrosoftTeams, Mattermost, Webhook or CustomChannel). Allowed values: `Call`, `SMS`, `Push`, `Email`, `WhatsApp`, `Telegram`, `Slack`, `MicrosoftTeams`, `Mattermost`, `Webhook`, `CustomChannel`.
- `notification_channels` (String) On-call rules only: the channels members must be notified on, as a list - each member needs a rule on every one of them (Call, SMS, Push, Email, WhatsApp, Telegram, Slack, MicrosoftTeams, Mattermost, Webhook or CustomChannel, which any custom channel that is on meets). Leave empty to accept any channel. A JSON value: write it with `jsonencode()`.
- `options` (String) Additional options for this compliance rule. A JSON value: write it with `jsonencode()`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
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
